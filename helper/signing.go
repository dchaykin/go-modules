package helper

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"hash"
	"io"
	"strings"
)

// SignInput signiert beliebige Eingaben mit optionalem Secret (HMAC-SHA256 oder SHA256).
//
// - input: []byte, string, *bytes.Buffer oder io.Reader
// - secret: wenn nil oder leer, wird SHA256 ohne HMAC verwendet
func SignInput(input any, secret []byte) (string, error) {
	reader, err := buildReader(input)
	if err != nil {
		return "", err
	}

	var hash hash.Hash
	if len(secret) > 0 {
		hash = hmac.New(sha256.New, secret)
	} else {
		hash = sha256.New()
	}

	if _, err := io.Copy(hash, reader); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func buildReader(input any) (io.Reader, error) {
	var reader io.Reader

	switch v := input.(type) {
	case []byte:
		reader = bytes.NewReader(v)
	case string:
		reader = strings.NewReader(v)
	case *bytes.Buffer:
		// wichtig: nicht verbrauchen
		reader = bytes.NewReader(v.Bytes())
	case io.Reader:
		reader = v
	default:
		return nil, errors.New("unsupported input type")
	}

	return reader, nil
}
