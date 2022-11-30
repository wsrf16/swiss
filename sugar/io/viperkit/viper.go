package viperkit

import (
	"bytes"
	"github.com/spf13/viper"
	"github.com/wsrf16/swiss/sugar/io/pathkit"
)

func UnmarshalSToT[T any](text string) (*T, error) {
	v := viper.New()

	if err := v.ReadConfig(bytes.NewBuffer([]byte(text))); err != nil {
		return nil, err
	}
	t := new(T)
	if err := v.Unmarshal(t); err != nil {
		return nil, err
	}
	return t, nil
}

func UnmarshalS(text string) (*viper.Viper, error) {
	v := viper.New()

	if err := v.ReadConfig(bytes.NewBuffer([]byte(text))); err != nil {
		return nil, err
	}
	return v, nil
}

func UnmarshalFileToT[T any](configFile string) (*T, error) {
	v := viper.New()
	v.SetConfigFile(configFile)
	// v.SetConfigType(configType)

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}
	t := new(T)
	if err := v.Unmarshal(t); err != nil {
		return nil, err
	}
	return t, nil
}

func UnmarshalCurrentFileToT[T any](currentFile string) (*T, error) {
	file := pathkit.GetPWD(currentFile)
	return UnmarshalFileToT[T](file)
}

func UnmarshalFile(filePath string) (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigFile(filePath)

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}
	return v, nil
}

func UnmarshalCurrentFile(currentFile string) (*viper.Viper, error) {
	file := pathkit.GetPWD(currentFile)
	return UnmarshalFile(file)
}
