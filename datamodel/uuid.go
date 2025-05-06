package datamodel

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"time"
)

func GenerateUUID() (string, error) {
	timestamp := time.Now().Unix()
	timestampHex := fmt.Sprintf("%08x", timestamp)
	randomBytes := make([]byte, 12)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}
	randomHex := hex.EncodeToString(randomBytes)
	return timestampHex + randomHex, nil
}

func ExtractTimeFromUUID(uuid string) (time.Time, error) {
	if len(uuid) < 8 {
		return time.Time{}, errors.New("UUID is too short for fetching the timestamp")
	}
	timestampHex := uuid[:8]
	unixSeconds, err := strconv.ParseInt(timestampHex, 16, 64)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid Hex-Timestamp: %w", err)
	}
	return time.Unix(unixSeconds, 0), nil
}
