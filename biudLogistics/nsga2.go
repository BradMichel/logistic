package biudLogistics

import(
  "github.com/gorilla/context"
  "net/http"
  "log"
  "time"
)

func PostNsga2(w http.ResponseWriter, r *http.Request)  {
  start := time.Now()
  TAG:="PostNsga2(w http.ResponseWriter, r *http.Request)"
  var(  //  Declaracion de variables
    deposits Deposits  //  Vector de depositos
    clients Clients //  Clientes
    vehicles Vehicles //  Vehiculos
    routes SliceRouteList //  Vector de soluciones
    penalties [2*N][]bool
    err error
  )
  routes=make(SliceRouteList, 2*N)  //  Inicializacion del vector de 2N soluciones
  content :=context.Get(r,"body").(*Content)  // Obtiene los datos enviados del Usuario

  deposits.Get(content)  //  Obtiene los depositos enviados por el Usuario
  clients.Get(content)  //  Obtiene los clientes enviados por el Usuario
  vehicles.Get(content) //  Obtiene los vehiculos enviados por el Usuario

  if penalties,err=routes.GetInitialPoblation(&vehicles, &deposits, &clients); err!=nil{
    log.Println(err)
    return
  }

  routes.LoadPaths()
  for i := 0; i < R; i++ {
    log.Println("R :",i+1," of ",R)
    if err=routes.GetSons(&penalties,&vehicles, &deposits, &clients);err!=nil{
      log.Println(err)
      return
    }
    log.Println("GetSons")
    fronts,err := routes.Dominance()
    if err!=nil{
      log.Println(err)
      return
    }
    log.Println("Dominance")
    if err=routes.GetNextGeneration(fronts);err!=nil{
      log.Println(err)
      return
    }
    log.Println("GetNextGeneration")
    routes.LoadPaths()
    log.Println("LoadPaths")
  }
    //log.Printf("%+v\n",routes)

  if err:=routes.ValidateClientsRepeat(TAG);err!=nil{
    log.Println(err)
  }

  elapsed := time.Since(start)
  log.Println("time execution: ",elapsed)
  responseJson(w, 202, routes[:N])  //  Se response al usuario enviando las N rutas resultantes

}
