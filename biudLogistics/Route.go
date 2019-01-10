package biudLogistics

import (
  "time"
  "math/rand"
  "log"
  "errors"
  "strconv"
)

func (route *Route) GetNeighbor(kind string,vehicle *Vehicle,deposit *Deposit,clients *Clients)  {  //  Crea la ruta siguiendo la logica del vecino mas cercano
  //TAG:="(route *Route) NearestNeighbor(deposit *Deposit,vehicle *Vehicle,clients *Clients)"
  //log.Println(TAG)

  if deposit.Penalized==false { //  Verifica si el deposito seleccionado no esta penalizado
    route.InitLoads(vehicle, deposit)
    previousStation:=&deposit.Station
    actualStation:=new(Client)
    switch kind {
    case "random":
      actualStation=deposit.NextRandomClient(clients)
      break;
    case "nearest":
      actualStation=deposit.NextClient(clients,0)
      break;
    }

    for actualStation!=nil && (*route).Residue-actualStation.Demand>=0 {
      (*route).Residue-=actualStation.Demand
      (*route).clients=append((*route).clients,*(*clients)[actualStation.Id])
      previousStation=new(Station)
      previousStation=&actualStation.Station
      previousStation.Visited=true
      actualStation=new(Client)
      actualStation=previousStation.NextClient(clients,0)
      actualStation=route.AddLastClient(actualStation, previousStation, clients)

      }
  }}


func (route *Route)GetTimeWindows(vehicle *Vehicle,deposit *Deposit,clients *Clients)  {
  //TAG:="(route *Route)GetTimeWindows(vehicle *Vehicle,deposit *Deposit,clients *Clients)"
  //logPrintln(TAG)

  if deposit.Penalized==false{
    route.InitLoads(vehicle, deposit)
    previousStation:=&deposit.Station
    actualStation:=route.GetActualStationForTimeWindows(previousStation,clients)

    for (actualStation!=nil && (*route).Residue-(*actualStation).Demand>=0 ){
      (*route).Residue-=(*actualStation).Demand
      (*route).clients=append((*route).clients,*actualStation)
      distance:=previousStation.GetDistaceTo(&actualStation.Station)
      newTime:=previousStation.TimeVisited.Time.Add(time.Duration(distance) * time.Second)
      (*actualStation).TimeVisited=MyTime{newTime}
      log.Println(actualStation.TimeVisited)
      previousStation=&actualStation.Station
      previousStation.Visited=true
      actualStation=route.GetActualStationForTimeWindows(previousStation,clients)
    }
  }
}

func (route *Route)GetActualStationForTimeWindows(previousStation *Station,clients *Clients) (actualStation *Client) {
  //TAG:="(route *Route)GetActualStationForTimeWindows(previousStation *Station,clients *Clients) (actualStation *Client)"
  //logPrintln(TAG)

  nextClientForTimeWindow:=previousStation.NextClientForTimeWindow(clients)
  nextClient:=previousStation.NextClient(clients,0)
  if nextClientForTimeWindow!=nil {
    actualStation=nextClientForTimeWindow
    if previousStation.ValidatePreviousVisit(nextClientForTimeWindow,nextClient){
      actualStation=nextClient
    }
  }else if nextClient != nil{
    actualStation=nextClient
  }
  return
}

func (route *Route)InitLoads(vehicle *Vehicle,deposit *Deposit)  {
  //TAG:="(route *Route)InitLoads(vehicle *Vehicle,deposit *Deposit)"
  //logPrintln(TAG)
  (*route).deposit=*deposit
  (*route).Id=deposit.Id
    (*route).Residue=vehicle.Capacity // Inicializa la ruta con un residuo igual a la capacidad del vehiculo

    if deposit.Load-vehicle.Capacity < 0 {
      (*vehicle).Load=(*deposit).Load
      (*deposit).Load=0
    }else{
      (*vehicle).Load=vehicle.Capacity
      (*deposit).Load-=vehicle.Capacity
    }
}

func (route *Route)AddLastClient(actualStation *Client,previousStation *Station,clients *Clients) (newActualStation *Client) {
  if actualStation!=nil && (*route).Residue-actualStation.Demand<0 {
    for i := 0; i < len(*clients); i++  {
        if ((*route).Residue-actualStation.Demand>=0) {break}
        clientForValidate:=previousStation.NextClient(clients,i)
        if clientForValidate!=nil && (*route).Residue-clientForValidate.Demand>=0{
        actualStation=new(Client)
        actualStation=clientForValidate
      }
      }
    }
    newActualStation=new(Client)
    newActualStation=actualStation
    return
}

func (route *Route)AddPath(from *Station,to *Station)  {
  //TAG:="(route *Route)AddPath(from *Station,to *Station)"
  //logPrintln(TAG)

  distance:=from.GetDistaceTo(to)
  (*to).Visited=true
  newTime:=from.TimeVisited.Time.Add(time.Duration(distance) * time.Second)
  to.TimeVisited=MyTime{newTime}
  (*route).R=append((*route).R,Path{I:(*from).Id,J:(*to).Id,Time:distance})
  (*route).Time+=distance
}

func (route1 *Route)ExchangeLastCustomers(route2 *Route)  (success bool){
  //TAG:="(route1 *Route)ExchangeLastCustomers(route2 *Route)  (success bool)"
  //log.Println(TAG)

  positionLastClientRoute1:=len((*route1).clients)-1
  positionLastClientRoute2:=len((*route2).clients)-1

  if(positionLastClientRoute1>1 && positionLastClientRoute2>1){
    lastClienRoute1:=(*route1).clients[positionLastClientRoute1]
    lastClienRoute2:=(*route2).clients[positionLastClientRoute2]
    if (*route1).ValidateChangeClient(&lastClienRoute2,positionLastClientRoute1) && (*route2).ValidateChangeClient(&lastClienRoute1,positionLastClientRoute2) {
          (*route1).ChangeClient(lastClienRoute2,positionLastClientRoute1)
          (*route2).ChangeClient(lastClienRoute1,positionLastClientRoute2)
          success=true
    }
  }
  return
}

func (route *Route)ValidateChangeClient(clientChange *Client,clientChangePosition int)(success bool)  {
  //TAG:="(route *Route)ValidateChangeClient(clientChange *Client,clientChangePosition int)(success bool)"
  //log.Println(TAG)
  clientChanged:=(*route).clients[clientChangePosition]
  if clientChanged.Demand+(*route).Residue>=clientChange.Demand && clientChange.Id != clientChanged.Id{
    //log.Println("change",true)
    success=true
  }

  return
}

func (route *Route)ChangeClient(clientChange Client,clientChangePosition int) (success bool) {
  //TAG:="(route *Route)ChangeClient(clientChange *Client,clientChangePosition int) (success bool)"
  //log.Println(TAG)
  (*route).Residue+=clientChange.Demand
  if (*route).ValidateChangeClient(&clientChange,clientChangePosition) {
    (*route).clients[clientChangePosition]=clientChange
    (*route).Residue-=clientChange.Demand
    success=true
  }
  return
}

func (route *Route)RepositionClient(clientToReposition Client)  (success bool){
  //TAG:="(route *Route)RepositionClient(clientToReposition *Client)  (success bool)"
  //log.Println(TAG)

  position:=route.GetPositionClient(&clientToReposition)
  numClients:=len((*route).clients)

  clients:=&(*route).clients
  if numClients>1 {
    newPosition:=-1
    for (newPosition == -1 || newPosition == position) {
      t := time.Now()
      r := rand.New(rand.NewSource(int64(t.Nanosecond())))
      newPosition = r.Intn(numClients-1)
    }

  clientFromReposition:=(*clients)[newPosition]


  success1:=route.ChangeClient(clientToReposition, newPosition)
  success2:=route.ChangeClient(clientFromReposition, position)

  if(success1 == false|| success2 == false){
    route.LogClients()
    log.Println("success1,success2",success1,success2)
    log.Println("clientToReposition,newPosition",clientToReposition.Id,newPosition)
    log.Println("clientFromReposition,position",clientFromReposition.Id,position)

  }
  }
  return
}

func (route *Route)GetPositionClient(clientToCompare *Client)  (position int){
  //TAG:="(route *Route)GetPositionClient(client *Client)  (position int)"
  //logPrintln(TAG)
  position=-1
  for i,client := range (*route).clients {
    if client.Id==(*clientToCompare).Id {
      position=i
    }
  }
  return
}

func (route *Route)Retime()  {
  //TAG:="(route *Route)Retime()"
  //logPrintln(TAG)

  var time float64
  for _,path := range (*route).R {
    time+=path.Time
  }
  (*route).Time=time
}
func (route *Route)GetTime() (time float64) {
  for _,path := range (*route).R {
    time+=path.Time
  }
  return
}
func (route *Route)GetReliability() (reliability float64) {
  reliability=(*route).deposit.Reliability
  return
}
func (route *Route)LoadPaths()  {
  nClients:=len((*route).clients)

  if(nClients>0){
    (*route).R=make([]Path,0)
    previousStation:=(*route).deposit.Station
    var actualStation Station
    for _,client := range (*route).clients{
      actualStation=client.Station
      route.AddPath(&previousStation, &actualStation)
      previousStation=actualStation
    }
    route.AddPath(&previousStation, &(*route).deposit.Station)
    route.Retime()
  }
}
func (route *Route)ValidateClientsRepeat(lisClients *map[int]bool,tag string)  (err error){
  clients:=(*route).clients
  for _,client := range clients {
    _,ok:=(*lisClients)[client.Id]
    if ok==false {
      (*lisClients)[client.Id]=true
    }else{
      err=errors.New("Client Repeat in the solution: "+tag)
      log.Println("-------------------------------------------------------------------------------")
      log.Println("Client repeat:",client.Id)
      route.LogClients()
      break
    }
  }
  return
}
func (route *Route)LogClients()  {
  lClients:=""
  clients:=(*route).clients
  log.Println("# route: ",route.Id," # deposit: ",route.deposit.Id," # clients: ",len(clients))
  for _,client := range clients {
    lClients+=strconv.Itoa(client.Id)+", "
  }
  log.Println("clients: ",lClients)
}
