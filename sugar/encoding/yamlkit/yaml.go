package yamlkit

import (
	"gopkg.in/yaml.v3"
	"os"
)

func Marshal(t interface{}) (string, error) {
	bytes, err := yaml.Marshal(t)
	return string(bytes), err
}

func Unmarshal(js string, t interface{}) error {
	bytes := []byte(js)
	err := yaml.Unmarshal(bytes, t)
	return err
}

func UnmarshalSToT[T interface{}](js string) (*T, error) {
	bytes := []byte(js)
	t := new(T)
	if err := yaml.Unmarshal(bytes, t); err != nil {
		return nil, err
	} else {
		return t, err
	}
}

func UnmarshalBToT[T interface{}](bytes []byte) (*T, error) {
	t := new(T)
	if err := yaml.Unmarshal(bytes, t); err != nil {
		return nil, err
	} else {
		return t, err
	}
}

func MarshalToFile(info interface{}, file string) error {
	filePtr, err := os.Create(file)
	if err != nil {
		return err
	}
	defer filePtr.Close()

	encoder := yaml.NewEncoder(filePtr)
	err = encoder.Encode(info)
	if err != nil {
		return err
	}
	return nil
}

func UnmarshalFromFile(file string, info interface{}) error {
	filePtr, err := os.Open(file)
	if err != nil {
		return err
	}
	defer filePtr.Close()

	decoder := yaml.NewDecoder(filePtr)
	err = decoder.Decode(&info)
	if err != nil {
		return err
	}
	return nil
}

func UnmarshalFromFileToT[T any](file string) (*T, error) {
	filePtr, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer filePtr.Close()

	decoder := yaml.NewDecoder(filePtr)
	t := new(T)
	err = decoder.Decode(t)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func DeepCopy[T any](t any) (*T, error) {
	marshal, err := Marshal(t)
	if err != nil {
		return nil, err
	}
	return UnmarshalSToT[T](marshal)
}
