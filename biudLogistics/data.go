package biudLogistics

import(
	"time"
)

const(
	N=50
	R=50
	P=50
	CfitTime=7
	CfitResidue=14
	CselectS=0.7
	CselectI=1-CselectS
)

type(
	MyTimes []MyTime
	MyTime struct{
		time.Time
	}
	Stop interface{
		OrderDistanceAsc()
		NextClient() Client
	}
	Station struct{
		Id int  `json:"id,omitempty"`
		Distances []float64 `json:"distances,omitempty"`
		DistancesOrderAsc	[]int `json:"distances,omitempty"`
		Visited bool	`json:"visited,omitempty"`
		TimeVisited MyTime
	}
  Deposits map[int]*Deposit
  Deposit struct{
		Station
    Penalized bool `json:"penalized,omitempty"`
    Capacity int `json:"capacity,omitempty"`
		Load	int	`json:"load,omitempty"`
    Clients Clients `json:"clients,omitempty"`
		Reliability float64 `json:"reliability,omitempty"`
  }
  Clients map[int]*Client
  Client struct{
		Station
    Demand int  `json:"demand,omitempty"`
		TimeLimit MyTime
  }
	Vehicles map[int]*Vehicle
	Vehicle struct{
		Id int	`json:"id,omitempty"`
		Capacity int	`json:"capacity,omitempty"`
		Load int `json:"load,omitempty"`
	}

// Tipo de dato para la respuesta
Rta struct{
	// Tiempo resultante de las rutas
	Time float64 `json:"time,omitempty"`
	// Rutas resultantes
	Routes RouteList `json:"routes,omitempty"`
}

// Tipo de dato de los datos recibidos
Content struct{
	// Matriz de distancias
	Distances 		[][]float64 `json:"distances,omitempty"`
	// Vector de demanda
	Demand []int `json:"demand,omitempty"`
	// Vector de capacidad
	Capacity []int `json:"capacity,omitempty"`
	// Vector de confiabilidad
	Reliability []float64 `json:"reliabilitys,omitempty"`
	//Ventanas de tiempo
	TimeWindows MyTimes `json:"timewindows,omitempty"`
	//Tiempo de inicio de recorrido
	TimeStart MyTime `json:"timeStart,omitempty"`
	}

// Tipo de dato de trayecto
Path struct{
	// Punto de partida
	I int `json:"i,omitempty"`
	// Punto de llegada
	J int `json:"j,omitempty"`
	// Tiempo de trayecto
	Time float64 `json:"time,omitempty"`
}

ByTimeSliceRouteList  struct{
	SliceRouteList
}
ByReliabilitySliceRouteList struct{
	SliceRouteList
}
ByDistanceCrowding struct{
	SliceRouteList
}
// Soluciones
SliceRouteList []RouteList
// rutas
RouteList struct{
	Routes []Route `json:"routes,omitempty"`
	TimesVisited MyTimes
	Pos int `json:"pos,omitempty"`
	Comparison struct{
		Time float64 `json:"time,omitempty"`
		Reliability float64 `json:"reliability,omitempty"`	// Confiabilidad
	}
	StackingDistance struct{
		Time float64 `json:"time,omitempty"`
		Reliablity float64 `json:"reliability,omitempty"`
	} `json:"stackingDistance,omitempty"`
}



// ruta
Route struct{
	Id int `json:"id,omitempty"`
	// Vector de trayecto
	R []Path `json:"r,omitempty"`
	//	Deposito
	deposit Deposit
	// Clientes
	clients []Client
	// Tiempo total de la ruta
	Time float64 `json:"time,omitempty"`
	// Residuo de demanda de la ruta
	Residue int `json:"residue,omitempty"`
}
Individuals []Individual
Individual struct{
	n int
	p []int
}
Dominance struct{
	More int
	Equal int
	Less int
}
)
