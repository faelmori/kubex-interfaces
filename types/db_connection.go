package types

type ISQLiteDatabase interface {
	IDatabase
	Vacuum() error
}

type IPostgreSQLDatabase interface {
	IDatabase
	Listen(channel string) error
	Notify(channel string, payload string) error
}

type IConnectionPool interface {
	GetConnection() (IDatabase, error)
	ReleaseConnection(IDatabase) error
	GetActiveConnections() int
	GetConnectionLimit() int
	SetConnectionLimit(limit int) error
}
