package sstrans

import (
	"github.com/wsrf16/swiss/netkit/layer/transport/sstrans/core"
)

type ShadowSocksConfig struct {
	Cipher    *core.Cipher
	Algorithm string
	Key       []byte
	Password  string
}

func BuildConfig(algorithm string, key []byte, password string) (*ShadowSocksConfig, error) {
	//cipher, err := core.PickCipher(algorithm, key, password)
	//if err != nil {
	//    return nil, err
	//}
	config := new(ShadowSocksConfig)
	//config.Cipher = &cipher
	config.Algorithm = algorithm
	config.Key = key
	config.Password = password
	return config, nil
}
