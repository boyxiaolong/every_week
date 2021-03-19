/**
 * Definition for singly-linked list.
 * struct ListNode {
 *     int val;
 *     ListNode *next;
 *     ListNode() : val(0), next(nullptr) {}
 *     ListNode(int x) : val(x), next(nullptr) {}
 *     ListNode(int x, ListNode *next) : val(x), next(next) {}
 * };
 */
class Solution {
public:
    ListNode* mergeKLists(vector<ListNode*>& lists) {
		ListNode* phead = NULL;
		ListNode* pcur = NULL;
		while (true)
		{
			int min = 999999999;
			int min_index = -1;
			for (int i = 0; i < lists.size(); ++i)
			{
				ListNode* curnode = lists[i];
				if (NULL == curnode)
				{
					continue;
				}
				int cur_val = curnode->val;
				if (min > cur_val)
				{
					min = cur_val;
					min_index = i;
				}
			}

			if (min_index < 0)
			{
				if (pcur)
				{
					pcur->next = NULL;
				}
				break;
			}

			ListNode* curnodes = lists[min_index];
			lists[min_index] = curnodes->next;
			if (!phead)
			{
				phead = curnodes;
				pcur = phead;
			}
			else
			{
				pcur->next = curnodes;
				pcur = pcur->next;
			}
		}
		return phead;
    }
};
