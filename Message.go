package gohttp

import (
	"net/http"
	"strings"
)

const (
	protocal_version_1_0 = "1.0"
	protocal_version_1_1 = "1.1"
	protocal_version_2   = "2.0"
)

type Message struct {
	Headers     http.Header
	HeaderNames []string
	protocal    string
}

func (m *Message) GetProtocolVersion() string {
	if "" == m.protocal {
		return protocal_version_1_1
	}
	return m.protocal
}

func (m *Message) WithProtocolVersion(version string) *Message {
	newM := *m
	newM.protocal = version
	return &newM
}

func (m *Message) GetHeaders() http.Header {
	return m.Headers
}

func (m *Message) HasHeader(header string) bool {
	return nil != m.Headers[header] || nil != m.Headers[strings.ToLower(header)]
}

func (m *Message) GetHeader(header string) []string {
	if v, ok := m.Headers[header]; ok {
		return v
	}
	if v, ok := m.Headers[strings.ToLower(header)]; ok {
		return v
	}
	return []string{}
}

func (m *Message) WithHeader(header, value string) {

}