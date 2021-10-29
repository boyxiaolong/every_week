#include <chrono>
#include <cstring>
#include <iostream>
#include <cstdint>
#include <string>
#include <array>

constexpr std::uint64_t SIZE_BYTES = 1073741824; // 1GB
class Timer
{
public:
	Timer()
		: mStart(),
		mStop()
	{
		update();
	}

	void update()
	{
		mStart = std::chrono::high_resolution_clock::now();
		mStop = mStart;
	}

	double elapsedMs()
	{
		mStop = std::chrono::high_resolution_clock::now();
		std::chrono::milliseconds elapsed_ms =
			std::chrono::duration_cast<std::chrono::milliseconds>(mStop - mStart);
		return elapsed_ms.count();
	}

private:
	std::chrono::high_resolution_clock::time_point mStart;
	std::chrono::high_resolution_clock::time_point mStop;
};

std::string formatBytes(std::uint64_t bytes)
{
	static const int num_suffix = 5;
	static const char* suffix[num_suffix] = { "B", "KB", "MB", "GB", "TB" };
	double dbl_s_byte = bytes;
	int i = 0;
	for (; (int)(bytes / 1024.) > 0 && i < num_suffix;
		++i, bytes /= 1024.)
	{
		dbl_s_byte = bytes / 1024.0;
	}

	const int buf_len = 64;
	char buf[buf_len];

	// use snprintf so there is no buffer overrun
	int res = snprintf(buf, buf_len, "%0.2f%s", dbl_s_byte, suffix[i]);

	// snprintf returns number of characters that would have been written if n had
	//       been sufficiently large, not counting the terminating null character.
	//       if an encoding error occurs, a negative number is returned.
	if (res >= 0)
	{
		return std::string(buf);
	}
	return std::string();
}

void doMemmove(void* pDest, const void* pSource, std::size_t sizeBytes)
{
	memmove(pDest, pSource, sizeBytes);
}

int big_memory_test()
{
	std::cout << "big_memory_test\n";
	// big array to use for testing
	char* p_big_array = NULL;

	/////////////
	// malloc 
	{
		Timer timer;

		p_big_array = (char*)malloc(SIZE_BYTES * sizeof(char));
		if (p_big_array == NULL)
		{
			std::cerr << "ERROR: malloc of " << SIZE_BYTES << " returned NULL!"
				<< std::endl;
			return 1;
		}

		std::cout << "malloc for \t" << formatBytes(SIZE_BYTES) << " took "
			<< timer.elapsedMs() << "ms"
			<< std::endl;
	}
	{
		// cleanup
		Timer timer;
		free(p_big_array);
		p_big_array = NULL;
		double elapsed_ms = timer.elapsedMs();
		std::cout << "free for \t" << formatBytes(SIZE_BYTES) << " took "
			<< elapsed_ms << "ms " << std::endl;
	}

	return 0;
}

int small_memory_test()
{
	std::cout << "small_memory_test\n";
	// big array to use for testing
	char* p_big_array = NULL;

	constexpr int devide_num = 100000;
	constexpr int every_memory_size = SIZE_BYTES / devide_num;
	std::array<char*, devide_num> memory_array;
	/////////////
	// malloc 
	{
		Timer timer;
		for (int i = 0; i < devide_num; ++i)
		{
			char* small_ic = (char*)malloc(every_memory_size * sizeof(char));
			memory_array[i] = small_ic;
		}

		std::cout << "malloc for \t" << formatBytes(every_memory_size) << " " << devide_num << " took "
			<< timer.elapsedMs() << "ms"
			<< std::endl;
	}
	{
		// cleanup
		Timer timer;
		for (int i = 0; i < devide_num; ++i)
		{
			free(memory_array[i]);
		}

		std::cout << "free for \t" << formatBytes(every_memory_size) << " " << devide_num << " took "
			<< timer.elapsedMs() << "ms"
			<< std::endl;
	}

	return 0;
}

int main(int argc, char* argv[])
{
	printf("open asan\n");
	
	big_memory_test();

	std::cout << "\n\n------\n\n";
	small_memory_test();

	return 0;
}