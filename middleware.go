package plugin_log4shell

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

// Config the plugin configuration.
type Config struct {
	ErrorCode int `json:"errorCode"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		ErrorCode: http.StatusOK,
	}
}

// Log4J a plugin.
type Log4J struct {
	next      http.Handler
	name      string
	ErrorCode int
}

// New created a new plugin.
func New(_ context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	return &Log4J{
		name:      name,
		next:      next,
		ErrorCode: config.ErrorCode,
	}, nil
}

func (l *Log4J) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	for _, values := range req.Header {
		for _, value := range values {
			if containsJNDI(value) {
				rw.WriteHeader(l.ErrorCode)
				return
			}
		}
	}

	l.next.ServeHTTP(rw, req)
}

func containsJNDI(value string) bool {
	if len(value) < 8 {
		return false
	}

	lower := strings.ToLower(value)

	if !strings.Contains(lower, "${") {
		return false
	}

	if strings.Contains(lower, "${jndi") {
		return true
	}

	root := Parse(lower)

	for _, node := range root.Values {
		if containsJNDINode(node) {
			return true
		}
	}

	return false
}

func containsJNDINode(node fmt.Stringer) bool {
	v, ok := node.(*NodeExpression)
	if !ok {
		return false
	}

	fmt.Println(v.Key.String())
	if strings.Contains(v.Key.String(), "jndi") {
		return true
	}

	for _, k := range v.Key {
		if containsJNDINode(k) {
			return true
		}
	}

	return false
}
