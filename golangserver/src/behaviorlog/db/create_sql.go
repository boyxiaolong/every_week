package db

import (
	"behaviorlog/items"
	"database/sql"
	"fmt"
	"public/common"

	_ "github.com/go-sql-driver/mysql"
)

func init() {

}

func CreateObjectTable(db *sql.DB, table_name string) error {
	sql := "CREATE TABLE " +
		table_name +
		" (`id` bigint(20) NOT NULL AUTO_INCREMENT," +
		"`object_type` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL," +
		"`createtime` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL," +
		"`player_id` bigint(20) NOT NULL," +
		"`player_name` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL," +
		"`sword` bigint(20) UNSIGNED NOT NULL," +
		"`townlevel` int(11) UNSIGNED NOT NULL," +
		"`kingdom_id` int(11) UNSIGNED NOT NULL," +
		"`type` int(11) UNSIGNED NOT NULL," +
		"`type_name` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL," +
		"`change_type` int(11) UNSIGNED NOT NULL," +
		"`change_num` bigint(20) UNSIGNED NOT NULL," +
		"`left_num` bigint(20) UNSIGNED NOT NULL," +
		"`action_id` int(11) UNSIGNED NOT NULL," +
		"`action_name` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL," +
		"`para_int_1` bigint(20) UNSIGNED NOT NULL," +
		"`para_int_2` bigint(20) UNSIGNED NOT NULL," +
		"`para_int_3` bigint(20) UNSIGNED NOT NULL," +
		"`para_int_4` bigint(20) UNSIGNED NOT NULL," +
		"`para_int_5` bigint(20) UNSIGNED NOT NULL," +
		"`para_str_1` text COLLATE utf8mb4_unicode_ci NOT NULL," +
		"`para_str_2` text COLLATE utf8mb4_unicode_ci NOT NULL," +
		"`para_str_3` text COLLATE utf8mb4_unicode_ci NOT NULL," +
		"`para_str_4` text COLLATE utf8mb4_unicode_ci NOT NULL," +
		"`para_str_5` text COLLATE utf8mb4_unicode_ci NOT NULL," +
		"PRIMARY KEY (`id`)) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;"

	_, err := db.Exec(sql)

	if err != nil {
		return err
	}

	fmt.Println("create table success, tablename:", table_name)

	return nil
}

func getSqlString(value string) string {
	return "'" + value + "'"
}

func InsertObjectTable(tablename string, info *common.StringParse) string {
	sql := "Insert Into " + tablename +
		" (object_type,createtime,player_id,player_name,sword,townlevel,kingdom_id,type,type_name,change_type,change_num,left_num,action_id,action_name,para_int_1,para_int_2,para_int_3,para_int_4,para_int_5,para_str_1,para_str_2,para_str_3,para_str_4,para_str_5)" +
		" values (" +
		getSqlString(info.GetString(0)) + "," +
		getSqlString(info.GetString(1)) + "," +
		info.GetString(2) + "," +
		getSqlString(info.GetString(3)) + "," +
		info.GetString(4) + "," +
		info.GetString(5) + "," +
		info.GetString(6) + "," +
		info.GetString(7) + "," +
		getSqlString(items.GetObjectItem(info.GetString(0), info.GetString(7))) + "," +
		info.GetString(8) + "," +
		info.GetString(9) + "," +
		info.GetString(10) + "," +
		info.GetString(11) + "," +
		getSqlString(items.GetObjectItem("action", info.GetString(11))) + "," +
		info.GetString(12) + "," +
		info.GetString(13) + "," +
		info.GetString(14) + "," +
		info.GetString(15) + "," +
		info.GetString(16) + "," +
		getSqlString(info.GetString(17)) + "," +
		getSqlString(info.GetString(18)) + "," +
		getSqlString(info.GetString(19)) + "," +
		getSqlString(info.GetString(20)) + "," +
		getSqlString(info.GetString(21)) + ");"

	return sql
}

func CreateEventTable(db *sql.DB, table_name string) error {
	sql := "CREATE TABLE " +
		table_name +
		" (`id` bigint(20) NOT NULL AUTO_INCREMENT," +
		"`event_type` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL," +
		"`createtime` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL," +
		"`player_id` bigint(20) UNSIGNED NOT NULL," +
		"`player_name` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL," +
		"`sword` bigint(20) UNSIGNED NOT NULL," +
		"`townlevel` int(11) UNSIGNED NOT NULL," +
		"`kingdom_id` int(11) UNSIGNED NOT NULL," +
		"`action_id` int(11) UNSIGNED NOT NULL," +
		"`action_name` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL," +
		"`para_int_1` bigint(20) UNSIGNED NOT NULL," +
		"`para_int_2` bigint(20) UNSIGNED NOT NULL," +
		"`para_int_3` bigint(20) UNSIGNED NOT NULL," +
		"`para_int_4` bigint(20) UNSIGNED NOT NULL," +
		"`para_int_5` bigint(20) UNSIGNED NOT NULL," +
		"`para_str_1` text COLLATE utf8mb4_unicode_ci NOT NULL," +
		"`para_str_2` text COLLATE utf8mb4_unicode_ci NOT NULL," +
		"`para_str_3` text COLLATE utf8mb4_unicode_ci NOT NULL," +
		"`para_str_4` text COLLATE utf8mb4_unicode_ci NOT NULL," +
		"`para_str_5` text COLLATE utf8mb4_unicode_ci NOT NULL," +
		"PRIMARY KEY (`id`)) ENGINE=MyISAM DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;"

	_, err := db.Exec(sql)

	if err != nil {
		return err
	}

	fmt.Println("create table success")

	return nil
}

func InsertEventTable(tablename string, info *common.StringParse) string {
	sql := "Insert Into " + tablename +
		" (event_type,createtime,player_id,player_name,sword,townlevel,kingdom_id,action_id,action_name,para_int_1,para_int_2,para_int_3,para_int_4,para_int_5,para_str_1,para_str_2,para_str_3,para_str_4,para_str_5)" +
		" values (" +
		getSqlString(info.GetString(0)) + "," +
		getSqlString(info.GetString(1)) + "," +
		info.GetString(2) + "," +
		getSqlString(info.GetString(3)) + "," +
		info.GetString(4) + "," +
		info.GetString(5) + "," +
		info.GetString(6) + "," +
		info.GetString(7) + "," +
		getSqlString(items.GetEventItem(info.GetString(0), info.GetString(7))) + "," +
		info.GetString(8) + "," +
		info.GetString(9) + "," +
		info.GetString(10) + "," +
		info.GetString(11) + "," +
		info.GetString(12) + "," +
		getSqlString(info.GetString(13)) + "," +
		getSqlString(info.GetString(14)) + "," +
		getSqlString(info.GetString(15)) + "," +
		getSqlString(info.GetString(16)) + "," +
		getSqlString(info.GetString(17)) + ");"

	return sql
}
