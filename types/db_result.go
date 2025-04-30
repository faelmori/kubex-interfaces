package types

type IDBResult interface {
	RowsAffected() int
	LastInsertID() (int, error)
}

type IRowSet interface {
	Next() bool
	Scan(dest ...any) error
	Close() error
}
