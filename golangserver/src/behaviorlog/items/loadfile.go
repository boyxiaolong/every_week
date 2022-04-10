package items

import (
	"fmt"
	"behaviorlog/operat_file"
)



func GetItemsDirectory() {
	logDirectory := operat_file.GetCurrentDirectory()
	itemsDirectory := logDirectory + "/items"

	filenames,err  := operat_file.GetAllFileName(itemsDirectory)

	if err != nil {
		fmt.Printf("get items file error %d \n", err)
	}

	for _, v := range filenames  {
		xmlData := LoadFile(v)

		if xmlData == nil {
			fmt.Printf("loadfile fail filename %s \n", v)
			return
		}

		AddItemMgr(xmlData)
	}
}

