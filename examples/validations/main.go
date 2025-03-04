package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/go-msvc/errors/v2"
	"github.com/go-msvc/errors/v2/examples/validations/users"
)

func main() {
	portFlag := flag.Int("p", 8080, "HTTP port")
	flag.Parse()
	http.HandleFunc("/add", addUser)
	http.HandleFunc("/upd", updUser)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *portFlag), nil); err != nil {
		panic(err)
	}
}

func addUser(httpRes http.ResponseWriter, httpReq *http.Request) {
	if httpReq.Method != http.MethodPost {
		http.Error(httpRes, "this is not a post", http.StatusMethodNotAllowed)
		return
	}
	if httpReq.Body == nil {
		http.Error(httpRes, "missing body", http.StatusBadRequest)
		return
	}
	var req users.AddUserRequest
	if err := json.NewDecoder(httpReq.Body).Decode(&req); err != nil {
		http.Error(httpRes, fmt.Sprintf("cannot parse JSON body into %T: %+s", req, err), http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		err = errors.Wrap(err, "invalid request")
		fmt.Fprintf(os.Stderr, "HTTP %s %s: %+v\n",
			httpReq.Method,
			httpReq.URL.Path,
			err) //note use of %+v in log, but %+s below in user message
		http.Error(httpRes, fmt.Sprintf("%+s", err), http.StatusBadRequest)
		return
	}
}

func updUser(httpRes http.ResponseWriter, httpReq *http.Request) {
	if httpReq.Method != http.MethodPut {
		http.Error(httpRes, "this is not a put", http.StatusMethodNotAllowed)
		return
	}
	if httpReq.Body == nil {
		http.Error(httpRes, "missing body", http.StatusBadRequest)
		return
	}
	var req users.UpdateUserRequest
	if err := json.NewDecoder(httpReq.Body).Decode(&req); err != nil {
		http.Error(httpRes, fmt.Sprintf("cannot parse JSON body into %T: %+s", req, err), http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		err = errors.Wrap(err, "invalid request")
		fmt.Fprintf(os.Stderr, "HTTP %s %s: %+v\n",
			httpReq.Method,
			httpReq.URL.Path,
			err) //note use of %+v in log, but %+s below in user message
		http.Error(httpRes, fmt.Sprintf("%+s", err), http.StatusBadRequest)
		return
	}
}
