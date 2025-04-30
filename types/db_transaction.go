package types

type ITransaction interface {
	Commit() error
	Rollback() error
	Execute(query string, args ...any) (IResult, error)
	Query(query string, args ...any) (IRowSet, error)
}
