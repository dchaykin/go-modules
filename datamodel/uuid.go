package datamodel

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/dchaykin/go-modules/database"
	"github.com/dchaykin/go-modules/log"
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

func EnsureUUID(domainEntity database.DomainEntity) error {
	uuid := domainEntity.UUID()
	if len(uuid) > 0 && len(uuid) != 32 {
		log.Info("invalid uuid: %s. A new value will be generated", uuid)
	} else if len(uuid) == 32 {
		return nil
	}

	uuid, err := GenerateUUID()
	if err != nil {
		return fmt.Errorf("could not generate a uuid: %v", err)
	}
	domainEntity.SetUUID(uuid)
	return nil
}
