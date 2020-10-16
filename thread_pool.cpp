#include <iostream>
#include <utility>
#include <thread>
#include <chrono>
#include <functional>
#include <atomic>
#include <queue>
#include <mutex>
#include <condition_variable>
#include <csignal>
#include <atomic>
#include <vector>
#include <list>

class TaskData
{
public:
	int data;
};

static long long get_cur_time()
{
	return std::chrono::duration_cast<std::chrono::milliseconds>(std::chrono::system_clock::now().time_since_epoch()).count();
}

typedef std::chrono::system_clock::time_point task_tm;

class WorkerThread
{
public:
	WorkerThread(int thread_id, int max_queue_size) :thread_id_(thread_id)
		, max_queue_size_(max_queue_size)
		, task_size_(0)
		, thread_runing_(true)
	{

	}
	virtual ~WorkerThread()
	{
		printf("~WorkerThread %d\n", thread_id_);
		stop();
		if (thd_)
		{
			thd_->join();
			delete thd_;
		}
	}
	void run()
	{
		while (thread_runing_)
		{
			std::queue<TaskData> temp;
			{
				std::unique_lock<std::mutex> guard(task_lock_);
				if (!tasks_.empty())
				{
					std::swap(temp, tasks_);
				}
			}
			{
				std::unique_lock<std::mutex> guard(thread_lock_);
				if (temp.empty() && thread_runing_)
				{
					con_.wait(guard, [] {std::chrono::milliseconds(10); return true; });
					continue;
				}
			}
			while (!temp.empty())
			{
				TaskData& t = temp.front();
				temp.pop();
				printf("task id:%d\n", t.data);
				--task_size_;
			}
		}
	}

	void start()
	{
		thd_ = new std::thread(&WorkerThread::run, this);
	}

	bool push(TaskData& data) {
		if (!thread_runing_)
		{
			return false;
		}

		bool is_notify = false;
		{
			std::unique_lock<std::mutex> guard(task_lock_);
			int total_size = tasks_.size();
			if (total_size >= max_queue_size_)
			{
				printf("max size\n");
				return false;
			}

			is_notify = total_size == 0;
			tasks_.push(data);
			++task_size_;
		}

		if (!thread_runing_)
		{
			return false;
		}

		last_active_ms_ = std::chrono::system_clock::now();
		if (is_notify)
		{
			std::unique_lock<std::mutex> guard(thread_lock_);
			con_.notify_one();
		}
		return true;
	}
	void wait()
	{
		if (thd_)
		{
			thd_->join();
		}
	}

	void stop()
	{
		printf("try stop thread %d\n", thread_id_);
		thread_runing_ = false;
	}

	bool is_full()
	{
		return task_size_ >= max_queue_size_;
	}

	bool is_empty()
	{
		return task_size_ < 1;
	}
	int get_thread_id()
	{
		return thread_id_;
	}

	task_tm get_last_active_ms()
	{
		return last_active_ms_;
	}

private:
	std::queue<TaskData> tasks_;
	std::atomic<bool> thread_runing_;
	std::mutex task_lock_;
	std::mutex thread_lock_;
	std::condition_variable con_;
	int max_queue_size_;
	std::thread* thd_ = NULL;
	int thread_id_;
	std::atomic<int> task_size_;
	std::chrono::system_clock::time_point last_active_ms_ = std::chrono::system_clock::now();
};

typedef WorkerThread* pWorkerThread;

class ThreadPool
{
public:
	typedef std::list<pWorkerThread> thread_vec;
	ThreadPool(int min_thread_num, int max_thread_num, int max_queue_size) :
		min_thread_num_(min_thread_num)
		, max_thread_num_(max_thread_num)
		, max_queue_size_(max_queue_size)
		, is_runing_(true)
	{
	}
	~ThreadPool()
	{
		stop();
	}
	void start()
	{
		check_min_threads();
	}

	bool push(TaskData& t)
	{
		if (!is_runing_)
		{
			printf("threadpool has stop\n");
			return false;
		}
		bool is_create_new = check_min_threads();
		if (!is_create_new)
		{
			shrink_threads();
		}
		return handle_task(t);
	}

	bool handle_task(TaskData& t)
	{
		bool is_handled = false;
		std::unique_lock<std::mutex> guard(thread_lock_);
		for (thread_vec::iterator iter = thd_vec_.begin();
			iter != thd_vec_.end(); ++iter)
		{
			pWorkerThread pw = *iter;
			if (NULL == pw)
			{
				continue;
			}
			if (pw->is_full())
			{
				continue;
			}
			is_handled = pw->push(t);
			break;
		}

		if (is_handled)
		{
			return true;
		}

		int thd_size = thd_vec_.size();
		if (thd_size < max_thread_num_)
		{
			int thd_id = thd_size + 1;
			WorkerThread* pw = new WorkerThread(thd_id, max_queue_size_);
			pw->start();
			thd_vec_.push_back(pw);
			pw->push(t);
			printf("create thread %d\n", thd_id);
			return true;
		}
		return false;
	}

	void stop()
	{
		printf("stop threadpool\n");
		is_runing_ = false;
		std::unique_lock<std::mutex> guard(thread_lock_);
		for (thread_vec::iterator iter = thd_vec_.begin();
			iter != thd_vec_.end(); ++iter)
		{
			pWorkerThread pw = *iter;
			if (pw)
			{
				pw->stop();
			}
		}
		for (thread_vec::iterator iter = thd_vec_.begin();
			iter != thd_vec_.end(); ++iter)
		{
			pWorkerThread pw = *iter;
			if (pw)
			{
				delete pw;
			}
		}
	}

private:
	bool check_min_threads()
	{
		std::unique_lock<std::mutex> guard(thread_lock_);
		int thd_size = thd_vec_.size();
		if (thd_size < min_thread_num_)
		{
			for (int i = thd_size + 1; i <= min_thread_num_; ++i)
			{
				WorkerThread* pw = new WorkerThread(i, max_queue_size_);
				pw->start();
				thd_vec_.push_back(pw);
				printf("create thread %d\n", i);
			}
			return true;
		}
		return false;
	}

	bool shrink_threads()
	{
		std::unique_lock<std::mutex> guard(thread_lock_);
		int thd_size = thd_vec_.size();
		if (thd_size <= min_thread_num_)
		{
			return false;
		}
		int i = 0;
		auto now = std::chrono::system_clock::now();
		for (thread_vec::iterator iter = thd_vec_.begin();
			iter != thd_vec_.end(); )
		{
			++i;
			if (i <= min_thread_num_)
			{
				++iter;
				continue;
			}
			pWorkerThread pw = *iter;
			if (pw && pw->is_empty())
			{
				long long unlive_ms = std::chrono::duration_cast<std::chrono::milliseconds>(now - pw->get_last_active_ms()).count();
				if (unlive_ms < keep_live_time_)
				{
					continue;
				}
				printf("shrink_threads %d unlive_ms:%lld\n", pw->get_thread_id(), unlive_ms);
				pw->stop();
				delete pw;
				iter = thd_vec_.erase(iter);
				continue;
			}
			++iter;
		}
		return true;
	}
private:
	int min_thread_num_;
	int max_thread_num_;
	int max_queue_size_;
	thread_vec thd_vec_;
	std::mutex thread_lock_;
	std::atomic<bool> is_runing_;
	int keep_live_time_ = 1000;
};

volatile std::sig_atomic_t gSignalStatus;
std::atomic_bool is_running(true);
void sig_handler(int sig)
{
	printf("get sig\n");
	is_running = false;
}
int main()
{
	std::signal(SIGINT, sig_handler);
	int min_thread_num = 5;
	int max_thread_num = 10;
	int max_queue_size = 10;
	ThreadPool tp(min_thread_num, max_thread_num, max_queue_size);
	for (int i = 0; i < 999; ++i)
	{
		TaskData t;
		t.data = i;
		tp.push(t);
	}

	while (is_running)
	{
		std::this_thread::sleep_for(std::chrono::milliseconds(10));
	}
	TaskData t;
	t.data = 1;
	tp.push(t);
#ifdef _WIN32
	system("pause");
#endif
}