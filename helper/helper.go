package helper

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"golang.org/x/net/html"

	"github.com/dchaykin/go-modules/log"
)

func LoadAccessData(fileName string) {
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

func ValueString(fields map[string]any, fieldName string) string {
	value, ok := fields[fieldName]
	if !ok || value == nil {
		return ""
	}
	return fmt.Sprintf("%s", value)
}

func FloatFromString(value string) float64 {
	if value == "" {
		return 0
	}
	result, err := strconv.ParseFloat(value, 64)
	if err == nil {
		return result
	}
	log.Errorf("Could not parse %s into float: %v", value, err)
	return 0
}

func Int64FromString(value string) int64 {
	if value == "" {
		return 0
	}
	result, err := strconv.ParseInt(value, 10, 64)
	if err == nil {
		return result
	}
	log.Errorf("Could not parse %s into int64: %v", value, err)
	return 0
}

func HtmlToText(htmlInput string) string {
	doc, err := html.Parse(strings.NewReader(htmlInput))
	if err != nil {
		log.Error(log.WrapError(err))
		return htmlInput
	}
	var sb strings.Builder
	extractText(doc, &sb)
	return sb.String()
}

func extractText(n *html.Node, sb *strings.Builder) {
	if n.Type == html.TextNode {
		sb.WriteString(n.Data)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		extractText(c, sb)
	}
}
