#include <string>
#include <set>
#include <map>
using namespace std;

class Solution {
public:
	int lengthOfLongestSubstring(string s) {
		map<int, int> pre;
		int max = 0;
		for (int i = 0; i < s.size(); ++i)
		{
			int value = s[i] - 'a';
			map<int, int>::iterator iter = pre.find(value);
			if (iter != pre.end())
			{
				//保存最大个数
				if (max < pre.size())
				{
					max = pre.size();
				}
				int pre_index = iter->second;
				iter->second = i;
				for (int j = i - pre.size(); j < pre_index; ++j)
				{
					pre.erase(s[j] - 'a');
				}
			}
			else
			{
				pre[value] = i;
			}
		}
		if (max < pre.size())
		{
			max = pre.size();
		}
		return max;
	}
};

int main()
{
	string str("umvejcuuk");
	Solution s;
	printf("%d\n", s.lengthOfLongestSubstring(str));
}
