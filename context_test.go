package rmux_test

import (
	"net/http"
	"testing"

	"github.com/piotrkowalczuk/rmux"
)

func TestParams(t *testing.T) {
	req, err := http.NewRequest("GET", "something", nil)
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	par := rmux.Params(req)

	if par.Path != nil {
		t.Error("expected nil")
	}
}
