#include <iostream>
#include <vector>
#include <map>
using namespace std;

class Solution {
public:
	bool containsNearbyDuplicate(vector<int>& nums, int k) {
		map<int, int> v_2_index_map;
		for (int i = 0; i < nums.size(); ++i)
		{
			int val = nums[i];
			map<int, int>::iterator iter = v_2_index_map.find(val);
			if (iter == v_2_index_map.end())
			{
				v_2_index_map.insert(std::make_pair(val, i));
			}
			else
			{
				int pre_index = iter->second;
				if (i - pre_index <= k)
				{
					return true;
				}

				iter->second = i;
			}
		}
		return false;
	}
};

int main()
{
	vector<int> vec = { 1, 2, 3, 1 };
	Solution s;
	bool res = s.containsNearbyDuplicate(vec, 2);
	printf("res %d\n", res);
}
