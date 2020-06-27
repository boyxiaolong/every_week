#include <iostream>
#include <vector>
#include <algorithm>

using namespace std;
class Solution {
public:
	int find_index(vector<int>& nums, int begin, int end, int value)
	{
		if (begin == end) {
			//printf("equal %d\n", begin);
			return begin;
		}
		int mid = (begin + end) / 2;
		int mid_value = nums[mid];
		if (mid_value == value){
			//printf("mid %d\n", mid);
			return mid;
		}
		if (mid_value < value) {
			return find_index(nums, mid + 1, end, value);
		}

		return find_index(nums, begin, mid - 1, value);
	}
	vector<int> twoSum(vector<int>& nums, int target) {
		vector<int> res;
		if (nums.empty())
		{
			return res;
		}

		std::sort(nums.begin(), nums.end());
		for (int i = 0; i < nums.size(); ++i)
		{
			int cur_num = nums[i];
			int left = target - cur_num;
			//printf("cur i %d left %d\n", i, left);
			int index = find_index(nums, i + 1, nums.size() - 1, left);
			//printf("find index %d value %d\n", index, nums[index]);
			if (nums[index] == left)
			{
				res.push_back(i);
				res.push_back(index);
				break;
			}
		}
		return res;
	}
};

int main()
{
	vector<int> nums;
	nums.push_back(3);
	nums.push_back(2);
	nums.push_back(4);
	Solution s;
	s.twoSum(nums, 6);
}
