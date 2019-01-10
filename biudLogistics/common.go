package biudLogistics

import(
	"net/http"
	"fmt"
	"encoding/json"
)


func responseJson(w http.ResponseWriter , errorCode int,data interface{}){
	// //logPrintln("responseJson")
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w,)
		http.Error(w,"",errorCode)
		jsonData,_ := json.Marshal(data)
		w.Write(jsonData)
}
