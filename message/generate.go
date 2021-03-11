// +build ignore

package main

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"go/format"
	"html/template"
	"io"
	"net/http"
	"os"
	"strings"
)

const faceDownloadUrl = `https://downv6.qq.com/qqface/config/face_config_8.5.0.15.zip?mType=Other` //? 好像是会自动更新的

type config struct {
	SystemFace []face `json:"sysface"`
}

type face struct {
	QSid string `json:"QSid"`
	QDes string `json:"QDes"`
}

const codeTemplate = `// Code generated by message/generate.go DO NOT EDIT.

package message

var faceMap = map[int]string{
{{range .SystemFace}}	{{.QSid}}:	"{{.QDes}}",
{{end}}
}
`

func main() {
	f, _ := os.OpenFile("face.go", os.O_WRONLY|os.O_CREATE|os.O_SYNC|os.O_TRUNC, 0755)
	defer func() { _ = f.Close() }()
	resp, err := http.Get(faceDownloadUrl)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	rsp, _ := io.ReadAll(resp.Body)
	reader, _ := zip.NewReader(bytes.NewReader(rsp), resp.ContentLength)
	file, _ := reader.Open("face_config.json")
	data, _ := io.ReadAll(file)
	faceConfig := config{}
	_ = json.Unmarshal(data, &faceConfig)
	for i := range faceConfig.SystemFace {
		faceConfig.SystemFace[i].QDes = strings.TrimPrefix(faceConfig.SystemFace[i].QDes, "/")
	}
	tmpl, _ := template.New("template").Parse(codeTemplate)
	buffer := &bytes.Buffer{}
	_ = tmpl.Execute(buffer, &faceConfig)
	source, _ := format.Source(buffer.Bytes())
	f.Write(source)
}