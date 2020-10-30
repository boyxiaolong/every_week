class Solution {
public:
    char firstUniqChar(string s) {
		if (s.empty())
		{
			return ' ';
		}
#define max_size 60
		int num[max_size] = { 0 };
		int pos[max_size] = { 0 };
		int min_char = 2147483647;
		bool is_find = false;
		for (int i = 0; i < s.size(); ++i)
		{
			int t = s[i] - 'a';
			++num[t];
			if (num[t] == 1)
			{
				pos[t] = i+1;
			}
			else if (num[t] > 1)
			{
				pos[t] = 0;
			}
		}

		
		for (int i = 0; i < max_size; ++i)
		{
			int index = pos[i];
			if (index > 0 && min_char > index)
			{
				min_char = index;
				is_find = true;
			}
		}

		if (!is_find)
		{
			return ' ';
		}
		return s[min_char-1];
    }
};