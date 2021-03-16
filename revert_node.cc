class Solution {
public:
    TreeNode* invertTree(TreeNode* root) {
		if (NULL == root)
		{
			return NULL;
		}

		TreeNode* right_node = root->right;
		root->right = invertTree(root->left);
		root->left = invertTree(right_node);
		return root;
    }
};
