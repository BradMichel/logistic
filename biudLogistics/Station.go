package biudLogistics

import(
//  "log"
  "time"
  "math/rand"
)

func (station *Station)OrderDistanceAsc()  {
  ////TAG:="(station *Station)OrderDistanceAsc()"
  ////logPrintln(TAG)

  l:=len((*station).Distances)
  distances:=make([]float64,l,l*2)
  copy(distances[:],(*station).Distances)
  distancesOrderAsc:=&(*station).DistancesOrderAsc
  for index := 0;  index< len(distances); index++ {
    var(
      men int
      valMen float64
    )
    men,valMen=0,-1
    for id,distance := range distances {
      if ((valMen<0 || distance<valMen) && distance > 0) {
        men,valMen=id,distance
      }
    }
    *distancesOrderAsc=append(*distancesOrderAsc,men)
    distances[men]=-1
  }
}
func (station *Station)NextClient(clients *Clients,start int) (client *Client)  {
  //TAG:="(station *Station)NextClient(clients *Clients) (client *Client)"
  //logPrintln(TAG)

  distancesOrderAsc:=&(*station).DistancesOrderAsc
  for i,id := range (*distancesOrderAsc) {
    if i>=start {
      nextClient,ok:=(*clients)[id]
      if ok==true && nextClient.Visited==false{
        client=nextClient
        break
      }
    }
  }
  return
}

func (station *Station)NextClientForTimeWindow(clients *Clients) (minClient *Client)  {
  //TAG:="(station *Station)NextClientForTimeWindow(clients *Clients) (client *Client)"
  //logPrintln(TAG)
  minClient=nil
  for _,nextClient:=range (*clients){
    if (nextClient.Visited==false && nextClient.TimeLimit.IsZero()==false) && (minClient==nil || minClient.TimeLimit.Time.After(nextClient.TimeLimit.Time)) {
      minClient=nextClient
    }
  }
  return
}
func (station *Station)NextRandomClient(clients *Clients)(client *Client){
  nClients:=len(*clients)
  if(nClients>0){
    repetitions:=0
    for client == nil || (client.Visited==false  && repetitions<1000){
      t:=time.Now()
      r:=rand.New(rand.NewSource(int64(t.Nanosecond())))
      position:=r.Intn(nClients)
      client=(*clients)[position]
      repetitions++
    }
  }
  return
}
func (station *Station)GetDistaceTo(to *Station)  (distance float64){
  //TAG:="(station *Station)GetDistaceTo(to *Station)  (distace float64)"
  //logPrintln(TAG)

  distance=(*station).Distances[(*to).Id]
  return
}
func (station *Station)InitializeVisit()  {
  ////TAG:="(station *Station)InitializeVisit()"
  ////logPrintln(TAG)
  station.Visited=false
}
func (station *Station)ValidatePreviousVisit(clientForTimeWindow *Client,client *Client)  bool{
  //TAG:="(station *Station)ValidatePreviousVisit(clientForTimeWindow *Client,client *Client)  bool"
  //logPrintln(TAG)

  tt:=(*station).TimeVisited.Time
  t:=MyTime{tt}
  t.Time=t.Add(time.Duration((*station).Distances[(*client).Id]) * time.Second)
  t.Time=t.Add(time.Duration((*client).Distances[(*clientForTimeWindow).Id]) * time.Second)
  return clientForTimeWindow.TimeLimit.Time.Before(t.Time)
}
