package biudLogistics

import (
	"log"
	"reflect"
)

func (routes *RouteList) GetNeighbors(kind string,vehicles *Vehicles, deposits *Deposits, clients *Clients) {
	TAG:="(routes *RouteList)GetNeighbors(vehicles *Vehicles,deposits *Deposits,clients *Clients)"
	//log.Println(TAG)

	for _, vehicle := range *vehicles {
		vehicle.Load = 0
	}
	for _, deposit := range *deposits {
		deposit.InitializeVisit()
		deposit.Load = deposit.Capacity
	}
	for _, client := range *clients {
		client.InitializeVisit()
	}

	idDeposit := -1
	for _, vehicle := range *vehicles {
		route := new(Route)
		deposit := deposits.GetNext(idDeposit)
		if deposit == nil {
			//log.Println("deposit",deposit)
			break
		} else {
			route.GetNeighbor(kind,vehicle, deposit, clients)
			if(len((*route).clients)>0){
				(*routes).Routes = append((*routes).Routes, *route)
				(*routes).LoadComparison()
				if err:=routes.ValidateClientsRepeat(TAG);err!=nil{
					log.Println(err)
				}
			}
		}
	}
}

func (routes *RouteList) GetTimeWindows(vehicles *Vehicles, deposits *Deposits, clients *Clients, timeStart *MyTime) {
	//TAG:="RouteList.GetTimeWindows"
	//logPrintln(TAG)
	idDeposit := -1
	for _, vehicle := range *vehicles {
		route := new(Route)
		//newClients:=new(Clients)
		//(*route).clients=*newClients
		deposit := deposits.GetNext(idDeposit)
		deposit.TimeVisited = MyTime{(*timeStart).Time}
		if deposit == nil {
			break
		} else {
			route.GetTimeWindows(vehicle, deposit, clients)
			(*routes).Routes = append((*routes).Routes, *route)

		}
	}
}
func (routes *RouteList) LoadComparison() {
	var time, reliability float64
	for _, route := range (*routes).Routes {
		time += route.GetTime()
		reliability += route.GetReliability()
		//log.Println("reliability",reliability)
	}
	(*routes).Comparison.Time, (*routes).Comparison.Reliability = time, reliability
}
func (routes *RouteList) Swap() {
	//TAG:="(routes *RouteList)Swap()"
	//log.Println(TAG)

	restricted := make(map[int]bool)
	success := false
	maxLoop := 0
	for success == false && maxLoop <= 10 && len(routes.Routes) > 1 {
		firstLongerRoute := routes.FindRouteMoreTime(&restricted)
		secondLongerRoute := routes.FindRouteMoreTime(&restricted)
		success = firstLongerRoute.ExchangeLastCustomers(secondLongerRoute)
		maxLoop++
	}
}
func (routes *RouteList) InsertionLastClientLongerRoute() {
	//TAG:="(routes *RouteList)InsertionLastClientLongerRoute()"
	//log.Println(TAG)

	if len((*routes).Routes) > 0 {
		restricted := make(map[int]bool)
		longerRoute := routes.FindRouteMoreTime(&restricted)
		//log.Println("longerRoute")
		//log.Printf("%+v\n",longerRoute)
		numClients := len((*longerRoute).clients)
		lastClient := (*longerRoute).clients[numClients-1]
		longerRoute.RepositionClient(lastClient)

	}
}
func (routes *RouteList) FindRouteMoreTime(restricted *map[int]bool) (routeMax *Route) {
	//TAG:="(routes *RouteList)FindRouteMoreTime(restricted *map[int]bool)(routeMax *Route)"
	//log.Println(TAG)

	routeMax = routes.GetInitialRoute(restricted)
	for i, route := range (*routes).Routes {
		_, ok := (*restricted)[route.Id]

		if route.Time > routeMax.Time && ok == false {
			routeMax = &(*routes).Routes[i]
		}
	}
	(*restricted)[routeMax.Id] = true
	return
}
func (routes *RouteList) GetInitialRoute(restricted *map[int]bool) (route *Route) {
	//TAG:="(routes *RouteList)GetInitialRoute(restricted *map[int]bool) (route *Route)"
	//log.Println(TAG)
	route = nil
	condition := len((*routes).Routes) - 1
	stop := 0
	//logPrintln("condition,restricted",condition,restricted)

	for stop <= condition {
		_, ok := (*restricted)[stop]
		//logPrintln("ok",ok)
		if ok == false {
			//logPrintln("stop,(*routes).Routes",stop,(*routes).Routes)
			route = &(*routes).Routes[stop]
			stop = condition + 1
		}
		stop++
	}
	//log.Println(route)
	return
}
func (routes *RouteList) GetTime() (time float64) {
	//TAG:="(routes *RouteList)GetTime()  (time float64)"
	//logPrintln(TAG)

	for _, route := range (*routes).Routes {
		time += route.Time
	}
	return
}
func (routes *RouteList) FirstFront(routesComparate *RouteList, j int, individual *Individual) {
	//TAG:="(routes *RouteList)FirstFront(routesComparate *RouteList, j int, individual *Individual)"
	//logPrintln(TAG)
	//logPrintln(len((*routes).Routes))
	if len((*routes).Routes) > 0 {
		dominance := new(Dominance)
		comparisonA := reflect.ValueOf(&(*routes).Comparison).Elem()
		comparisonB := reflect.ValueOf(&(*routesComparate).Comparison).Elem()
		numComparison := comparisonA.NumField()
		for i := 0; i < numComparison; i++ {
			if comparisonA.Field(i).Interface().(float64) < comparisonB.Field(i).Interface().(float64) {
				(*dominance).Less++
			} else if comparisonA.Field(i).Interface().(float64) == comparisonB.Field(i).Interface().(float64) {

				(*dominance).Equal++
			} else {
				(*dominance).More++
			}
		}

		if dominance.Less == 0 && dominance.Equal != numComparison {
			(*individual).n++
		} else if dominance.More == 0 && dominance.Equal != numComparison {
			(*individual).p = append((*individual).p, j)
		}

	}else{
		log.Println("0 routes")
		log.Println((*routes).Routes)
	}
}

func (routes *RouteList) GetLenObjetives() (len int) {
	//TAG:="(routes *RouteList)GetLenObjetives() (len int)"
	//logPrintln(TAG)

	len = reflect.ValueOf(&(*routes).Comparison).Elem().NumField()
	return
}

func (routes *RouteList) GetValueIndexField(t string, index int) (value float64) {
	//TAG:="(routes *RouteList)GetValueIndexField(t string,index int)  (value float64)"
	//log.Println(TAG)
	obj := routes.GetReflect(t)
	value = obj.Field(index).Interface().(float64)
	return
}

func (routes *RouteList) SetValueIndexField(t string, index int, value float64) {
	//TAG := "(routes *RouteList)SetValueIndexField(t string, index int,value float64)"
	//log.Println(TAG)

	obj := routes.GetReflect(t)
	obj.Field(index).SetFloat(value)
}

func (routes *RouteList) GetReflect(t string) (obj reflect.Value) {
	//TAG := "(routes *RouteList)GetReflect(t string) (obj reflect.Value)"
	//log.Println(TAG)

	if t == "Comparison" {
		obj = reflect.ValueOf(&(*routes).Comparison).Elem()
	}
	if t == "StackingDistance" {
				obj = reflect.ValueOf(&(*routes).StackingDistance).Elem()

	}
	return
}
func (routes *RouteList) GetTimesVisited(deposits *Deposits, clients *Clients) {
	//TAG:="(routes *RouteList)GetTimesVisited(deposits *Deposits,clients *Clients)"
	//logPrintln(TAG)

	lenStation := len(*deposits) + len(*clients)
	var station Station
	for i := 0; i < lenStation-1; i++ {
		_, ok := (*clients)[i]
		if ok {
			station = (*clients)[i].Station
		} else {
			_, ok = (*deposits)[i]
			if ok {
				station = (*deposits)[i].Station
			}
		}
		(*routes).TimesVisited = append((*routes).TimesVisited, station.TimeVisited)
	}
}
func (routes *RouteList)GetReliabilitys() (reliability float64) {
	for _,route := range (*routes).Routes {
		reliability+=route.GetReliability()
	}
  return
}
func (routes *RouteList)SetReliability()  {
	reliability:=routes.GetReliabilitys()
	(*routes).Comparison.Reliability=reliability
}
func (routes *RouteList) LoadPaths() {
	for pos := range (*routes).Routes {
		route := &(*routes).Routes[pos]
		route.LoadPaths()
	}
	(*routes).Comparison.Time = (*routes).GetTime()
}
func (routes *RouteList)ValidateClientsRepeat(tag string)  (err error){

	var lisClients map[int]bool
	lisClients=make(map[int]bool)

	for _,route := range (*routes).Routes {
		if err = route.ValidateClientsRepeat(&lisClients,tag); err!=nil{
			break
		}
	}
	return
}
