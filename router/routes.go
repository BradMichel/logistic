package router
import(
	"net/http"
	"github.com/justinas/alice"
	"github.com/gorilla/context"
	biudLogistics "biudLogistics"
	"log"
)

type Route struct {
	Method 		string
	Pattern 	string
	HandlerFunc http.HandlerFunc
	Content v
}
type v interface{}
type Routes []Route

var routes = Routes{
/*	Route{
		"POST",
		"/api/bee/contents",
		biudLogistics.PostBee,
		biudLogistics.Content{},
	},
	Route{
		"POST",
		"/api/ecwa/contents",
		biudLogistics.PostEcwa,
		biudLogistics.Content{},
	},
*/
Route{
	"POST","/api/nsga2/contents",biudLogistics.PostNsga2,biudLogistics.Content{},
},
Route{
	"POST","/api/timeWindows/contents",biudLogistics.PostTimeWindows,biudLogistics.Content{},
},
/*	Route{
		"GET",
		"/contents/:K",
		getContent,
	},
	Route{
		"GET",
		"/contents",
		getContent,
	},*/
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile )
	commonHandlers := alice.New(context.ClearHandler,loggingHandler, recoverHandler)

	router := NewRouter()


	for _, route := range routes{
		//logPrintln(route)
		switch route.Method {
			case "GET", "DELETE":
				router.handle(route.Method,route.Pattern, commonHandlers.ThenFunc(route.HandlerFunc))

			case "POST", "PUT":
				router.handle(route.Method,route.Pattern, commonHandlers.Append(contentTypeHandler, bodyHandler(route.Content)).ThenFunc(route.HandlerFunc))
		}

	}


	http.Handle("/", router)

}
