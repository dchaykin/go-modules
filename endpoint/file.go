package endpoint

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"os"

	"github.com/dchaykin/go-modules/datamodel"
	"github.com/dchaykin/go-modules/user"
	"github.com/dchaykin/mygolib/httpcomm"
	"github.com/dchaykin/mygolib/log"
)

func retrieveMetaData(fileUUID string, userIdentity user.UserIdentity) (*datamodel.MetaData, error) {
	log.Debug("Downloading metadata for file %s", fileUUID)

	ep := fmt.Sprintf("https://%s/app-cloudfile/api/metadata/%s", os.Getenv("MYHOST"), fileUUID)
	hr := httpcomm.Get(ep, userIdentity, nil, nil)
	if err := hr.GetError(); err != nil {
		return nil, log.WrapError(err)
	}

	sr := httpcomm.ServiceResponse{}
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

	log.Debug("Metadata for file %s is downloaded: %s", fileUUID, string(data))

	md := datamodel.MetaData{}
	if err := json.Unmarshal(data, &md); err != nil {
		return nil, log.WrapError(err)
	}

	return &md, nil
}

func DownloadFile(fileUUID, path string, userIdentity user.UserIdentity) (*datamodel.MetaData, error) {
	md, err := retrieveMetaData(fileUUID, userIdentity)
	if err != nil {
		return nil, log.WrapError(err)
	}

	log.Debug("Downloading file %s into %s", md.OriginalFileName, path)

	ep := fmt.Sprintf("https://%s/app-cloudfile/api/file/%s", os.Getenv("MYHOST"), fileUUID)
	hr := httpcomm.Get(ep, userIdentity, nil, nil)
	if err := hr.GetError(); err != nil {
		return nil, log.WrapError(err)
	}
	if err := os.WriteFile(path+"/"+md.OriginalFileName, hr.Answer, 0644); err != nil {
		return nil, log.WrapError(err)
	}

	log.Debug("%s/%s downloaded", path, md.OriginalFileName)

	return md, nil
}

func UploadFile(pathToFile string, userIdentity user.UserIdentity) (string, *datamodel.MetaData, error) {
	log.Debug("Uploading file %s", pathToFile)
	if userIdentity == nil {
		return "", nil, log.WrapError(fmt.Errorf("userIdentity is nil, could not upload file %s", pathToFile))
	}

	file, err := os.Open(pathToFile)
	if err != nil {
		return "", nil, log.WrapError(fmt.Errorf("error opening file %s: %w", pathToFile, err))
	}
	defer file.Close()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	// File-Part hinzufügen (entspricht -F "filename=@...")
	part, err := writer.CreateFormFile("filename", file.Name())
	if err != nil {
		return "", nil, log.WrapError(fmt.Errorf("error creating form file: %w", err))
	}

	if _, err = io.Copy(part, file); err != nil {
		return "", nil, log.WrapError(fmt.Errorf("error copying file to form: %w", err))
	}

	if err := writer.Close(); err != nil {
		return "", nil, log.WrapError(fmt.Errorf("error closing writer: %w", err))
	}

	ep := fmt.Sprintf("https://%s/app-cloudfile/api/upload", os.Getenv("MYHOST"))
	hr := httpcomm.PostBuffer(ep, userIdentity, map[string]string{
		"Content-Type": writer.FormDataContentType(),
	}, &body)

	if err := hr.GetError(); err != nil {
		return "", nil, log.WrapError(err)
	}

	sr := httpcomm.ServiceResponse{}
	err = json.Unmarshal([]byte(hr.Answer), &sr)
	if err != nil {
		return "", nil, log.WrapError(err)
	}

	fileUUID := fmt.Sprintf("%v", sr.Data)
	log.Debug("%s uploaded, uuid: %v", pathToFile, fileUUID)

	md, err := retrieveMetaData(fileUUID, userIdentity)
	return fileUUID, md, err
}
