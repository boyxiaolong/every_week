#include <vector>
#include <bitset>
#include <limits>
using namespace std;

class BitMap
{
public:
	BitMap(int range) :bits_(range, 0)
	{
	}
	void set(int val)
	{
		int index = val >> 5;
		int tmp = val % 32;
		bits_[index] |= 1 << tmp;
	}
	int is_has(int val)
	{
		int index = val >> 5;
		int tmp = val % 32;
		return bits_[index] & 1 << tmp;
	}
private:
	vector<int> bits_;
};
class Solution {
public:
	bool containsDuplicate(vector<int>& nums) {
		BitMap b(std::numeric_limits<int>::max() >> 5 + 1);
		for (int i = 0; i < nums.size(); ++i)
		{
			int val = nums[i];
			if (b.is_has(val))
			{
				return true;
			}
			b.set(val);
		}
		return false;
	}
};
//只能处理整数的情况！！！
int main()
{
	vector<int> vec = { 1, 2, 34, 2, 34 };
	Solution s;
	bool res = s.containsDuplicate(vec);
	printf("res %d\n", res);
}
