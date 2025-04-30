package drivers

import (
	"database/sql"
	"fmt"
	. "github.com/faelmori/kubex-interfaces/databases/types"
	"github.com/faelmori/logz"
	"time"
)

// ConnectDB establishes a connection to the database using the provided configuration.
func ConnectDB(config *DBConfigBase) (*sql.DB, error) {

	db, err := sql.Open(config.Connection.Driver, config.Connection.Dsn)
	if err != nil {
		return nil, fmt.Errorf("falha ao conectar ao banco de dados: %w", err)
	}

	// Verifica a conexão
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("falha ao verificar a conexão com o banco de dados: %w", err)
	}

	return db, nil
}

// ReconnectDB attempts to reconnect to the database using the provided configuration.
func ReconnectDB(config DBConfigBase) (*sql.DB, error) {
	var db *sql.DB
	var err error
	maxRetries := 5
	retryInterval := 10 * time.Second

	for i := 0; i < maxRetries; i++ {
		db, err = sql.Open(config.Connection.Driver, config.Connection.Dsn)
		if err == nil {
			if pingErr := db.Ping(); pingErr == nil {
				return db, nil
			}
		}
		time.Sleep(retryInterval)
	}

	return nil, fmt.Errorf("falha ao reconectar ao banco de dados após %d tentativas: %w", maxRetries, err)
}

// VacuumDatabase executes the VACUUM command on the specified SQLite database.
func VacuumDatabase(dbPath string) error {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("falha ao abrir o banco de dados: %w", err)
	}
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	_, err = db.Exec("VACUUM")
	if err != nil {
		return fmt.Errorf("falha ao executar VACUUM: %w", err)
	}

	logz.InfoCtx("VACUUM executado com sucesso", map[string]interface{}{})
	return nil
}

// CloseDB closes the database connection.
func CloseDB(db *sql.DB) error {
	if db != nil {
		if err := db.Close(); err != nil {
			return fmt.Errorf("falha ao fechar a conexão com o banco de dados: %w", err)
		}
	}
	return nil
}
