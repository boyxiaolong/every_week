package items

import (
	"fmt"
	"strings"
)

var  ItemMgrInst *ItemsMgr

type Item struct {
	Name string
	Desc string
}

type Items struct {
	Name string
	Desc string
	ItemMap map[string]Item
}

type ItemsMgr struct {
	items map[string]Items
}


func init()  {
	ItemMgrInst = &ItemsMgr{
		items : make(map[string]Items,0),
	}
}


func AddItemMgr(data *XmlLogData) {
	if data == nil {
		fmt.Printf("xmllogdata is nil")
		return
	}
	
	for _, v := range data.ItemsInfo  {
		items :=Items{}
		items.ItemMap = make(map[string]Item)
		items.Name = strings.ToLower(v.Name)
		items.Desc = v.Desc

		for _, k := range v.ItemInfo  {
			item := Item{}
			item.Name = k.Name
			item.Desc = k.Desc
			items.ItemMap[item.Name] = item
		}

		ItemMgrInst.items[items.Name] = items
		fmt.Printf("add items %s\n", items.Name)
	}
}

func GetObjectItem(items_name string, item_name string) string  {
	items_name = "object_" + items_name
	if items, ok := ItemMgrInst.items[items_name]; ok {
		if item, ok := items.ItemMap[item_name]; ok {
			return item.Desc
		}
	}

	return ""
}

func GetEventItem(items_name string, item_name string) string  {
	items_name =  items_name + "_action"
	items_name = strings.ToLower(items_name)
	if items, ok := ItemMgrInst.items[items_name]; ok {
		if item, ok := items.ItemMap[item_name]; ok {
			return item.Desc
		}
	}

	return ""
}