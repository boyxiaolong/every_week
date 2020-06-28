#include <iostream>

struct ListNode {
	int val;
	struct ListNode *next;
	ListNode(int val) :val(val), next(NULL) {}
};

class Solution {
public:
	/**
	 *
	 * @param l1 ListNode类
	 * @param l2 ListNode类
	 * @return ListNode类
	 */
	ListNode* addTwoNumbers(ListNode* l1, ListNode* l2) {
		// write code here
		if (l1 == NULL)
		{
			return l2;
		}

		if (l2 == NULL)
		{
			return l1;
		}
		ListNode* sumlist = new ListNode(0);
		ListNode* curNode = sumlist;
		int last_add = 0;
		bool is_first = true;
		while (l1 && l2)
		{
			int cur = l1->val + l2->val + last_add;
			last_add = cur / 10;
			cur -= last_add * 10;

			l1 = l1->next;
			l2 = l2->next;

			if (!is_first)
			{
				curNode->next = new ListNode(0);
				curNode = curNode->next;
			}
			curNode->val = cur;
			is_first = false;
		}

		while (l1)
		{
			int cur = l1->val + last_add;
			last_add = cur / 10;
			cur -= last_add * 10;
			curNode->next = new ListNode(0);
			curNode = curNode->next;
			curNode->val = cur;
			l1 = l1->next;
		}
		while (l2)
		{
			int cur = l2->val + last_add;
			last_add = cur / 10;
			cur -= last_add * 10;
			curNode->next = new ListNode(0);
			curNode = curNode->next;
			curNode->val = cur;

			l2 = l2->next;
		}
		if (last_add > 0)
		{
			curNode->next = new ListNode(0);
			curNode = curNode->next;
			curNode->val = last_add;
		}
		return sumlist;
	}
};
int main()
{
	int a1[] = {5};
	ListNode* l1 = new ListNode(0);
	ListNode* cur1 = l1;
	int size1 = sizeof(a1) / sizeof(a1[0]);
	for (size_t i = 0; i < size1; i++)
	{
		cur1->val = a1[i];
		if (i != size1 - 1)
		{
			cur1->next = new ListNode(0);
			cur1 = cur1->next;
		}
	}

	ListNode* l2 = new ListNode(0);
	ListNode* cur2 = l2;
	int a2[] = {5};
	int size2 = sizeof(a2) / sizeof(a2[0]);
	for (size_t i = 0; i < size2; i++)
	{
		cur2->val = a2[i];
		if (i != size2 - 1)
		{
			cur2->next = new ListNode(0);
			cur2 = cur2->next;
		}
	}

	Solution s;
	auto res = s.addTwoNumbers(l1, l2);
	while (res)
	{
		printf("%d ", res->val);
		res->next = res;
	}
}
