package rmux

import (
	"context"
	"net/http"
	"net/url"
)

const (
	contextKey = "rmux_arguments"
)

type arguments struct {
	keys []string
	vals []string
	len  int
}

func newContext(ctx context.Context, args arguments) context.Context {
	return context.WithValue(ctx, contextKey, args)
}

func fromContext(ctx context.Context) url.Values {
	v, ok := ctx.Value(contextKey).(arguments)
	if !ok {
		return nil
	}
	res := make(url.Values, len(v.keys))
	for i, k := range v.keys {
		if res[k] == nil {
			res[k] = make([]string, 0, 1)
		}
		res[k] = append(res[k], v.vals[i])
	}
	return res
}

// Values is a handy container for url.Values.
type Values struct {
	Path  url.Values
	Query url.Values
}

// Params allocates new Values based on given http Request object.
func Params(r *http.Request) *Values {
	return &Values{
		Path:  fromContext(r.Context()),
		Query: r.URL.Query(),
	}
}
