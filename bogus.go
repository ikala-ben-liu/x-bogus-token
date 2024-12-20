// Package bogus Package plugindemo a demo plugin.
package bogus

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"
)

// Config the plugin configuration.
type Config struct {
	Headers map[string]string `json:"headers,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		Headers: make(map[string]string),
	}
}

// Demo a Demo plugin.
type Demo struct {
	next     http.Handler
	headers  map[string]string
	name     string
	template *template.Template
}

// New created a new Demo plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if len(config.Headers) == 0 {
		return nil, fmt.Errorf("headers cannot be empty")
	}

	return &Demo{
		headers:  config.Headers,
		next:     next,
		name:     name,
		template: template.New("demo").Delims("[[", "]]"),
	}, nil
}

func (a *Demo) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	for key, value := range a.headers {
		tmpl, err := a.template.Parse(value)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		writer := &bytes.Buffer{}

		err = tmpl.Execute(writer, req)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		req.Header.Set(key, writer.String())
	}
	// 如果 request body 為空，則不處理
	if req.Body == nil {
		return
	}
	// 讀取 request body
	bodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	bodyString := string(bodyBytes)

	log.Println("bodyBytes: %v", bodyString)
	// 將 bodyString 轉換為 base64 字串
	encodedBody := base64.StdEncoding.EncodeToString([]byte(bodyString))
	log.Println("encodedBody: %v", encodedBody)
	req.Header.Set("Encoded-Body", encodedBody)
	a.next.ServeHTTP(rw, req)
}
