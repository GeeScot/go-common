package sqlio

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"database/sql"
)

type Object interface {
	Map(*sql.Rows) (interface{}, error)
	GetID() int64
	SetID(id int64) interface{}
}

type SQLio struct {
	obj *interface{}
}

var db *sql.DB

func Open(driverName string, dataSourceName string) {
	conn, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		panic(err)
	}

	db = conn
}

func (sqlio *SQLio) CanExecute() (bool, error) {
	if db == nil {
		return false, errors.New("DB is nil or not connected")
	}

	if err := db.Ping(); err != nil {
		return false, err
	}

	return true, nil
}

func (sqlio *SQLio) tableName() string {
	objType := reflect.TypeOf(*sqlio.obj)
	str := strings.ToLower(objType.String())
	startIndex := strings.Index(str, ".")
	return str[startIndex+1:]
}

func (sqlio *SQLio) Get() interface{} {
	return *sqlio.obj
}

func New(obj interface{}) *SQLio {
	sqlio := SQLio{}
	sqlio.obj = &obj

	return &sqlio
}

func (sqlio *SQLio) GetID() int64 {
	obj := (*sqlio.obj).(Object)

	return obj.GetID()
}

func (sqlio *SQLio) SetID(id int64) {
	obj := (*sqlio.obj).(Object)
	objInterface := obj.SetID(id)
	sqlio.obj = &objInterface
}

func (sqlio *SQLio) Select(id int64) error {
	if _, err := sqlio.CanExecute(); err != nil {
		return err
	}

	table := sqlio.tableName()
	query := fmt.Sprintf("SELECT * FROM %s WHERE id=? LIMIT 1", table)

	rows, err := db.Query(query, id)
	if err != nil {
		return err
	}

	defer rows.Close()

	if !rows.Next() {
		return errors.New("no results")
	}

	obj := (*sqlio.obj).(Object)
	retVal, err := obj.Map(rows)
	if err != nil {
		return err
	}

	sqlio.obj = &retVal

	return nil
}

func (sqlio *SQLio) SelectWhere(values Values) error {
	if _, err := sqlio.CanExecute(); err != nil {
		return err
	}

	table := sqlio.tableName()
	where, params := sqlio.whereClause(values)

	query := fmt.Sprintf("SELECT * FROM %s", table)
	query += where
	query += " LIMIT 1"

	rows, err := db.Query(query, params...)
	if err != nil {
		return err
	}

	defer rows.Close()

	if !rows.Next() {
		return errors.New(fmt.Sprintf("no results from select query %v", values))
	}

	obj := (*sqlio.obj).(Object)
	retVal, err := obj.Map(rows)
	if err != nil {
		return err
	}

	sqlio.obj = &retVal

	return nil
}

func (sqlio SQLio) Exists(values Values) bool {
	if _, err := sqlio.CanExecute(); err != nil {
		return false
	}

	table := sqlio.tableName()
	where, params := sqlio.whereClause(values)

	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", table)
	query += where
	query += " LIMIT 1"

	rows, err := db.Query(query, params...)
	if err != nil {
		return false
	}

	defer rows.Close()

	if !rows.Next() {
		return false
	}

	var count int64
	rows.Scan(&count)

	return count > 0
}

func (sqlio *SQLio) Save() error {
	if _, err := sqlio.CanExecute(); err != nil {
		return err
	}

	var data []interface{}
	values := ToMap(sqlio.obj)

	columns := ""
	inputValues := ""
	index := 0

	keys := values.SortedKeys()
	for _, key := range keys {
		if key == "id" {
			continue
		}

		columns += key + ","
		inputValues += "?,"

		data = append(data, values[key])
		index++
	}

	columns = string(columns[:len(columns)-1])
	inputValues = string(inputValues[:len(inputValues)-1])

	table := sqlio.tableName()
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table, columns, inputValues)

	stmnt, err := db.Prepare(query)
	if err != nil {
		return err
	}

	result, err := stmnt.Exec(data...)
	if err != nil {
		return err
	}

	if _, err := result.RowsAffected(); err != nil {
		return err
	}

	newID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	sqlio.SetID(newID)

	return nil
}

func (sqlio *SQLio) Update(newObj interface{}) error {
	if _, err := sqlio.CanExecute(); err != nil {
		return err
	}

	var data []interface{}
	values := ToMap(&newObj)

	sets := ""
	index := 0

	keys := values.SortedKeys()
	for _, key := range keys {
		if key == "id" {
			continue
		}

		sets += key + "=?,"

		data = append(data, values[key])
		index++
	}

	id := sqlio.GetID()

	data = append(data, id)
	strLength := len(sets) - 1

	table := sqlio.tableName()
	query := fmt.Sprintf("UPDATE %s SET %s WHERE id=?", table, sets[:strLength])

	stmnt, err := db.Prepare(query)
	if err != nil {
		return err
	}

	result, err := stmnt.Exec(data...)
	if err != nil {
		return err
	}

	if _, err := result.RowsAffected(); err != nil {
		return err
	}

	sqlio.obj = &newObj
	sqlio.SetID(id)

	return nil
}

func (sqlio *SQLio) Delete() error {
	if _, err := sqlio.CanExecute(); err != nil {
		return err
	}

	table := sqlio.tableName()
	query := fmt.Sprintf("DELETE FROM %s WHERE id=?", table)

	stmnt, err := db.Prepare(query)
	if err != nil {
		return err
	}

	id := sqlio.GetID()

	if id < 0 {
		return errors.New("sql.delete: cannot delete what is not there")
	}

	result, err := stmnt.Exec(id)
	if err != nil {
		return err
	}

	if _, err = result.RowsAffected(); err != nil {
		return err
	}

	return nil
}

func (sqlio *SQLio) whereClause(values Values) (string, []interface{}) {
	params := make([]interface{}, len(values))
	index := 0

	where := " WHERE"

	keys := values.SortedKeys()
	for _, key := range keys {
		where += " " + key + "=? AND"
		params[index] = values[key]

		index++
	}

	where = where[:len(where)-4]

	return where, params
}
