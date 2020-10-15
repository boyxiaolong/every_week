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

class TaskData
{
public:
	int data;
};
class WorkerThread
{
public:
	virtual ~WorkerThread()
	{
		printf("~WorkerThread\n");
		if (thd_)
		{
			thd_->join();
			delete thd_;
		}
	}
	void run()
	{		
		while (is_runing_)
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
				if (temp.empty() && is_runing_)
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
		}
		if (is_notify)
		{
			std::unique_lock<std::mutex> guard(thread_lock_);
			con_.notify_one();
		}
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
		is_runing_ = true;
	}
private:
	std::queue<TaskData> tasks_;
	volatile bool is_runing_(true);
	std::mutex task_lock_;
	std::mutex thread_lock_;
	std::condition_variable con_;
	int max_queue_size_ = 100;
	std::thread* thd_ = NULL;
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
	WorkerThread w;
	w.start();
	for (int i = 0; i < 999;++i)
	{
		TaskData t;
		t.data = i;
		w.push(t);
	}
	
	while (is_running)
	{
		std::this_thread::sleep_for(std::chrono::milliseconds(10));
	}
}