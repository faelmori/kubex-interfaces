package databases

import (
	"database/sql"
	"fmt"
	c "github.com/faelmori/kubex-interfaces/config"
	config2 "github.com/faelmori/kubex-interfaces/databases/core/config"
	"github.com/faelmori/kubex-interfaces/databases/etl"
	t "github.com/faelmori/kubex-interfaces/databases/types"
	tp "github.com/faelmori/kubex-interfaces/types"
	l "github.com/faelmori/logz"
	"github.com/google/uuid"
	"reflect"
	"sync"
	"time"
)

type DatabaseService[T any] struct {
	tp.IDatabaseService[T]
	mu sync.RWMutex
	wg sync.WaitGroup
	//mp         sync.Map
	logger      l.Logger
	mapper      tp.IMapper[t.DBConfigBase]
	ID          string
	Properties  map[string]c.Property[any]
	Connections map[string]any
}

func setupDatabaseBasicProperties[T any](dbSrv *DatabaseService[T]) error {
	// Adiciona propriedades da configuração do banco de dados CENTRAL
	dbSrv.Properties["connectionDetails"] = c.NewProperty[t.DBConfigBase]("connectionDetails", &t.DBConfigBase{})

	// Definição de valores básicos
	dbSrv.Properties["activeConnections"] = c.NewProperty[int]("activeConnections", nil)
	if err := dbSrv.Properties["activeConnections"].SetValue(0, nil); err != nil {
		dbSrv.logger.ErrorCtx("Error setting activeConnections", map[string]any{
			"context":  "DatabaseService",
			"action":   "SetValue",
			"error":    err,
			"showData": true,
		})
		return nil
	}
	dbSrv.Properties["connectionLimit"] = c.NewProperty[int]("connectionLimit", nil)
	if err := dbSrv.Properties["connectionLimit"].SetValue(100, nil); err != nil {
		dbSrv.logger.ErrorCtx("Error setting connectionLimit", map[string]any{
			"context":  "DatabaseService",
			"action":   "SetValue",
			"error":    err,
			"showData": true,
		})
		return nil
	}

	// Adiciona validadores para as propriedades
	if err := dbSrv.Properties["connectionLimit"].AddValidator("connectionLimit", func(value any) error {
		if value.(int) <= 0 {
			return fmt.Errorf("Connection limit must be greater than 0")
		}
		return nil
	}); err != nil {
		dbSrv.logger.ErrorCtx("Error adding validator for connectionLimit", map[string]any{
			"context":  "DatabaseService",
			"action":   "AddValidator",
			"error":    err,
			"showData": true,
		})
		return nil
	}

	return nil
}

func NewDatabaseService[T any](logger l.Logger) tp.IDatabaseService[T] {
	if logger == nil {
		logger = l.GetLogger("DatabaseService")
	}
	dbSrv := DatabaseService[T]{
		mu: sync.RWMutex{},
		wg: sync.WaitGroup{},
		//mp:         sync.Map{},
		ID:          uuid.NewString(),
		logger:      logger,
		mapper:      tp.NewMapper[t.DBConfigBase](),
		Properties:  make(map[string]c.Property[any]),
		Connections: make(map[string]any),
	}

	if err := setupDatabaseBasicProperties(&dbSrv); err != nil {
		dbSrv.logger.ErrorCtx("Error setting up database basic properties", map[string]any{
			"context":  "DatabaseService",
			"action":   "setupDatabaseBasicProperties",
			"error":    err,
			"showData": true,
		})
		return nil
	}

	return &dbSrv
}

func (ds *DatabaseService[T]) GetID() string {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	return ds.ID
}

func (ds *DatabaseService[T]) GetName() string {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	return ds.ID
}

func (ds *DatabaseService[T]) GetConfig(format string) ([]byte, error) {
	// Valor principal do Property fica num atomic.Pointer[T], não é necessário fazer o lock de leitura agora (certo?)
	config := ds.Properties["connectionDetails"].GetValue()

	ds.mu.RLock()
	defer ds.mu.RUnlock()

	if serializedData, err := ds.mapper.Serialize(nil, config.(*t.DBConfigBase), format); err != nil {
		ds.logger.ErrorCtx("Error serializing config", map[string]any{
			"context":  "DatabaseService",
			"action":   "Serialize",
			"error":    err,
			"showData": true,
		})
		return nil, err
	} else {
		return serializedData, nil
	}
}

func (ds *DatabaseService[T]) SaveConfig(format string, data []byte) ([]byte, error) {
	var config t.DBConfigBase

	ds.mu.Lock()
	defer ds.mu.Unlock()

	if err := ds.mapper.Deserialize(data, &config, format); err != nil {
		return nil, err
	}

	if setValErr := ds.Properties["connectionDetails"].SetValue(config, nil); setValErr != nil {
		ds.logger.ErrorCtx("Error setting connectionDetails", map[string]any{
			"context":  "DatabaseService",
			"action":   "SetValue",
			"error":    setValErr,
			"showData": true,
		})
		return nil, setValErr
	}

	if serializedData, err := ds.mapper.Serialize(nil, &config, format); err != nil {
		return nil, err
	} else {
		return serializedData, nil
	}
}

func (ds *DatabaseService[T]) Connect() error {
	if connectionDetails := ds.Properties["connectionDetails"].GetValue().(*t.DBConfigBase); connectionDetails == nil {
		return fmt.Errorf("no configuration details provided")
	} else {
		ds.mu.Lock()
		defer ds.mu.Unlock()

		if db, dbErr := ds.connectDB(connectionDetails); dbErr != nil {
			ds.logger.ErrorCtx("Failed to connect to database: "+dbErr.Error(), nil)
			return dbErr
		} else {
			ds.Connections[connectionDetails.Name] = tp.NewProperty[sql.DB]("SQLServer_Execmplo", db)
		}
		return nil
	}
}

func (ds *DatabaseService[T]) LoadConfig() error {
	cfgT := config2.NewDynamicConfigManager("config.json", ds.logger)
	var ddd *sql.DB

	if databaseConfig, databaseConfigErr := cfgT.GetConfig(); databaseConfigErr != nil {
		return databaseConfigErr
	} else {
		// Protegido com ponteiro ainda, a conexão.. De forma a protegê-lo nbão só de concorrentes, mas também de acesso indevido
		// Porém continua com o chan, não tirei né... rssr
		if err := ds.Properties["some_connextion_details"].SetValue(databaseConfig, nil); err != nil {
			return err
		}
	}

	var config t.DBConfigBase

	ds.mu.Lock()
	defer ds.mu.Unlock()

	if err := ds.mapper.Deserialize(ddd, &config, "json"); err != nil {
		ds.logger.ErrorCtx(fmt.Sprintf("Error deserializing config: %v", err), nil)
		return err
	}

	if setValErr := ds.Properties["connectionDetails"].SetValue(config, nil); setValErr != nil {
		ds.logger.ErrorCtx(fmt.Sprintf("Error setting connectionDetails: %v", setValErr), nil)
		return setValErr
	}
	return nil
}

func (ds *DatabaseService[T]) ExecuteQuery(query string, args ...interface{}) ([]T, error) {
	queryChan := ds.Properties["queryChannel"].GetValue().(tp.IChannel[string, int])
	if sendErr := queryChan.Send(fmt.Sprintf(query, args...)); sendErr != nil {
		ds.logger.ErrorCtx("Error sending query", map[string]any{
			"context":  "DatabaseService",
			"action":   "Send",
			"error":    sendErr,
			"showData": true,
		})
		return nil, sendErr
	}

	resultChan := ds.Properties["resultChannel"].GetValue().(tp.IChannel[[]T, int])
	if _, typ, err := resultChan.Listen(); err != nil {
		return nil, err
	} else if typ != reflect.TypeFor[[]T]() {
		typeErr := fmt.Sprintf("Type mismatch: expected %v, got %v", reflect.TypeFor[T](), typ)
		ds.logger.ErrorCtx(typeErr, nil)
		return nil, fmt.Errorf(typeErr)
	} else {
		return reflect.MakeSlice(reflect.TypeFor[[]T](), 0, 0).Interface().([]T), nil
	}
}

func (ds *DatabaseService[T]) GetProperty(name string) (any, error) {
	ds.mu.RLock()
	defer ds.mu.RUnlock()
	if prop, ok := ds.Properties[name]; ok {
		return prop.GetValue(), nil
	}
	return nil, fmt.Errorf("property %s not found", name)
}

// connectDB establishes a connection to the database using the provided configuration.
func (ds *DatabaseService[T]) connectDB(config *t.DBConfigBase) (*sql.DB, error) {

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

// reconnectDB attempts to reconnect to the database using the provided configuration.
func (ds *DatabaseService[T]) reconnectDB(config *t.DBConfigBase) (*sql.DB, error) {
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

// vacuumDatabase executes the VACUUM command on the specified SQLite database.
func (ds *DatabaseService[T]) vacuumDatabase(dbPath string) error {
	if db, err := sql.Open("sqlite3", dbPath); err != nil {
		return fmt.Errorf("falha ao abrir o banco de dados: %w", err)
	} else {
		defer func(db *sql.DB) {
			_ = db.Close()
		}(db)

		if _, err := db.Exec("VACUUM"); err != nil {
			return fmt.Errorf("falha ao executar VACUUM: %w", err)
		}
	}

	l.InfoCtx("VACUUM executado com sucesso", map[string]interface{}{})

	return nil
}

// closeDB closes the database connection.
func (ds *DatabaseService[T]) closeDB(db *sql.DB) error {
	if db != nil {
		if err := db.Close(); err != nil {
			return fmt.Errorf("falha ao fechar a conexão com o banco de dados: %w", err)
		}
	}
	return nil
}

func main() {
	l.GetLogger("TESTE ETL")

	configManager := etl.ConfigManager{}
	if loadCfgErr := configManager.LoadConfig("config.json"); loadCfgErr != nil {
		l.ErrorCtx(fmt.Sprintf("Failed to load config: %v", loadCfgErr), nil)
		return
	}

	etlPipeline := &etl.ETLPipeline{
		ConfigManager: configManager,
		DataExtractor: etl.DataExtractor{Config: configManager},
		DataTransformer: etl.DataTransformer{
			Transformations: configManager.Config.Transformations,
		},
		DataLoader: etl.DataLoader{
			DestinationType: configManager.Config.DestinationType,
			DestinationConn: configManager.Config.DestinationConnectionString,
		},
		Logger: l.GetLogger("ETLPipeline"),
	}

	if err := etlPipeline.RunPipeline(); err != nil {
		fmt.Println("ETL pipeline failed:", err)
	}
}
