#include<iostream>
#include<cstdlib>
#include<string>
#include<cstdio>
using namespace std;
const int T_S = 200;
class HashTableEntry {
public:
	int k;
	int v;
	HashTableEntry(int k, int v) {
		this->k = k;
		this->v = v;
		pnext = NULL;
	}
	HashTableEntry* pnext;
};
class HashMapTable {
private:
	HashTableEntry **t;
public:
	HashMapTable() {
		t = new HashTableEntry *[T_S];
		for (int i = 0; i < T_S; i++) {
			t[i] = NULL;
		}
	}
	int HashFunc(int k) {
		return k % T_S;
	}
	void Insert(int k, int v) {
		int h = HashFunc(k);
		HashTableEntry* porigin = t[h];
		if (porigin == NULL) {
			porigin = new HashTableEntry(k,v);
			t[h] = porigin;
			return;
		}
		if (porigin->k == k) {
			porigin->v = v;
			printf("新建 k %d v %d\n", k, v);
			return;
		}
		HashTableEntry* pnew = new HashTableEntry(k, v);
		HashTableEntry* ppre = porigin;
		HashTableEntry* pnext = porigin;
		while (pnext)
		{
			ppre = pnext;
			pnext = ppre->pnext;
		}
		ppre->pnext = pnew;
		printf("碰撞 k %d v %d\n", k, v);
	}
	int SearchKey(int k) {
		int h = HashFunc(k);
		HashTableEntry* porigin = t[h];
		while (porigin)
		{
			if (porigin->k == k)
			{
				printf("find k %d v %d\n", k,porigin->v);
				return porigin->v;
			}
			porigin = porigin->pnext;
		}
		printf("not find key:%d\n", k);
		return 0;
	}
	void Remove(int k) {
		int h = HashFunc(k);
		HashTableEntry* ppre = t[h];
		if (NULL == ppre)
		{
			printf("rm k %d failed \n", k);
			return;
		}

		if (ppre->k == k)
		{
			t[h] = ppre->pnext;
			delete ppre;
			printf("rm k %d\n", k);
			return;
		}

		HashTableEntry* pnext = ppre->pnext;
		while (pnext)
		{
			if (pnext->k == k)
			{
				ppre->pnext = pnext->pnext;
				printf("rm k %d\n", k);
				delete pnext;
				return;
			}
			ppre = pnext;
			pnext = pnext->pnext;
		}
		printf("rm k %d failed \n", k);
	}
	~HashMapTable() {
		for (int i = 0; i < T_S; i++) {
			HashTableEntry* pcur = t[i];
			HashTableEntry* pnext = pcur;
			while (pnext)
			{
				delete pcur;
				pcur = pnext;
				pnext = pcur->pnext;
			}
		}
	}
};
int main() {
	HashMapTable hash;
	for (int i = 0; i < 10000; ++i)
	{
		hash.Insert(i, i + 1);
	}

	for (int i = 0; i < 10005; ++i)
	{
		hash.Remove(i);
	}
	for (int i = 0; i < 10005; ++i)
	{
		int v = hash.SearchKey(i);
	}
	return 0;
}
