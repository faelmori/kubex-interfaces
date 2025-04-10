package config

import (
	icc "github.com/faelmori/kbxutils/utils/interfaces"
	ig "github.com/faelmori/kubex-interfaces/config"
)

func NewBaseProperty[T any](name string) icc.Property[T] { return ig.NewBaseProperty[T](name) }
