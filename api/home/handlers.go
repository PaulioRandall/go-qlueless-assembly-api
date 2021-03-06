package home

import (
	"net/http"

	uhttp "github.com/PaulioRandall/go-cookies/uhttp"
	wrapped "github.com/PaulioRandall/go-qlueless-api/shared/wrapped"
	writers "github.com/PaulioRandall/go-qlueless-api/shared/writers"
)

// HomeHandler handles requests to the root path and requests to nothing (404s)
func HomeHandler(res http.ResponseWriter, req *http.Request) {
	uhttp.LogRequest(req)
	notFound(&res, req)
}

// notFound handles requests nothing (404s)
func notFound(res *http.ResponseWriter, req *http.Request) {
	r := wrapped.WrappedReply{
		Message: "Resource not found",
	}

	writers.WriteWrappedReply(res, req, http.StatusNotFound, r)
}
