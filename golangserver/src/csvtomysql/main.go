package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"

	//"reflect"
	"runtime/debug"
	"strings"

	_ "public/db"

	"github.com/axgle/mahonia"
	"github.com/golang/protobuf/proto"
)

type (
	// ColumnInfo aa
	ColumnInfo struct {
		ColumnName, ColumnType string
		Valid                  bool
	}
)

var (
//validColumn = []string{"id", "pcid", "sub_goods_id", "comment", "next_goods", "name", "type", "subtype", "duration", "ad_pic", "level", "start_time", "end_time", "limit_time", "currency_type", "original_price", "price", "label", "desc", "goods", "icon", "gift", "discontinue"}
)

//"sub_goods_id","comment","next_goods","name","type","subtype","duration","ad_pic","level","start_time","end_time","limit_time","currency_type","original_price","price","label","desc","goods","icon","gift","discontinue"}

func makeTableSQL(tableName string, tableColumns []ColumnInfo) string {
	sql := fmt.Sprintf("DROP TABLE IF EXISTS `%v`; \nCREATE TABLE `%v` (\n", tableName, tableName)
	for _, v := range tableColumns {
		if !v.Valid {
			continue
		}
		column := ""
		if v.ColumnType == "uint" {
			column = fmt.Sprintf("\t`%v` bigint(20) unsigned NOT NULL,\n", v.ColumnName)
		} else if v.ColumnType == "string" {
			column = fmt.Sprintf("\t`%v` text COLLATE utf8mb4_unicode_ci,\n", v.ColumnName)
		}
		sql += column
	}

	sql = sql + "\tPRIMARY KEY (`id`) \n) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;\n"
	return sql
}

func makeDataSQL(tableName string, tableColumns []ColumnInfo, columnsData [][]string, delOldData bool) []string {
	dec := mahonia.NewDecoder("GB18030")
	sqlList := make([]string, 0)
	idset := make(map[uint64]bool)
	if delOldData {
		sql := fmt.Sprintf("DELETE FROM %v;\n", tableName)
		sqlList = append(sqlList, sql)
	}
	for _, v := range columnsData {
		insertValues := ""
		updateValues := ""
		for i, columninfo := range tableColumns {
			if !columninfo.Valid {
				continue
			}

			if i == 0 {
				insertValues = insertValues + fmt.Sprintf("'%v'", dec.ConvertString(v[i]))
				updateValues = updateValues + fmt.Sprintf("`%v` = '%v'", columninfo.ColumnName, v[i])
				id, err := strconv.ParseUint(v[i], 10, 64)
				if err != nil {
					fmt.Println("init id [" + v[i] + "] error!!")
					os.Exit(3)
				}
				_, ok := idset[id]
				if ok {
					fmt.Printf("repead id %v", id)
					os.Exit(3)
				}
				idset[id] = true
			} else {
				insertValues = insertValues + fmt.Sprintf(",'%v'", dec.ConvertString(v[i]))
				updateValues = updateValues + fmt.Sprintf(",`%v` = '%v'", columninfo.ColumnName, dec.ConvertString(v[i]))
			}
		}
		sql := ""
		if delOldData {
			sql = fmt.Sprintf("INSERT INTO %v VALUES(%v);\n", tableName, insertValues)
		} else {
			sql = fmt.Sprintf("INSERT INTO %v VALUES(%v) ON DUPLICATE KEY UPDATE %v;\n", tableName, insertValues, updateValues)
		}
		sqlList = append(sqlList, sql)
	}
	return sqlList
}

func makeTableProto(tableName string, tableColumns []ColumnInfo) string {
	protoData := fmt.Sprintf("syntax = \"proto3\"; \n\npackage db; \nmessage %v\n{\n", tableName)

	protoIndex := 1
	for _, v := range tableColumns {
		if !v.Valid {
			continue
		}
		column := ""
		if v.ColumnType == "uint" {
			column = fmt.Sprintf("\tuint64 %v = %v;\n", v.ColumnName, protoIndex)
		} else if v.ColumnType == "string" {
			column = fmt.Sprintf("\tstring %v = %v;\n", v.ColumnName, protoIndex)
		}
		protoData += column
		protoIndex = protoIndex + 1
	}

	protoData = protoData + fmt.Sprintf("}\n\nmessage %v_set\n{ \n\trepeated %v set = 1;\n}", tableName, tableName)
	return protoData
}

func makeCSVValidColumn(tableName string) ([]string, bool) {
	validColumn := make([]string, 0)
	msgName := "db." + tableName

	msgReflect := proto.MessageType(msgName)
	if msgReflect == nil {
		return nil, false
	}

	properties := proto.GetProperties(msgReflect.Elem())
	if properties == nil {
		return nil, false
	}
	for _, v := range properties.Prop {
		validColumn = append(validColumn, v.OrigName)
	}
	return validColumn, true
}

func makeCSVSqlTable(file string, delOldData bool) {
	tempfilename := filepath.Base(file)
	fileName := strings.Replace(tempfilename, path.Ext(tempfilename), "", 1)
	tableName := fileName + "_config"
	tableColumns := make([]ColumnInfo, 0)
	columnsName := make([]string, 0)
	columnsType := make([]string, 0)
	columnsData := make([][]string, 0)

	validColumn, ok := makeCSVValidColumn(tableName)
	if !ok {
		panic(fmt.Sprintf("make csv %v valid column error", file))
	}

	csvFile, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()

	r := csv.NewReader(csvFile)

	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	for index, record := range records {
		if index == 1 {
			for _, columnType := range record {
				if columnType != "uint" {
					columnType = "string"
				}
				columnsType = append(columnsType, columnType)
			}
		} else if index == 2 {
			for _, columnName := range record {
				if strings.Contains(columnName, "#") {
					columnName = "id"
				}
				columnsName = append(columnsName, columnName)
			}
		} else if index >= 3 {
			columnsData = append(columnsData, record)
		}
	}

	for index := 0; index < len(columnsName); index++ {
		var column ColumnInfo
		column.ColumnName = columnsName[index]
		column.ColumnType = columnsType[index]
		column.Valid = false
		for _, v := range validColumn {
			if v == column.ColumnName {
				column.Valid = true
			}
		}
		tableColumns = append(tableColumns, column)
	}

	tableSQL := makeTableSQL(tableName, tableColumns)
	dataSQLList := makeDataSQL(tableName, tableColumns, columnsData, delOldData)

	tableFileName := tableName + ".sql"

	tableFile, err := os.Create(tableFileName)
	if err != nil {
		panic(err)
	}
	defer tableFile.Close()

	tableFile.WriteString(tableSQL)

	recordFileName := tableName + "_record" + ".sql"
	recordFile, err := os.Create(recordFileName)
	if err != nil {
		panic(err)
	}
	defer recordFile.Close()

	for _, data := range dataSQLList {
		recordFile.WriteString(data)
	}

	// tableProto := makeTableProto(tableName, tableColumns)

	// protoFileName := tableName + ".proto"

	// protoFile, err := os.Create(protoFileName)
	// if err != nil {
	// 	panic(err)
	// }
	// defer protoFile.Close()

	// protoFile.WriteString(tableProto)

	fmt.Println("read line ", len(dataSQLList))
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
		}
	}()

	fmt.Println("parse file : " + os.Args[1])
	delOldData, error := strconv.ParseBool(os.Args[2])
	if error != nil {
		fmt.Println("param error")
		os.Exit(3)
	}
	makeCSVSqlTable(os.Args[1], delOldData)
	fmt.Println("parse end")
}
