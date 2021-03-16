/**
 * Definition for singly-linked list.
 * struct ListNode {
 *     int val;
 *     ListNode *next;
 *     ListNode(int x) : val(x), next(NULL) {}
 * };
 */
class Solution {
public:
    void deleteNode(ListNode* node) {
		if (NULL == node)
		{
			return;
		}

		ListNode* cur = node;
		
		while (cur)
		{
			ListNode* next = cur->next;
			if (next == NULL)
			{
				break;
			}

			cur->val = next->val;
			if (next->next == NULL)
			{
				cur->next = NULL;
				break;
			}
			cur = next;
		}
    }
};
