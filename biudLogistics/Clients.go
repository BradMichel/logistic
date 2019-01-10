package biudLogistics

import (
  "time"
  "log"

)

func (clients *Clients) Get(content *Content)  {  //  Obtiene los clientes de la variable content que agrupa los datos enviados por el usuario
  //TAG:="(clients *Clients) Get(content *Content)"
  //logPrintln(TAG)

  distances:=(*content).Distances //  Obtiene la matriz de distacias
  *clients=make(Clients)  //  Inicializa la variable clientes
  lenTimeWindows:=len((*content).TimeWindows)
  i:=0
  for id,demand := range (*content).Demand {  //  Recorre el vector demandas
    if demand>0 {  //  Verifica si la demanda es positiva la cual indica que es una estacion de tipo cliente
      timeInit:=time.Time{}
      timeVisited:=MyTime{Time:timeInit}
      var timeWindow MyTime
      if i<lenTimeWindows {
      timeWindow=(*content).TimeWindows[i]
      }else{
        timeInit2:=time.Time{}
        timeWindow=MyTime{Time:timeInit2}
      }

      (*clients)[id]=&Client{Station:Station{Id:id,Distances:distances[id],TimeVisited:timeVisited},Demand:demand,TimeLimit:timeWindow}  //  Crea al cliente y lo agrega al listado de clientes

      (*clients)[id].OrderDistanceAsc() //  Crea el vector ordenado de estaciones vecinas de la mas cercana a la mas lejana
      //log.Println("id",id)
      i++
    }

  }
}
func (clients *Clients)GetNoVisited()  {
  TAG:="(clients *Clients)GetNoVisited()"
  log.Println(TAG)
  con:=0
  for _,client := range *clients {
    if client.Visited==false{
      con++
    }
  }
  log.Println("# clients no visited",con)
}
func (clients *Clients)Log()  {
  TAG:="(clients *Clients)Log()"
  log.Println(TAG)
  log.Println("len",len(*clients))
  log.Println("clients",clients)

  for id,client := range *clients {
    log.Println("id",id)
    log.Printf("%+v\n",client)
  }

}
