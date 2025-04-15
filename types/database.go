package types

import (
	"time"
)

type IDatabaseService[T any] interface {
	Connect() error
	Disconnect() error
	Reconnect() error
	CheckHealth() error
	ExecuteQuery(query string, args ...interface{}) ([]T, error)
	GetConnectionDetails() Property[any]              // Retorna uma propriedade dinâmica para detalhes da conexão
	Monitor(interval time.Duration) error             // Monitoramento contínuo do banco
	AddListener(event string, listener func(T) error) // Escuta eventos no banco
}
