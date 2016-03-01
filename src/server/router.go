package server

import (
	"net/http"
	"github.com/gorilla/mux"
	"log"
	"time"
)

type Route struct {
    Name        string
    Method      string
    Pattern     string
    HandlerFunc http.HandlerFunc
}

type Routes []Route


func NewRouter() *mux.Router {
router := mux.NewRouter().StrictSlash(true)
    for _, route := range routes {
		handler := Logger(route.HandlerFunc, route.Name)
        router.
            Methods(route.Method).
            Path(route.Pattern).
            Name(route.Name).
            Handler(handler)
    }

    return router

}


var routes = Routes{
    Route{
        "RecordDelete",
        "DELETE",
        "/",
        RecordDelete,
    },
	Route{
        "RecordCreate",
        "POST",
        "/",
        RecordCreate,
    },
    Route{
        "RecordIndex",
        "GET",
        "/records",
        RecordIndex,
    },
    Route{
        "RecordShow",
        "GET",
        "/records/{recordId}",
        RecordShow,
    },
}


func Logger(inner http.Handler, name string) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()

        inner.ServeHTTP(w, r)

        log.Printf(
            "%s\t%s\t%s\t%s",
            r.Method,
            r.RequestURI,
            name,
            time.Since(start),
        )
    })
}
