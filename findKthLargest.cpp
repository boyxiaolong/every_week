

#include <iostream>
#include <vector>
using namespace std;
void exchange_val(int& v1, int& v2)
{
	int tmp = v1;
	v1 = v2;
	v2 = tmp;
}
void heapify(vector<int>& vec, int i)
{
	int left = (i << 1) + 1;
	int righ = left + 1;
	if (vec[i] < vec[left])
	{
		exchange_val(vec[i], vec[left]);
	}
	if (vec.size() <= righ)
	{
		return;
	}
	if (vec[i] < vec[righ])
	{
		exchange_val(vec[i], vec[righ]);
	}
}
//构建最大堆
void build_heap(vector<int>& vec, int total)
{
	for (int i = total / 2 - 1; i >= 0; --i)
	{
		heapify(vec, i);
	}
}

class Solution {
public:
	int findKthLargest(vector<int>& nums, int k) {
		int total = nums.size();
		build_heap(nums, total);
		int k_max = nums[0];
		for (int i = total - 1; i >= total - k; --i)
		{
			build_heap(nums, i);
			k_max = nums[0];
			nums[0] = nums[i];
		}
		return k_max;
	}
};
int main()
{
	vector<int> vec = { 3, 2, 3, 1, 2, 4, 5, 5, 6 };
	Solution s;
	int res = s.findKthLargest(vec, 5);
	printf("res %d\n", res);
}
