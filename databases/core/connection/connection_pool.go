package connection

import (
	"fmt"
	c "github.com/faelmori/kubex-interfaces/config"
	t "github.com/faelmori/kubex-interfaces/databases/types"
	tp "github.com/faelmori/kubex-interfaces/types"
	l "github.com/faelmori/logz"
	"github.com/google/uuid"
	"sync"
)

type DBConnectionPool struct {
	mu          sync.Mutex
	wg          sync.WaitGroup
	logger      l.Logger
	ID          string
	name        string
	Properties  map[string]c.Property[any]
	connections map[string]tp.IChannel[any, int]
}

func setupDBBasicProperties(pool *DBConnectionPool) error {
	// Adiciona propriedades da configuração do banco de dados CENTRAL
	pool.Properties["connectionDetails"] = c.NewProperty[t.DBConfigBase]("connectionDetails", &t.DBConfigBase{})

	pool.Properties["configRegistry"] = c.NewProperty[map[string]any]("configRegistry", nil)
	if setValErr := pool.Properties["configRegistry"].SetValue(make(map[string]any), nil); setValErr != nil {
		pool.logger.ErrorCtx("Error setting configRegistry", map[string]any{
			"context":  "DBConnectionPool",
			"action":   "SetValue",
			"error":    setValErr,
			"showData": true,
		})
		return nil
	}

	pool.Properties["connections"] = c.NewProperty[map[string]tp.IChannel[any, int]]("connections", nil)
	if setValErr := pool.Properties["connections"].SetValue(make(map[string]tp.IChannel[any, int]), nil); setValErr != nil {
		pool.logger.ErrorCtx("Error setting connections", map[string]any{
			"context":  "DBConnectionPool",
			"action":   "SetValue",
			"error":    setValErr,
			"showData": true,
		})
		return nil
	}

	pool.Properties["activeConnections"] = c.NewProperty[int]("activeConnections", nil)
	if err := pool.Properties["activeConnections"].SetValue(0, nil); err != nil {
		pool.logger.ErrorCtx("Error setting activeConnections", map[string]any{
			"context":  "DBConnectionPool",
			"action":   "SetValue",
			"error":    err,
			"showData": true,
		})
		return nil
	}
	if err := pool.Properties["activeConnections"].AddValidator("activeConnections", func(value any) error {
		if value.(int) < 0 {
			return fmt.Errorf("activeConnections cannot be negative")
		}
		return nil
	}); err != nil {
		pool.logger.ErrorCtx("Error adding validator for activeConnections", map[string]any{
			"context":  "DBConnectionPool",
			"action":   "AddValidator",
			"error":    err,
			"showData": true,
		})
		return nil
	}

	return nil
}

func NewDBConnectionPool(limit *int, logger l.Logger) *DBConnectionPool {
	if logger == nil {
		logger = l.GetLogger("Kubex")
	}
	pool := DBConnectionPool{
		mu:          sync.Mutex{},
		wg:          sync.WaitGroup{},
		logger:      logger,
		ID:          uuid.NewString(),
		name:        "DBConnectionPool",
		Properties:  make(map[string]c.Property[any]),
		connections: make(map[string]tp.IChannel[any, int]),
	}
	pool.Properties["connectionLimit"] = c.NewProperty[int]("connectionLimit", limit)

	if err := setupDBBasicProperties(&pool); err != nil {
		pool.logger.ErrorCtx("Error setting up DB connection pool properties", map[string]any{
			"context":  "DBConnectionPool",
			"action":   "setupDBBasicProperties",
			"error":    err,
			"showData": true,
		})
		return nil
	}

	return &pool
}

func (dbPool *DBConnectionPool) AddConnection(name string, conn tp.IChannel[t.DBConfigBase, int]) error {
	dbPool.mu.Lock()
	defer dbPool.mu.Unlock()

	if len(dbPool.connections) >= dbPool.Properties["connectionLimit"].GetValue().(int) {
		return fmt.Errorf("Connection limit reached")
	}

	dbPool.connections[name] = conn
	active := len(dbPool.connections)
	dbPool.Properties["activeConnections"].SetValue(active, nil)
	return nil
}
