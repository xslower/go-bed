package orm

import (
	`database/sql`
	// `fmt`

	_ "github.com/go-sql-driver/mysql"
	// "log"
)

var (
	gDbConn     IConn
	gDbRegistry = make(map[string]func(config map[string]string) IConn, 1)
	gDbType     = map[uint8]string{1: `tinyint`, 2: `smallint`, 3: `int`, 4: `bigint`, 5: `utinyint`, 6: `usmallint`, 7: `uint`, 8: `ubigint`, 9: `char`, 10: `varchar`, 11: `timestamp`, 12: `enum`}
	gDbTypeR    = map[string]uint8{}
)

func initVariable() {
	for i, typ := range gDbType {
		gDbTypeR[typ] = i
	}
}

func init() {
	initVariable()

}

func getColumnsInfo(table, db string) {

}

type IConn interface {
	Query(sql string, args ...interface{}) (columns []string, rrows []RawRow)
	Exec(sql string, args ...interface{}) int
	Prepare(query string) *sql.Stmt
	Begin() error
	Rollback() error
	Commit() error
	Close() error
	//SetMaxIdleConns/SetMaxOpenConns
}

type RawRow [][]byte

type Conn struct {
	db *sql.DB
	tx *sql.Tx
}

func (this *Conn) Query(query string, args ...interface{}) (columns []string, rrows []RawRow) {
	rows, err := this.db.Query(query, args...)
	throw(err)
	//var rrows []RawRow
	columns, _ = rows.Columns()
	throw(err)
	num := len(columns)
	for rows.Next() {
		rr := make([][]byte, num)
		ref := make([]interface{}, num)
		for i, _ := range ref {
			ref[i] = &rr[i]
		}
		err = rows.Scan(ref...)
		throw(err)
		rrows = append(rrows, rr)
	}
	rows.Close()
	return
}

func (this *Conn) Exec(query string, args ...interface{}) int {
	result, err := this.db.Exec(query, args...)
	throw(err)
	id, _ := result.LastInsertId()
	if id == 0 {
		id, _ = result.RowsAffected()
	}
	return int(id)
}

func (this *Conn) Begin() error {
	var err error
	this.tx, err = this.db.Begin()
	return err
}

func (this *Conn) Close() error {
	return this.db.Close()
}

func (this *Conn) Rollback() error {
	return this.tx.Rollback()
}

func (this *Conn) Commit() error {
	return this.tx.Commit()
}

func (this *Conn) Prepare(query string) *sql.Stmt {
	stmt, err := this.db.Prepare(query)
	throw(err)
	return stmt
}
