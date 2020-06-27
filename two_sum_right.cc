#include <iostream>
#include <vector>
#include <algorithm>
#include <map>

using namespace std;
class Solution {
public:
	vector<int> twoSum(vector<int>& nums, int target) {
		vector<int> res;
		res.resize(2, -1);
		if (nums.empty())
		{
			return res;
		}

		map<int, int> store_map;
		for (int i = 1; i < nums.size(); ++i)
		{
			int value = nums[i];
			map<int, int>::iterator iter = store_map.find(value);
			if (iter != store_map.end())
			{
				if (value * 2 == target)
				{
					res[0] = iter->second;
					res[1] = i;
					return res;
				}
			}
			else
			{
				store_map.insert(std::make_pair(value, i));
			}
		}

		for (int i = 0; i < nums.size(); ++i)
		{
			int value = nums[i];
			int left = target - value;
			map<int, int>::iterator iter = store_map.find(left);
			if (iter != store_map.end())
			{
				int pre = iter->second;
				if (pre == i)
				{
					continue;
				}
				res[0] = i;
				res[1] = pre;
				return res;
			}
		}
		return res;
	}
};

int main()
{
	vector<int> nums;
	int arr[] = { 1,3,4,2 };
	for (int i = 0; i < sizeof(arr) / sizeof(arr[0]); ++i)
	{
		nums.push_back(arr[i]);
	}
	Solution s;
	s.twoSum(nums, 6);
}
