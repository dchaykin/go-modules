package httpcomm

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/dchaykin/go-modules/auth"
	"github.com/dchaykin/go-modules/log"
)

type MetaData struct {
	OriginalFileName string    `json:"filename"`
	FileSize         int64     `json:"size"`
	Author           string    `json:"author"`
	Partner          string    `json:"partner"`
	MimeType         string    `json:"mime"`
	Error            *string   `json:"error,omitempty"`
	Timestamp        time.Time `json:"createdAt"`
}

func DownloadFile(fileUUID, path string, userIdentity auth.UserIdentity) (*MetaData, error) {
	ep := fmt.Sprintf("https://%s/app-cloudfile/api/metadata/%s", os.Getenv("MYHOST"), fileUUID)
	hr := Get(ep, userIdentity, nil, nil)
	if err := hr.GetError(); err != nil {
		return nil, log.WrapError(err)
	}

	sr := ServiceResponse{}
	if err := json.Unmarshal(hr.Answer, &sr); err != nil {
		return nil, log.WrapError(err)
	}

	if sr.Error != nil {
		return nil, log.WrapError(fmt.Errorf(*sr.Error))
	}

	data, err := sr.GetPayload()
	if err != nil {
		return nil, log.WrapError(err)
	}

	md := MetaData{}
	if err := json.Unmarshal(data, &md); err != nil {
		return nil, log.WrapError(err)
	}

	ep = fmt.Sprintf("https://%s/app-cloudfile/api/file/%s", os.Getenv("MYHOST"), fileUUID)
	hr = Get(ep, userIdentity, nil, nil)
	if err := hr.GetError(); err != nil {
		return nil, log.WrapError(err)
	}
	if err := os.WriteFile(path+"/"+md.OriginalFileName, hr.Answer, 0644); err != nil {
		return nil, log.WrapError(err)
	}

	return &md, nil
}
