package config

import (
	ig "github.com/faelmori/kubex-interfaces/config"
)

func NewBaseProperty[T any](name string) ig.Property[T] { return ig.NewProperty[T](name, nil) }
