#include <iostream>
#include <vector>
#include <algorithm>

using namespace std;
class Solution {
public:
	int find_index(vector<int>& nums, int begin, int end, int value)
	{
		if (begin == end) {
			return begin;
		}
		int mid = (begin + end) / 2;
		int mid_value = nums[mid];
		if (mid_value == value){
			return mid;
		}
		if (mid_value < value) {
			return find_index(nums, mid + 1, end, value);
		}

		return find_index(nums, begin, mid - 1, value);
	}
	vector<int> twoSum(vector<int>& nums, int target) {
		vector<int> res;
		res.resize(2, -1);
		if (nums.empty())
		{
			return res;
		}

		auto new_vec = nums;
		std::sort(new_vec.begin(), new_vec.end());
		for (int i = 0; i < new_vec.size(); ++i)
		{
			int cur_num = new_vec[i];
			int left = target - cur_num;
			//printf("cur i %d left %d\n", i, left);
			int index = find_index(new_vec, i + 1, new_vec.size() - 1, left);
			//printf("find index %d value %d\n", index, nums[index]);
			if (new_vec[index] == left)
			{
				int sum = 0;
				for (int j = 0; j < nums.size(); ++j)
				{
						if (nums[j] == new_vec[i] && res[0] == -1)
						{
							res[0] = j;
							++sum;
						}
						else if (nums[j] == new_vec[index] && res[1] == -1)
						{
							res[1] = j;
							++sum;
						}
						if (sum == 2)
						{
							break;
						}
				}
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
	nums.push_back(3);
	Solution s;
	s.twoSum(nums, 6);
}
