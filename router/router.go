package router

import (
	"fmt"
	"reflect"
	"encoding/json"
	"net/http"
	"github.com/julienschmidt/httprouter"
	"github.com/gorilla/context"

)


// Middlewares

func recoverHandler(next http.Handler) http.Handler {
	// //logPrintln("recoverHandler")
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				responseJson(w,500,map[string]string{"error":"Internal server error"})
				return
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func loggingHandler(next http.Handler) http.Handler {
	// //logPrintln(//"loggingHandler")
	fn := func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func contentTypeHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			responseJson(w,415,map[string]string{"messange":"Data validation failed","error":"Content type unsupported "})
			return
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func bodyHandler(v interface{}) func(http.Handler) http.Handler {

	t := reflect.TypeOf(v)

	m := func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {

			val := reflect.New(t).Interface()
			err := json.NewDecoder(r.Body).Decode(val)


			if err != nil {
				responseJson(w,406,map[string]string{"messange":"Data validation failed","error":err.Error()})
				return
			}

			if next != nil {
				context.Set(r, "body", val)
				next.ServeHTTP(w, r)
			}
		}
		return http.HandlerFunc(fn)
	}

	return m
}

// Router

type router struct {
	*httprouter.Router
}

func (r *router) handle(method string, path string, handler http.Handler) {

	r.Handle(method,path, wrapHandler(handler))
}



func NewRouter() *router {
	// //logPrintln("NewRouter")
	return &router{httprouter.New()}
}

func wrapHandler(h http.Handler) httprouter.Handle {
	// //logPrintln("wrapHandler")
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		context.Set(r, "params", ps)
		h.ServeHTTP(w, r)
	}
}

func responseJson(w http.ResponseWriter , errorCode int,data interface{}){
	// //logPrintln("responseJson")
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w,)
		http.Error(w,"",errorCode)
		errorJson, _ := json.Marshal(data)
		w.Write(errorJson)
}
