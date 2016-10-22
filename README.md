# rmux [![GoDoc](https://godoc.org/github.com/piotrkowalczuk/rmux?status.svg)](http://godoc.org/github.com/piotrkowalczuk/rmux)&nbsp;[![Build Status](https://travis-ci.org/piotrkowalczuk/rmux.svg)](https://travis-ci.org/piotrkowalczuk/rmux)&nbsp;[![codecov.io](https://codecov.io/github/piotrkowalczuk/rmux/coverage.svg?branch=master)](https://codecov.io/github/piotrkowalczuk/rmux?branch=master)

RESTful router, that supports path variables. Requires Go version 1.7+.

## Features

* sinatra like path parameters - `/users/:id`
* no extra dependencies
* compatible with [http.ServeMux](https://golang.org/pkg/net/http/#ServeMux)
* middleware like system using [rmux.Interceptor](https://godoc.org/github.com/piotrkowalczuk/rmux#Interceptor) interface
* context handled by [http.Request.WithContext](https://golang.org/pkg/net/http/#Request.WithContext)

## Example 

```go
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
```

## Benchmarks

### rmux __without__ context manipulation
```
BenchmarkRMUX_GithubAll               	   20000       	     85549 ns/op       	   12208 B/op  	     203 allocs/op
BenchmarkRMUX_GithubParam             	 3000000       	       434 ns/op       	      64 B/op  	       1 allocs/op
BenchmarkRMUX_GithubStatic            	 5000000       	       303 ns/op       	      32 B/op  	       1 allocs/op
BenchmarkRMUX_GPlusParam              	10000000       	       258 ns/op       	      32 B/op  	       1 allocs/op
BenchmarkRMUX_GPlusStatic             	20000000       	       129 ns/op       	      16 B/op  	       1 allocs/op
BenchmarkRMUX_Param                   	10000000       	       217 ns/op       	      32 B/op  	       1 allocs/op
BenchmarkRMUX_Param20                 	 1000000       	      1170 ns/op       	     320 B/op  	       1 allocs/op
BenchmarkRMUX_Param5                  	 5000000       	       377 ns/op       	      80 B/op  	       1 allocs/op
BenchmarkRMUX_ParamWrite              	 1000000       	      2514 ns/op       	     880 B/op  	       8 allocs/op
```

### rmux __with__ context manipulation

```
BenchmarkRMUX_GithubAll               	   10000       	    266757 ns/op       	   98224 B/op  	    1182 allocs/op
BenchmarkRMUX_GithubParam             	 1000000       	      1173 ns/op       	     496 B/op  	       6 allocs/op
BenchmarkRMUX_GithubStatic            	 2000000       	       818 ns/op       	     400 B/op  	       5 allocs/op
BenchmarkRMUX_GPlus2Params            	 1000000       	      1061 ns/op       	     496 B/op  	       6 allocs/op
BenchmarkRMUX_GPlusParam              	 2000000       	       865 ns/op       	     432 B/op  	       6 allocs/op
BenchmarkRMUX_GPlusStatic             	 2000000       	       653 ns/op       	     384 B/op  	       5 allocs/op
BenchmarkRMUX_Param                   	 2000000       	       999 ns/op       	     432 B/op  	       6 allocs/op
BenchmarkRMUX_Param20                 	 1000000       	      2451 ns/op       	    1008 B/op  	       6 allocs/op
BenchmarkRMUX_Param5                  	 1000000       	      1438 ns/op       	     528 B/op  	       6 allocs/op
BenchmarkRMUX_ParamWrite              	 1000000       	      1322 ns/op       	     432 B/op  	       6 allocs/op
```

### httprouter

```
BenchmarkHttpRouter_GithubAll          	   30000       	     49574 ns/op       	   13792 B/op  	     167 allocs/op
BenchmarkHttpRouter_GithubParam        	 5000000       	       289 ns/op       	      96 B/op  	       1 allocs/op
BenchmarkHttpRouter_GithubStatic       	20000000       	        62.6 ns/op     	       0 B/op  	       0 allocs/op
BenchmarkHttpRouter_GPlus2Params       	10000000       	       215 ns/op       	      64 B/op  	       1 allocs/op
BenchmarkHttpRouter_GPlusParam         	10000000       	       225 ns/op       	      64 B/op  	       1 allocs/op
BenchmarkHttpRouter_GPlusStatic        	50000000       	        36.9 ns/op     	       0 B/op  	       0 allocs/op
BenchmarkHttpRouter_Param              	20000000       	       124 ns/op       	      32 B/op  	       1 allocs/op
BenchmarkHttpRouter_Param20            	 1000000       	      1208 ns/op       	     640 B/op  	       1 allocs/op
BenchmarkHttpRouter_Param5             	 5000000       	       409 ns/op       	     160 B/op  	       1 allocs/op
BenchmarkHttpRouter_ParamWrite         	10000000       	       165 ns/op       	      32 B/op  	       1 allocs/op
```

### httptreemux

```
BenchmarkHttpTreeMux_GithubAll         	   10000       	    185566 ns/op       	   65856 B/op  	     671 allocs/op
BenchmarkHttpTreeMux_GithubParam       	 1000000       	      1086 ns/op       	     384 B/op  	       4 allocs/op
BenchmarkHttpTreeMux_GithubStatic      	20000000       	        66.7 ns/op     	       0 B/op  	       0 allocs/op
BenchmarkHttpTreeMux_GPlus2Params      	 2000000       	       969 ns/op       	     384 B/op  	       4 allocs/op
BenchmarkHttpTreeMux_GPlusParam        	 2000000       	       736 ns/op       	     352 B/op  	       3 allocs/op
BenchmarkHttpTreeMux_GPlusStatic       	50000000       	        39.7 ns/op     	       0 B/op  	       0 allocs/op
BenchmarkHttpTreeMux_Param             	 2000000       	       835 ns/op       	     352 B/op  	       3 allocs/op
BenchmarkHttpTreeMux_Param20           	  200000       	      9043 ns/op       	    3196 B/op  	      10 allocs/op
BenchmarkHttpTreeMux_Param5            	 1000000       	      1664 ns/op       	     576 B/op  	       6 allocs/op
BenchmarkHttpTreeMux_ParamWrite        	 2000000       	       812 ns/op       	     352 B/op  	       3 allocs/op
```