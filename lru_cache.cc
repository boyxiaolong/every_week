#include <iostream>
#include <map>

class LRUCache {
public:
	LRUCache(int capacity) {
		capacity_ = capacity;
		next_index_ = 0;
	}

	int get(int key) {
		std::map<int, store_tv>::iterator iter = kv_map.find(key);
		if (iter == kv_map.end())
		{
			return -1;
		}

		int index = iter->second.index;
		time_map_.erase(index);
		++next_index_;
		time_map_[next_index_] = key;
		iter->second.index = next_index_;
		return iter->second.value;
	}

	void put(int key, int value) {
		std::map<int, store_tv>::iterator iter = kv_map.find(key);
		//新key
		if (iter == kv_map.end())
		{
			//如果键值超过上限
			if (time_map_.size() >= capacity_)
			{
				int pre_index = time_map_.begin()->first;
				int key = time_map_.begin()->second;
				time_map_.erase(time_map_.begin());
				kv_map.erase(key);
			}
			++next_index_;
			time_map_[next_index_] = key;
			kv_map[key] = store_tv(next_index_, value);
			return;
		}

		int pre_index = iter->second.index;
		time_map_.erase(pre_index);

		++next_index_;
		iter->second.index = next_index_;
		iter->second.value = value;
		time_map_[next_index_] = key;
	}

	int capacity_;
	int next_index_;
	struct store_tv
	{
		int index;
		int value;
		store_tv() :index(0), value(0){}
		store_tv(int i, int v) :index(i), value(v){

		}
	};
	std::map<int, store_tv> kv_map;
	std::map<int, int> time_map_;
};

/**
* Your LRUCache object will be instantiated and called as such:
* LRUCache* obj = new LRUCache(capacity);
* int param_1 = obj->get(key);
* obj->put(key,value);
*/

int main()
{
	int maxlen = 10;
	LRUCache* obj = new LRUCache(maxlen);
	int param_1 = obj->get(1);
	for (int i = 0; i < maxlen * 2; ++i){
		obj->put(i, i * 2);
	}
	for (int i = 0; i < maxlen * 2; ++i){
		obj->get(i);
	}
	delete obj;
}
