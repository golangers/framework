package activerecord

import (
	"database/sql"
	"golanger.com/framework/utils"
)

type DataBase interface {
	SetTable(tableName string) DataBase
	getTableName(rowsSlicePtr interface{}) string
	SetPK(pk string) DataBase
	Where(querystring interface{}, args ...interface{}) DataBase
	Limit(start int, size ...int) DataBase
	Offset(offset int) DataBase
	OrderBy(order string) DataBase
	Select(colums string) DataBase
	ScanPK(output interface{}) DataBase
	Join(join_operator, tablename, condition string) DataBase
	GroupBy(keys string) DataBase
	Having(conditions string) DataBase
	Find(output interface{}) error
	FindAll(rowsSlicePtr interface{}) error
	FindMap() ([]utils.M, error)
	GenerateSql() string
	Exec(finalQueryString string, args ...interface{}) (sql.Result, error)
	Save(output interface{}) error
	Insert(properties utils.M) (int64, error)
	InsertBatch(rows []utils.M) ([]int64, error)
	Update(properties utils.M) (int64, error)
	Delete(output interface{}) (int64, error)
	DeleteAll(rowsSlicePtr interface{}) (int64, error)
	DeleteRow()
	Close()
}
