package rmux

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
)

const (
	nodeKindMethod = iota
	nodeKindStatic
	nodeKindComplete
	nodeKindNonStatic
)

// ServeMux is an HTTP request multiplexer.
// It matches the URL of each incoming request
// against a list of registered patterns and calls the handler for the pattern that most closely matches the URL.
//
// It support RESTful naming convention:
// 	/users
//	/users/action
// 	/users/:id
// 	/users/:id/comments
// 	/users/:id/comments/hide
type ServeMux struct {
	methods  map[string]*node
	notFound http.Handler
	inter    Interceptor
	ctx      context.Context
}

// ServeMuxOpts allow to pass extra options to the muxer.
type ServeMuxOpts struct {
	NotFound    http.Handler
	Interceptor Interceptor
	Context     context.Context
}

// Interceptor is a function that decorate handler with custom logic.
type Interceptor func(http.ResponseWriter, *http.Request, http.Handler)

// NewServeMux allocates and returns a new ServeMux with default not found handler and context.
func NewServeMux(opts ServeMuxOpts) *ServeMux {
	if opts.NotFound == nil {
		opts.NotFound = http.NotFoundHandler()
	}
	if opts.Context == nil {
		opts.Context = context.Background()
	}
	sm := &ServeMux{
		methods:  map[string]*node{},
		notFound: opts.NotFound,
		inter:    opts.Interceptor,
		ctx:      opts.Context,
	}

	return sm
}

// Handle registers the handler for the given pattern.
func (sm *ServeMux) Handle(p string, h http.Handler) {
	args := sm.split(p)

	if _, ok := sm.methods[args.vals[0]]; !ok {
		sm.methods[args.vals[0]] = newNode(args.vals[0], nodeKindMethod)
	}

	sm.methods[args.vals[0]].add(args.vals[1:], h, 0)
	return

}

// ServeHTTP dispatches the request to the handler whose pattern most closely matches the request URL.
func (sm *ServeMux) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	args := sm.split(r.URL.Path)

	var (
		h  http.Handler
		ok bool
	)

	if root := sm.methods[r.Method]; root != nil {
		if !ok {
			h, ok = sm.search(root, &args)
		}
		if !ok {
			sm.notFound.ServeHTTP(rw, r.WithContext(sm.ctx))
			return
		}

		if args.len > 0 {
			r = r.WithContext(newContext(sm.ctx, args))
		} else {
			r = r.WithContext(sm.ctx)
		}
		if sm.inter != nil {
			sm.inter(rw, r, h)
		} else {
			h.ServeHTTP(rw, r)
		}

		return
	}

	sm.notFound.ServeHTTP(rw, r.WithContext(sm.ctx))
}

func (sm *ServeMux) search(n *node, args *arguments) (http.Handler, bool) {
	if args.len == 0 {
		return nil, false
	}

	var (
		nn *node
		ok bool
	)

	for idx, res := range (*args).vals {
		if n.static != nil {
			// Find static or non-static node.
			if nn, ok = n.static[res]; ok {
				if nn.end && args.len-1 == idx {
					if nn.res == res {
						return nn.handler, true
					}

					return nil, false
				}
			}
		}

		if n.nonstatic != nil {
			ok = true
			nn = n.nonstatic
			// slice of keys, values are stored in path variable
			if args.keys == nil {
				args.keys = make([]string, args.len)
			}
			args.keys[idx] = nn.res

			if nn.end && args.len-1 == idx {
				return nn.handler, true
			}
		}

		if !ok {
			return nil, false
		}

		n = nn
	}

	return nil, false
}

func (sm *ServeMux) split(p string) arguments {
	if p == "" {
		return arguments{}
	}

	if strings.HasPrefix(p, "/") {
		p = p[1:]
	}

	if strings.HasSuffix(p, "/") {
		p = p[:len(p)-1]
	}

	vals := strings.Split(p, "/")
	return arguments{
		vals: vals,
		len:  len(vals),
	}
}

// GoString implements fmt GoStringer interface.
func (sm *ServeMux) GoString() string {
	b, err := json.MarshalIndent(sm.methods, "", "\t")
	if err != nil {
		return err.Error()
	}

	return string(b)
}

type node struct {
	res       string
	end       bool
	kind      int
	handler   http.Handler
	static    map[string]*node
	nonstatic *node
}

func (n *node) MarshalJSON() ([]byte, error) {
	d := struct {
		Resource  string
		End       bool
		Kind      int
		Handler   bool
		NonStatic *node
		Static    map[string]*node
	}{
		Resource:  n.res,
		End:       n.end,
		Kind:      n.kind,
		NonStatic: n.nonstatic,
		Static:    n.static,
		Handler:   n.handler != nil,
	}

	return json.Marshal(d)
}

func newNode(res string, kind int) *node {
	return &node{
		res:    strings.TrimLeft(res, ":"),
		kind:   kind,
		static: map[string]*node{},
	}
}

func (n *node) add(path []string, h http.Handler, idx int) *node {
	if n == nil {
		n = &node{}
	}

	// for example GET/
	if len(path) == 0 && idx == 0 {
		nn := &node{
			end:     true,
			handler: h,
			kind:    nodeKindStatic,
		}
		n.static[""] = nn
		return nn
	}

	if len(path) == 0 || idx >= len(path) {
		return nil
	}

	var (
		ok  bool
		nn  *node
		res string
	)

	if len(path) > 0 {
		res = path[idx]
	}

	if isStatic(res) {
		if n.static == nil {
			n.static = map[string]*node{}
		}

		if nn, ok = n.static[res]; !ok {
			nn = newNode(res, nodeKindStatic)
		}

		if idx == len(path)-1 {
			nn.handler = h
			nn.end = true
		}

		if !ok {
			n.static[res] = nn
		}

		return nn.add(path, h, idx+1)
	}

	if n.nonstatic == nil {
		n.nonstatic = newNode(res, nodeKindNonStatic)
	}

	if idx == len(path)-1 {
		n.nonstatic.handler = h
		n.nonstatic.end = true
	}

	return n.nonstatic.add(path, h, idx+1)

}

func isStatic(p string) bool {
	return !strings.HasPrefix(p, ":")
}
