package db

import (
	"database/sql"
	"fmt"
	"behaviorlog/loadconfig"
)

var OperateSqlInst *OperateSql

func init(){
	OperateSqlInst = &OperateSql{
		insertsqls:make(chan string, 1000),
	}

	OperateSqlInst.OpenMysql()
	go OperateSqlInst.OperateSql()
}

type OperateSql struct {
	Db *sql.DB
	insertsqls chan string
}

func(m* OperateSql) OpenMysql(){
	var err error
	m.Db, err = sql.Open("mysql",
	loadconfig.SqlAddr)

	if err != nil {
		fmt.Println("open mysql error", err)
		return
	}

	err = m.Db.Ping()

	if err != nil {
		fmt.Println("ping mysql error", err)
		return
	}

	fmt.Println("ping mysql success")
}

func (m *OperateSql) AddSend(sql string) error {
	select {
	case m.insertsqls <- sql:
	}

	return nil
}

func (m *OperateSql) OperateSql() {
	for {
		select {
		case msg := <-m.insertsqls:
			m.InsertSql(msg)
		}
	}
}

func (m *OperateSql) InsertSql(sql string) {
	_, err := m.Db.Exec(sql)

	if err != nil {
		fmt.Println("Insert table error", err)
		fmt.Println("Insert table error", sql)
	}
}