package httpcomm

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/dchaykin/go-modules/auth"
	"github.com/dchaykin/go-modules/log"
	"github.com/stretchr/testify/require"
)

func loadAccessData(fileName string) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	accessData := map[string]string{}
	if err = json.Unmarshal(data, &accessData); err != nil {
		panic(err)
	}

	for k, v := range accessData {
		os.Setenv(k, v)
	}
}

func TestDownloadFile(t *testing.T) {
	loadAccessData("../.do-not-commit/env.vars")
	log.SetLevel(log.LevelDebug)

	user := auth.GetTestUserIdentity()

	md, err := DownloadFile("6887368bb26123efc0d0840ec3db3d94", "/tmp", user)
	require.NoError(t, err)

	require.FileExists(t, "/tmp/"+md.OriginalFileName)
	os.Remove("/tmp/" + md.OriginalFileName)
	require.NoFileExists(t, "/tmp/"+md.OriginalFileName)

	require.NotNil(t, md)
	require.Equal(t, "Invoice-E01KUPOQ-0001.pdf", md.OriginalFileName)
}

func TestUploadFile(t *testing.T) {
	loadAccessData("../.do-not-commit/env.vars")
	log.SetLevel(log.LevelDebug)

	user := auth.GetTestUserIdentity()

	fileUUID, md, err := UploadFile(os.TempDir()+"inquiry/IntroductionToSemiconductorModule.pdf", user)
	require.NoError(t, err)

	require.NotNil(t, md)
	require.NotEmpty(t, fileUUID)
	require.Equal(t, "IntroductionToSemiconductorModule.pdf", md.OriginalFileName)
	require.Equal(t, md.FileSize, int64(2248586))
}
