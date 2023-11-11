package tokenkit

import (
	"errors"
	"github.com/wsrf16/swiss/sugar/base/timekit"
	"github.com/wsrf16/swiss/sugar/crypto/rsakit"
	"github.com/wsrf16/swiss/sugar/encoding/jsonkit"
	"time"
)

// const INTERVAL string  = "|||"
const DEFAULT_EXPIRED time.Duration = 60 * time.Minute

func Sign(publicKeyText string, pwd string, salt string, signedTime int64) (string, error) {
	combine, err := marshal(pwd, salt, signedTime)
	if err != nil {
		return "", err
	}
	pub := rsakit.ParsePublicKey(publicKeyText)

	base64, err := rsakit.EncryptIntoBase64(*pub, combine)
	if err != nil {
		return "", err
	}
	return base64, nil
}

type TokenStruct struct {
	Pwd  string `json:"pwd"`
	Salt string `json:"salt"`
	Time int64  `json:"time"`
}

func marshal(pwd string, salt string, accessTime int64) (string, error) {
	tokenStruct := TokenStruct{Pwd: pwd, Salt: salt, Time: accessTime}
	if mar, err := jsonkit.MarshalToJson(tokenStruct); err != nil {
		return "", err
	} else {
		return mar, nil
	}
}

func unmarshal(mar string) (*TokenStruct, error) {
	t := new(TokenStruct)
	if err := jsonkit.Unmarshal(mar, t); err != nil {
		return nil, err
	} else {
		return t, nil
	}
}

func Verify(token string, privateKeyText string, pwd string, salt string, expired int64) error {
	return verify(token, privateKeyText, pwd, salt, time.Duration(expired)*time.Nanosecond)
}

func verify(token string, privateKeyText string, pwd string, salt string, expired time.Duration) error {
	privateKey := rsakit.ParsePrivateKey(privateKeyText)
	plainText, err := rsakit.DecryptFromBase64(*privateKey, token)
	if err != nil {
		return err
	}

	t, err := unmarshal(plainText)
	if err != nil {
		return err
	}

	if t.Pwd != pwd || t.Salt != salt {
		return errors.New("unauthorized")
	}
	past := timekit.TimeStampNanoSecond() - t.Time
	if past > expired.Nanoseconds() {
		return errors.New("request expired")
	}
	return nil
}

func SignNow(publicKeyText string, password string, salt string) (string, error) {
	now := timekit.TimeStampNanoSecond()
	return Sign(publicKeyText, password, salt, now)
}

func VerifyNow(token string, privateKeyText string, password string, salt string) error {
	return verify(token, privateKeyText, password, salt, DEFAULT_EXPIRED)
}
