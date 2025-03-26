package aheadinator

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type Config struct {
	HeaderCapture string `json:"headers"`
}

type headerDetails struct {
	OriginalFrom  string
	LowerCaseFrom string
	To            string
}

type Aheadinator struct {
	next          http.Handler
	config        *Config
	name          string
	headerNameMap map[string]headerDetails
}

type responseWriter struct {
	http.ResponseWriter
	headers http.Header
}

func (rw *responseWriter) WriteHeader(code int) {
	for key, values := range rw.ResponseWriter.Header() {
		rw.headers[key] = values
	}

	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Header() http.Header {
	return rw.ResponseWriter.Header()
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	return rw.ResponseWriter.Write(b)
}

func (p *Aheadinator) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	rpw := &responseWriter{
		ResponseWriter: rw,
		headers:        make(http.Header),
	}

	p.next.ServeHTTP(rpw, req)

	log.Printf("%+v\n", p.headerNameMap)

	for header, v := range req.Header {
		log.Printf("matchan headeran %s\n\n", header)
		lHeader := strings.ToLower(header)
		if p.headerNameMap[lHeader].OriginalFrom != "" {
			log.Println("match")
			if p.headerNameMap[lHeader].To != "" {
				rw.Header().Set(p.headerNameMap[lHeader].To, v[0])
			} else {
				rw.Header().Set(p.headerNameMap[lHeader].OriginalFrom, v[0])
			}
		}
	}
}

func CreateConfig() *Config {
	return &Config{}
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	headerNameMap := make(map[string]headerDetails)

	for _, v := range strings.Split(strings.ReplaceAll(config.HeaderCapture, " ", ""), ",") {
		// [origName, newName]
		s := strings.Split(v, ":")
		switch len(s) {
		case 0:
			continue
		case 1:
			headerNameMap[strings.ToLower(s[0])] = headerDetails{
				OriginalFrom:  s[0],
				LowerCaseFrom: s[0],
			}
		case 2:
			headerNameMap[s[0]] = headerDetails{
				OriginalFrom:  s[0],
				LowerCaseFrom: strings.ToLower(s[0]),
				To:            s[1],
			}
		default:
			return nil, fmt.Errorf("header map is does not match pattern name:new_name")
		}
	}

	return &Aheadinator{
		next:          next,
		config:        config,
		name:          name,
		headerNameMap: headerNameMap,
	}, nil
}
