package jsonkit

import (
	"github.com/json-iterator/go"
	"os"
)

var jsonIterator = jsoniter.ConfigCompatibleWithStandardLibrary

func Marshal(t interface{}) (string, error) {
	b, err := jsonIterator.Marshal(t)
	return string(b), err
}

func Unmarshal(j string, t interface{}) error {
	b := []byte(j)
	err := jsonIterator.Unmarshal(b, t)
	return err
}

func UnmarshalSToT[T interface{}](j string) (*T, error) {
	b := []byte(j)
	t := new(T)
	if err := jsonIterator.Unmarshal(b, t); err != nil {
		return nil, err
	} else {
		return t, err
	}
}

func UnmarshalBToT[T interface{}](bytes []byte) (*T, error) {
	t := new(T)
	if err := jsonIterator.Unmarshal(bytes, t); err != nil {
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

	encoder := jsonIterator.NewEncoder(filePtr)
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

	decoder := jsonIterator.NewDecoder(filePtr)
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

	decoder := jsonIterator.NewDecoder(filePtr)
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
