package activerecord

import (
	"database/sql"
	"errors"
	"fmt"
	"golanger.com/framework/utils"
	"reflect"
	"strconv"
	"strings"
)

var OnDebug = false

type ActiveRecord struct {
	Db              *sql.DB
	TableName       string
	LimitStr        int
	OffsetStr       int
	WhereStr        string
	ParamStr        []interface{}
	OrderStr        string
	ColumnStr       string
	PrimaryKey      string
	JoinStr         string
	GroupByStr      string
	HavingStr       string
	QuoteIdentifier string
	ParamIdentifier string
	ParamIteration  int
}

/**
 * Add New sql.DB in the future i will add ConnectionPool.Get() 
 */
func NewActiveRecord(db *sql.DB, options ...interface{}) ActiveRecord {
	var ar ActiveRecord
	if len(options) == 0 {
		ar = ActiveRecord{Db: db, ColumnStr: "*", PrimaryKey: "Id", QuoteIdentifier: "`", ParamIdentifier: "?", ParamIteration: 1}
	} else if options[0] == "pg" {
		ar = ActiveRecord{Db: db, ColumnStr: "id", PrimaryKey: "id", QuoteIdentifier: "\"", ParamIdentifier: options[0].(string), ParamIteration: 1}
	} else if options[0] == "mssql" {
		ar = ActiveRecord{Db: db, ColumnStr: "id", PrimaryKey: "id", QuoteIdentifier: "", ParamIdentifier: options[0].(string), ParamIteration: 1}
	}

	return ar
}

func (ar *ActiveRecord) SetTable(tbname string) *ActiveRecord {
	ar.TableName = tbname

	return ar
}

func (ar *ActiveRecord) getTableName(rowsSlicePtr interface{}) string {
	return utils.Strings(utils.Struct{rowsSlicePtr}.GetTypeName()).SnakeCasedName()
}

func (ar *ActiveRecord) SetPK(pk string) *ActiveRecord {
	ar.PrimaryKey = pk

	return ar
}

func (ar *ActiveRecord) Where(querystring interface{}, args ...interface{}) *ActiveRecord {
	switch querystring := querystring.(type) {
	case string:
		ar.WhereStr = querystring
	case int:
		if ar.ParamIdentifier == "pg" {
			ar.WhereStr = fmt.Sprintf("%v%v%v = $%v", ar.QuoteIdentifier, ar.PrimaryKey, ar.QuoteIdentifier, ar.ParamIteration)
		} else {
			ar.WhereStr = fmt.Sprintf("%v%v%v = ?", ar.QuoteIdentifier, ar.PrimaryKey, ar.QuoteIdentifier)
			ar.ParamIteration++
		}
		args = append(args, querystring)
	}

	ar.ParamStr = args

	return ar
}

func (ar *ActiveRecord) Limit(start int, size ...int) *ActiveRecord {
	ar.LimitStr = start
	if len(size) > 0 {
		ar.OffsetStr = size[0]
	}

	return ar
}

func (ar *ActiveRecord) Offset(offset int) *ActiveRecord {
	ar.OffsetStr = offset

	return ar
}

func (ar *ActiveRecord) OrderBy(order string) *ActiveRecord {
	ar.OrderStr = order

	return ar
}

func (ar *ActiveRecord) Select(colums string) *ActiveRecord {
	ar.ColumnStr = colums

	return ar
}

func (ar *ActiveRecord) ScanPK(output interface{}) *ActiveRecord {
	if reflect.TypeOf(reflect.Indirect(reflect.ValueOf(output)).Interface()).Kind() == reflect.Slice {
		sliceValue := reflect.Indirect(reflect.ValueOf(output))
		sliceElementType := sliceValue.Type().Elem()
		for i := 0; i < sliceElementType.NumField(); i++ {
			tag := sliceElementType.Field(i).Tag
			if tag.Get("index") == "PK" {
				ar.PrimaryKey = sliceElementType.Field(i).Name
			}
		}
	} else {
		tt := reflect.TypeOf(reflect.Indirect(reflect.ValueOf(output)).Interface())
		for i := 0; i < tt.NumField(); i++ {
			tag := tt.Field(i).Tag
			if tag.Get("index") == "PK" {
				ar.PrimaryKey = tt.Field(i).Name
			}
		}
	}

	return ar
}

//The join_operator should be one of INNER, LEFT OUTER, CROSS etc - this will be prepended to JOIN
func (ar *ActiveRecord) Join(join_operator, tablename, condition string) *ActiveRecord {
	ar.JoinStr = fmt.Sprintf("%v JOIN %v ON %v", join_operator, tablename, condition)

	return ar
}

func (ar *ActiveRecord) GroupBy(keys string) *ActiveRecord {
	ar.GroupByStr = fmt.Sprintf("GROUP BY %v", keys)

	return ar
}

func (ar *ActiveRecord) Having(conditions string) *ActiveRecord {
	ar.HavingStr = fmt.Sprintf("HAVING %v", conditions)

	return ar
}

func (ar *ActiveRecord) Find(output interface{}) error {
	st := utils.Struct{output}
	ar.ScanPK(output)
	var keys []string
	results := st.StructToSnakeKeyMap()
	if ar.TableName == "" {
		ar.TableName = ar.getTableName(output)
	}

	for key, _ := range results {
		keys = append(keys, key)
	}

	ar.ColumnStr = strings.Join(keys, ", ")
	ar.Limit(1)
	resultsSlice, err := ar.FindMap()
	if err != nil {
		return err
	}

	if len(resultsSlice) == 0 {
		return nil
	} else if len(resultsSlice) == 1 {
		results := resultsSlice[0]
		utils.M(results).MapToStruct(output)
	} else {
		return errors.New("More Then One Records")
	}

	return nil
}

func (ar *ActiveRecord) FindAll(rowsSlicePtr interface{}) error {
	ar.ScanPK(rowsSlicePtr)
	sliceValue := reflect.Indirect(reflect.ValueOf(rowsSlicePtr))
	if sliceValue.Kind() != reflect.Slice {
		return errors.New("needs a pointer to a slice")
	}

	sliceElementType := sliceValue.Type().Elem()
	st := utils.Struct{reflect.New(sliceElementType).Interface()}
	var keys []string
	results := st.StructToSnakeKeyMap()
	if ar.TableName == "" {
		ar.TableName = ar.getTableName(rowsSlicePtr)
	}

	for key, _ := range results {
		keys = append(keys, key)
	}

	ar.ColumnStr = strings.Join(keys, ", ")
	resultsSlice, err := ar.FindMap()
	if err != nil {
		return err
	}

	for _, results := range resultsSlice {
		newValue := reflect.New(sliceElementType)
		utils.M(results).MapToStruct(newValue.Interface())
		sliceValue.Set(reflect.Append(sliceValue, reflect.Indirect(reflect.ValueOf(newValue.Interface()))))
	}

	return nil
}

func (ar *ActiveRecord) FindMap() ([]utils.M, error) {
	var resultsSlice []utils.M
	defer ar.Init()
	sqls := ar.GenerateSql()
	if OnDebug {
		fmt.Println(sqls)
		fmt.Println(ar)
	}

	s, err := ar.Db.Prepare(sqls)
	if err != nil {
		return nil, err
	}

	defer s.Close()
	res, err := s.Query(ar.ParamStr...)
	if err != nil {
		return nil, err
	}

	defer res.Close()
	fields, err := res.Columns()
	if err != nil {
		return nil, err
	}

	for res.Next() {
		var scanResultContainers []interface{}
		for i := 0; i < len(fields); i++ {
			var scanResultContainer interface{}
			scanResultContainers = append(scanResultContainers, &scanResultContainer)
		}

		if err := res.Scan(scanResultContainers...); err != nil {
			return nil, err
		}

		result := utils.Slice{fields}.SliceToMap(scanResultContainers)
		resultsSlice = append(resultsSlice, result)
	}

	return resultsSlice, nil
}

func (ar *ActiveRecord) GenerateSql() string {
	var a string
	if ar.ParamIdentifier == "mssql" {
		if ar.OffsetStr > 0 {
			a = fmt.Sprintf("select ROW_NUMBER() OVER(order by %v )as rownum,%v from %v",
				ar.PrimaryKey,
				ar.ColumnStr,
				ar.TableName)
			if ar.WhereStr != "" {
				a = fmt.Sprintf("%v WHERE %v", a, ar.WhereStr)
			}
			a = fmt.Sprintf("select * from (%v) "+
				"as a where rownum between %v and %v",
				a,
				ar.OffsetStr,
				ar.LimitStr)
		} else if ar.LimitStr > 0 {
			a = fmt.Sprintf("SELECT top %v %v FROM %v", ar.LimitStr, ar.ColumnStr, ar.TableName)
			if ar.WhereStr != "" {
				a = fmt.Sprintf("%v WHERE %v", a, ar.WhereStr)
			}
			if ar.GroupByStr != "" {
				a = fmt.Sprintf("%v %v", a, ar.GroupByStr)
			}
			if ar.HavingStr != "" {
				a = fmt.Sprintf("%v %v", a, ar.HavingStr)
			}
			if ar.OrderStr != "" {
				a = fmt.Sprintf("%v ORDER BY %v", a, ar.OrderStr)
			}
		} else {
			a = fmt.Sprintf("SELECT %v FROM %v", ar.ColumnStr, ar.TableName)
			if ar.WhereStr != "" {
				a = fmt.Sprintf("%v WHERE %v", a, ar.WhereStr)
			}
			if ar.GroupByStr != "" {
				a = fmt.Sprintf("%v %v", a, ar.GroupByStr)
			}
			if ar.HavingStr != "" {
				a = fmt.Sprintf("%v %v", a, ar.HavingStr)
			}
			if ar.OrderStr != "" {
				a = fmt.Sprintf("%v ORDER BY %v", a, ar.OrderStr)
			}
		}
	} else {
		a = fmt.Sprintf("SELECT %v FROM %v", ar.ColumnStr, ar.TableName)
		if ar.JoinStr != "" {
			a = fmt.Sprintf("%v %v", a, ar.JoinStr)
		}
		if ar.WhereStr != "" {
			a = fmt.Sprintf("%v WHERE %v", a, ar.WhereStr)
		}
		if ar.GroupByStr != "" {
			a = fmt.Sprintf("%v %v", a, ar.GroupByStr)
		}
		if ar.HavingStr != "" {
			a = fmt.Sprintf("%v %v", a, ar.HavingStr)
		}
		if ar.OrderStr != "" {
			a = fmt.Sprintf("%v ORDER BY %v", a, ar.OrderStr)
		}
		if ar.OffsetStr > 0 {
			a = fmt.Sprintf("%v LIMIT %v, %v", a, ar.OffsetStr, ar.LimitStr)
		} else if ar.LimitStr > 0 {
			a = fmt.Sprintf("%v LIMIT %v", a, ar.LimitStr)
		}
	}

	return a
}

//Execute sql
func (ar *ActiveRecord) Exec(finalQueryString string, args ...interface{}) (sql.Result, error) {
	rs, err := ar.Db.Prepare(finalQueryString)
	if err != nil {
		return nil, err
	}

	defer rs.Close()

	res, err := rs.Exec(args...)
	if err != nil {
		return nil, err
	}

	return res, nil
}

//if the struct has PrimaryKey == 0 insert else update
func (ar *ActiveRecord) Save(output interface{}) (int64, error) {
	ar.ScanPK(output)
	st := utils.Struct{output}
	results := st.StructToSnakeKeyMap()
	if ar.TableName == "" {
		ar.TableName = ar.getTableName(output)
	}

	id := results[strings.ToLower(ar.PrimaryKey)]
	var opRes int64
	delete(results, strings.ToLower(ar.PrimaryKey))
	if reflect.ValueOf(id).Int() == 0 {
		structPtr := reflect.ValueOf(output)
		structVal := structPtr.Elem()
		structField := structVal.FieldByName(ar.PrimaryKey)
		id, err := ar.Insert(results)
		if err != nil {
			return 0, err
		}

		opRes = id
		structField.Set(reflect.ValueOf(id))

		return opRes, nil
	} else {
		var condition string
		if ar.ParamIdentifier == "pg" {
			condition = fmt.Sprintf("%v%v%v=$%v", ar.QuoteIdentifier, strings.ToLower(ar.PrimaryKey), ar.QuoteIdentifier, ar.ParamIteration)
		} else {
			condition = fmt.Sprintf("%v%v%v=?", ar.QuoteIdentifier, ar.PrimaryKey, ar.QuoteIdentifier)
		}

		ar.Where(condition, id)
		opRes, err := ar.Update(results)
		if err != nil {
			return opRes, err
		}
	}

	return opRes, nil
}

//inert one info
func (ar *ActiveRecord) Insert(properties utils.M) (int64, error) {
	defer ar.Init()
	var keys []string
	var placeholders []string
	var args []interface{}
	for key, val := range properties {
		keys = append(keys, key)
		if ar.ParamIdentifier == "pg" {
			ds := fmt.Sprintf("$%d", ar.ParamIteration)
			placeholders = append(placeholders, ds)
		} else {
			placeholders = append(placeholders, "?")
		}
		ar.ParamIteration++
		args = append(args, val)
	}
	ss := fmt.Sprintf("%v,%v", ar.QuoteIdentifier, ar.QuoteIdentifier)
	statement := fmt.Sprintf("INSERT INTO %v%v%v (%v%v%v) VALUES (%v)",
		ar.QuoteIdentifier,
		ar.TableName,
		ar.QuoteIdentifier,
		ar.QuoteIdentifier,
		strings.Join(keys, ss),
		ar.QuoteIdentifier,
		strings.Join(placeholders, ", "))
	if OnDebug {
		fmt.Println(statement)
		fmt.Println(ar)
	}

	if ar.ParamIdentifier == "pg" {
		statement = fmt.Sprintf("%v RETURNING %v", statement, ar.PrimaryKey)
		var id int64
		ar.Db.QueryRow(statement, args...).Scan(&id)
		return id, nil
	} else {
		res, err := ar.Exec(statement, args...)
		if err != nil {
			return -1, err
		}

		id, err := res.LastInsertId()

		if err != nil {
			return -1, err
		}
		return id, nil
	}

	return -1, nil
}

//insert batch info
func (ar *ActiveRecord) InsertBatch(rows []utils.M) ([]int64, error) {
	var ids []int64
	tablename := ar.TableName
	if len(rows) <= 0 {
		return ids, nil
	}

	for i := 0; i < len(rows); i++ {
		ar.TableName = tablename
		id, _ := ar.Insert(rows[i])
		ids = append(ids, id)
	}

	return ids, nil
}

// update info
func (ar *ActiveRecord) Update(properties utils.M) (int64, error) {
	defer ar.Init()
	var updates []string
	var args []interface{}
	for key, val := range properties {
		if ar.ParamIdentifier == "pg" {
			ds := fmt.Sprintf("$%d", ar.ParamIteration)
			updates = append(updates, fmt.Sprintf("%v%v%v = %v", ar.QuoteIdentifier, key, ar.QuoteIdentifier, ds))
		} else {
			updates = append(updates, fmt.Sprintf("%v%v%v = ?", ar.QuoteIdentifier, key, ar.QuoteIdentifier))
		}
		args = append(args, val)
		ar.ParamIteration++
	}

	args = append(args, ar.ParamStr...)
	if ar.ParamIdentifier == "pg" {
		if n := len(ar.ParamStr); n > 0 {
			for i := 1; i <= n; i++ {
				ar.WhereStr = strings.Replace(ar.WhereStr, "$"+strconv.Itoa(i), "$"+strconv.Itoa(ar.ParamIteration), 1)
			}
		}
	}

	var condition string
	if ar.WhereStr != "" {
		condition = fmt.Sprintf("WHERE %v", ar.WhereStr)
	} else {
		condition = ""
	}

	statement := fmt.Sprintf("UPDATE %v%v%v SET %v %v",
		ar.QuoteIdentifier,
		ar.TableName,
		ar.QuoteIdentifier,
		strings.Join(updates, ", "),
		condition)
	if OnDebug {
		fmt.Println(statement)
		fmt.Println(ar)
	}

	res, err := ar.Exec(statement, args...)
	if err != nil {
		return -1, err
	}

	id, err := res.RowsAffected()

	if err != nil {
		return -1, err
	}

	return id, nil
}

func (ar *ActiveRecord) Delete(output interface{}) (int64, error) {
	defer ar.Init()
	ar.ScanPK(output)

	st := utils.Struct{output}
	results := st.StructToSnakeKeyMap()
	if ar.TableName == "" {
		ar.TableName = ar.getTableName(output)
	}

	id := results[strings.ToLower(ar.PrimaryKey)]
	condition := fmt.Sprintf("%v%v%v='%v'", ar.QuoteIdentifier, strings.ToLower(ar.PrimaryKey), ar.QuoteIdentifier, id)
	statement := fmt.Sprintf("DELETE FROM %v%v%v WHERE %v",
		ar.QuoteIdentifier,
		ar.TableName,
		ar.QuoteIdentifier,
		condition)
	if OnDebug {
		fmt.Println(statement)
		fmt.Println(ar)
	}

	res, err := ar.Exec(statement)
	if err != nil {
		return -1, err
	}
	Affectid, err := res.RowsAffected()

	if err != nil {
		return -1, err
	}

	return Affectid, nil
}

func (ar *ActiveRecord) DeleteAll(rowsSlicePtr interface{}) (int64, error) {
	defer ar.Init()
	ar.ScanPK(rowsSlicePtr)
	if ar.TableName == "" {
		ar.TableName = ar.getTableName(rowsSlicePtr)
	}

	var ids []string
	val := reflect.Indirect(reflect.ValueOf(rowsSlicePtr))
	if val.Len() == 0 {
		return 0, nil
	}

	for i := 0; i < val.Len(); i++ {
		results := utils.Struct{val.Index(i).Interface()}.StructToSnakeKeyMap()
		id := results[strings.ToLower(ar.PrimaryKey)]
		switch id.(type) {
		case string:
			ids = append(ids, id.(string))
		case int, int64, int32:
			str := fmt.Sprintf("%v", id)
			ids = append(ids, str)
		}
	}

	condition := fmt.Sprintf("%v%v%v in ('%v')", ar.QuoteIdentifier, strings.ToLower(ar.PrimaryKey), ar.QuoteIdentifier, strings.Join(ids, "','"))
	statement := fmt.Sprintf("DELETE FROM %v%v%v WHERE %v",
		ar.QuoteIdentifier,
		ar.TableName,
		ar.QuoteIdentifier,
		condition)
	if OnDebug {
		fmt.Println(statement)
		fmt.Println(ar)
	}

	res, err := ar.Exec(statement)
	if err != nil {
		return -1, err
	}

	Affectid, err := res.RowsAffected()

	if err != nil {
		return -1, err
	}

	return Affectid, nil
}

func (ar *ActiveRecord) DeleteRow() (int64, error) {
	defer ar.Init()
	var condition string
	if ar.WhereStr != "" {
		condition = fmt.Sprintf("WHERE %v", ar.WhereStr)
	} else {
		condition = ""
	}

	statement := fmt.Sprintf("DELETE FROM %v%v%v %v",
		ar.QuoteIdentifier,
		ar.TableName,
		ar.QuoteIdentifier,
		condition)
	if OnDebug {
		fmt.Println(statement)
		fmt.Println(ar)
	}

	res, err := ar.Exec(statement, ar.ParamStr...)
	if err != nil {
		return -1, err
	}

	Affectid, err := res.RowsAffected()

	if err != nil {
		return -1, err
	}

	return Affectid, nil
}

func (ar *ActiveRecord) Close() {
	ar.Db.Close()
}

func (ar *ActiveRecord) Init() {
	ar.TableName = ""
	ar.LimitStr = 0
	ar.OffsetStr = 0
	ar.WhereStr = ""
	ar.ParamStr = make([]interface{}, 0)
	ar.OrderStr = ""
	ar.ColumnStr = "*"
	ar.PrimaryKey = "id"
	ar.JoinStr = ""
	ar.GroupByStr = ""
	ar.HavingStr = ""
	ar.ParamIteration = 1
}
