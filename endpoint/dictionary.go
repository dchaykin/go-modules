package endpoint

import (
	"fmt"
	"net/http"
	"os"

	"github.com/dchaykin/go-modules/log"
	"github.com/gorilla/mux"
)

func DownloadByLanguage(w http.ResponseWriter, r *http.Request, path string) {
	vars := mux.Vars(r)
	language := vars["language"]

	dictionaryFile := fmt.Sprintf("%s/%s.csv", path, language)
	content, err := os.ReadFile(dictionaryFile)
	if err != nil {
		log.Warn("Unable to read the dictionary file %s: %v ", dictionaryFile, err)
	}

	w.Header().Set("Content-Type", "text/csv")
	w.WriteHeader(http.StatusOK)
	w.Write(content)
}
