package databases

import (
	"fmt"
	c "github.com/faelmori/kubex-interfaces/config"
	t "github.com/faelmori/kubex-interfaces/databases/types"
	tp "github.com/faelmori/kubex-interfaces/types"
	l "github.com/faelmori/logz"
	"github.com/google/uuid"
	"reflect"
	"sync"
)

type DatabaseService[T any] struct {
	tp.IDatabaseService[T]
	mu         sync.RWMutex
	wg         sync.WaitGroup
	mp         sync.Map
	db         tp.IChannel[T, int] // Usamos IChannel para canalizar interações com o banco
	logger     l.Logger
	mapper     tp.IMapper[t.DBConfigBase]
	ID         string
	Properties map[string]c.Property[any]
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
		mu:         sync.RWMutex{},
		wg:         sync.WaitGroup{},
		ID:         uuid.NewString(),
		db:         nil,
		logger:     logger,
		mapper:     tp.NewMapper[t.DBConfigBase](),
		Properties: make(map[string]c.Property[any]),
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

func (ds *DatabaseService[T]) GetID() string { return ds.ID }

func (ds *DatabaseService[T]) GetName() string { return ds.ID }

func (ds *DatabaseService[T]) GetConfig(format string) ([]byte, error) {
	config := ds.Properties["connectionDetails"].GetValue()
	serializedData, err := ds.mapper.Serialize(nil, config.(*t.DBConfigBase), format)
	if err != nil {
		return nil, err
	}
	return serializedData, nil
}

func (ds *DatabaseService[T]) SaveConfig(format string, data []byte) ([]byte, error) {
	var config t.DBConfigBase
	err := ds.mapper.Deserialize(data, &config, format)
	if err != nil {
		return nil, err
	}
	if setValErr := ds.Properties["connectionDetails"].SetValue(config, nil); setValErr != nil {
		return nil, setValErr
	}
	serializedData, err := ds.mapper.Serialize(nil, &config, format)
	if err != nil {
		return nil, err
	}
	return serializedData, nil
}

func (ds *DatabaseService[T]) LoadConfig(format string, data []byte) error {
	var config t.DBConfigBase
	err := ds.mapper.Deserialize(data, &config, format)
	if err != nil {
		return err
	}
	if setValErr := ds.Properties["connectionDetails"].SetValue(config, nil); setValErr != nil {
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
	_, tp, err := resultChan.Listen()
	if err != nil {
		return nil, err
	}
	if tp != reflect.TypeFor[[]T]() {
		typeErr := fmt.Sprintf("Type mismatch: expected %v, got %v", reflect.TypeFor[T](), tp)
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
