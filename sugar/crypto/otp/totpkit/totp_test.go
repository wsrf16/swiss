package totpkit

import (
	"fmt"
	"testing"
	"time"
)

func TestTOTP(t *testing.T) {
	for true {
		if code, err := GenerateGoogleCode("ABCDEFGH234567MN"); err != nil {
			panic(err)
		} else {
			fmt.Println(code)
		}
		time.Sleep(time.Second * 1)
		time.Sleep(time.Duration(time.Second * 1))
	}
}
