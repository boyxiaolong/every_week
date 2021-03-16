	int maxDepth(TreeNode* root) {
		if (NULL == root)
		{
			return 0;
		}

		int total_len = 0;
		int left = maxDepth(root->left);
		int right = maxDepth(root->right);
		return left > right ? left + 1 : right + 1;
	}
