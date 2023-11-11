package sshkit

import (
	"errors"
)

type SSHPropertyPool struct {
	SSHParams []*SSHProperty
}

var propertyPool *SSHPropertyPool

func init() {
}

func GetSingleton() *SSHPropertyPool {
	if propertyPool == nil {
		propertyPool = newPool()
	}
	return propertyPool
}

func newPool() *SSHPropertyPool {
	pool := new(SSHPropertyPool)
	pool.SSHParams = make([]*SSHProperty, 0, 16)
	return pool
}

func (pool *SSHPropertyPool) Get(id string) (*SSHProperty, error) {
	for _, c := range pool.SSHParams {
		if id == c.Id {
			return c, nil
		}
	}
	return nil, errors.New("not found id " + id)
}

func (pool *SSHPropertyPool) Put(config *SSHProperty) {
	pool.SSHParams = append(pool.SSHParams, config)
}

func (pool *SSHPropertyPool) PutNew(id string, addr string) *SSHProperty {
	config := NewSSHConfigSpecifyId(id, addr)
	pool.Put(config)
	return config
}

func (pool *SSHPropertyPool) PutMulti(configs []*SSHProperty) {
	pool.SSHParams = append(pool.SSHParams, configs...)
}
