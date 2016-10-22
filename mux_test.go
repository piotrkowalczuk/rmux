package rmux_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/piotrkowalczuk/rmux"

	"io/ioutil"
)

func ExampleNewServeMux() {
	mux := rmux.NewServeMux(rmux.ServeMuxOpts{})
	mux.Handle("GET/user/deactivate", http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusForbidden)
	}))
	mux.Handle("GET/user/:id", http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		id := rmux.Params(r).Path.Get("id")

		rw.WriteHeader(http.StatusOK)
		io.WriteString(rw, `{"id": `+id+`}`)
	}))

	ts := httptest.NewServer(mux)

	var (
		res *http.Response
		err error
		pay []byte
	)

	if res, err = http.Get(ts.URL + "/user/9000"); err == nil {
		defer res.Body.Close()
		if pay, err = ioutil.ReadAll(res.Body); err == nil {
			fmt.Println(string(pay))
		}
	}
	if err != nil {
		fmt.Println(err.Error())
	}

	// Output:
	// {"id": 9000}
}

func TestServeMux_ServeHTTP(t *testing.T) {
	sm := rmux.NewServeMux(rmux.ServeMuxOpts{
		NotFound: http.NotFoundHandler(),
	})
	for pattern, given := range testPaths {
		sm.Handle(pattern, func(pat, exp string) http.Handler {
			return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
				if strings.Contains(pat, ":") {
					val := rmux.Params(r)
					if val.Path == nil {
						t.Error("context should not be empty")
					}
				}
				if r.URL.Path != exp {
					t.Errorf("executed handler do not match expected path, expected %s, got %s", exp, r.URL.Path)
				} else {
					t.Logf("proper handler executed for path: %s", exp)
				}
			})
		}(pattern, given))
	}

	ts := httptest.NewServer(sm)

	for pattern, path := range testPaths {
		t.Run(pattern, func(t *testing.T) {
			resp, err := http.Get(ts.URL + path)
			if err != nil {
				t.Fatalf("unexpected error: %s", err.Error())
			}
			if resp.StatusCode != http.StatusOK {
				t.Errorf("wrong status code: expected %s but got %s", http.StatusText(http.StatusOK), http.StatusText(resp.StatusCode))
			}
		})
	}

	t.Run("not found - wrong path", func(t *testing.T) {
		resp, err := http.Get(ts.URL + "/SOMETHING-THAT-DOES-NOT-EXISTS")
		if err != nil {
			t.Fatalf("unexpected error: %s", err.Error())
		}
		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("wrong status code: expected %s but got %s", http.StatusText(http.StatusNotFound), http.StatusText(resp.StatusCode))
		}
	})
	t.Run("not found - wrong method", func(t *testing.T) {
		resp, err := http.Head(ts.URL + "/SOMETHING-THAT-DOES-NOT-EXISTS")
		if err != nil {
			t.Fatalf("unexpected error: %s", err.Error())
		}
		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("wrong status code: expected %s but got %s", http.StatusText(http.StatusNotFound), http.StatusText(resp.StatusCode))
		}
	})
}

var testPaths = map[string]string{
	"GET/a/:a/b/:b/c/:c/d/:d/e/:e/f/:f/g/:g/h/:h": "/a/a/b/b/c/c/d/d/e/e/f/f/g/g/h/h",
	"GET/":                                             "/",
	"GET/users":                                        "/users",
	"GET/comments":                                     "/comments/",
	"GET/users/cleanup":                                "/users/cleanup",
	"GET/users/:id":                                    "/users/123",
	"GET/authorizations":                               "/authorizations",
	"GET/authorizations/:id":                           "/authorizations/1",
	"POST/authorizations":                              "/authorizations",
	"DELETE/authorizations/:id":                        "/authorizations/1",
	"GET/applications/:client_id/tokens/:access_token": "/applications/1/tokens/123456789",
}

func TestServeMux_GoString(t *testing.T) {
	mux := rmux.NewServeMux(rmux.ServeMuxOpts{})
	mux.Handle("GET/user/deactivate", http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusForbidden)
	}))

	got := mux.GoString()
	expected := `{
	"GET": {
		"Resource": "GET",
		"End": false,
		"Kind": 0,
		"Handler": false,
		"NonStatic": null,
		"Static": {
			"user": {
				"Resource": "user",
				"End": false,
				"Kind": 1,
				"Handler": false,
				"NonStatic": null,
				"Static": {
					"deactivate": {
						"Resource": "deactivate",
						"End": true,
						"Kind": 1,
						"Handler": true,
						"NonStatic": null,
						"Static": {}
					}
				}
			}
		}
	}
}`

	if got != expected {
		t.Errorf("wrong output, expected:\n	%s but got:\n	%s", expected, got)
	}
}
