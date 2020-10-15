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
class WorkerThread
{
public:
	WorkerThread(int thread_id) :thread_id_(thread_id)
	{

	}
	virtual ~WorkerThread()
	{
		printf("~WorkerThread\n");
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
		thread_runing_ = false;
	}

	bool is_full()
	{
		return task_size_ >= max_queue_size_;
	}
private:
	std::queue<TaskData> tasks_;
	volatile bool thread_runing_ = true;
	std::mutex task_lock_;
	std::mutex thread_lock_;
	std::condition_variable con_;
	int max_queue_size_ = 10;
	std::thread* thd_ = NULL;
	int thread_id_;
	//todo
	volatile int task_size_ = 0;
};

typedef WorkerThread* pWorkerThread;

class ThreadPool
{
public:
	typedef std::list<pWorkerThread> thread_vec;
	ThreadPool()
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
			return false;
		}
		check_min_threads();
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
			WorkerThread* pw = new WorkerThread(thd_id);
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
	}

private:
	void check_min_threads()
	{
		std::unique_lock<std::mutex> guard(thread_lock_);
		int thd_size = thd_vec_.size();
		if (thd_size < min_thread_num_)
		{
			for (int i = thd_size + 1; i <= min_thread_num_; ++i)
			{
				WorkerThread* pw = new WorkerThread(i);
				pw->start();
				thd_vec_.push_back(pw);
				printf("create thread %d\n", i);
			}
		}
	}

private:
	int min_thread_num_ = 1;
	int max_thread_num_ = 10;
	thread_vec thd_vec_;
	std::mutex thread_lock_;
	volatile bool is_runing_ = true;
};

volatile std::sig_atomic_t gSignalStatus;
std::atomic_bool is_running(true);
void sig_handler(int sig)
{
	is_running = false;
}
int main()
{
	std::signal(SIGINT, sig_handler);
	ThreadPool tp;
	for (int i = 0; i < 99; ++i)
	{
		TaskData t;
		t.data = i;
		tp.push(t);
	}

	while (is_running)
	{
		std::this_thread::sleep_for(std::chrono::milliseconds(10));
	}

	system("pause");
}