package jsonkit

import (
	"encoding/json"
	"os"
)

func Marshal(t interface{}) (string, error) {
	b, err := json.Marshal(t)
	return string(b), err
}

func Unmarshal(j string, t interface{}) error {
	b := []byte(j)
	err := json.Unmarshal(b, t)
	return err
}

func UnmarshalSToT[T interface{}](j string) (*T, error) {
	b := []byte(j)
	t := new(T)
	if err := json.Unmarshal(b, t); err != nil {
		return nil, err
	} else {
		return t, err
	}
}

func UnmarshalBToT[T interface{}](bytes []byte) (*T, error) {
	t := new(T)
	if err := json.Unmarshal(bytes, t); err != nil {
		return nil, err
	} else {
		return t, err
	}
}

func MarshalToFile(t interface{}, file string) error {
	filePtr, err := os.Create(file)
	if err != nil {
		return err
	}
	defer filePtr.Close()

	encoder := json.NewEncoder(filePtr)
	err = encoder.Encode(t)
	if err != nil {
		return err
	}
	return nil
}

func UnmarshalFromFile(file string, t interface{}) error {
	filePtr, err := os.Open(file)
	if err != nil {
		return err
	}
	defer filePtr.Close()

	decoder := json.NewDecoder(filePtr)
	err = decoder.Decode(&t)
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

	decoder := json.NewDecoder(filePtr)
	t := new(T)
	err = decoder.Decode(t)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func DeepClone[T any](t T) (*T, error) {
	marshal, err := Marshal(t)
	if err != nil {
		return nil, err
	}
	return UnmarshalSToT[T](marshal)
}

func DeepCloneSimilar[T any](t any) (*T, error) {
	marshal, err := Marshal(t)
	if err != nil {
		return nil, err
	}
	return UnmarshalSToT[T](marshal)
}
