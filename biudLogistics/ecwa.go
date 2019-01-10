package biudLogistics
/*
import(
	"math/rand"
	"time"
	"net/http"
	"github.com/gorilla/context"
	"math"
	//"log"
)



func PostEcwa(w http.ResponseWriter, r *http.Request) {
	//TAG:="PostEcwa"
	//logPrintln(TAG)
	// Obtiene los datos enviados del Usuario
	content :=context.Get(r,"body").(*Content)
	// Obtieneene el albergue mas cercano a cada cliente y que satisfaga su demanda
	routes,_ := nearestRefuge(content.Distances,content.Demand)

	//---------------------------May 23 2016-----------------------


	//	Se crean diferentes variantes de las rutas
	multipleRoutes:=routes.multipleSolutions(content.Distances,content.Demand);


	// Se hace el ruteo de la lista
	multipleRoutes.getFirstRoutes(content.Distances)

	// Se comparan las rutas para evaluar cada resultado
	fit:=multipleRoutes.fit()

	order:=orderMaxMin(fit)

	multipleRoutes.runGenetics(fit,order,content.Distances,content.Demand);

	//--------------------------------------------------------------

	fit=multipleRoutes.fit()
	order=orderMaxMin(fit)

	// Se formatea la rta en formato json para enviar la respuesta al usuario
	responseJson(w,201,map[string]interface{}{"List Routes":multipleRoutes,"fit":fit,"order":order})

}

// Funcion Albergue mas cercano que satesfaga demanda de cliente
func nearestRefuge(distances [][]float64, demand []int) (routes RouteList, initial map[int]int){
	//TAG:="nearestRefuge"
	//logPrintln(TAG)

	// Se crea la variable route que es un vector de trayecto
	var route []Path
	var path Path
	// Alvergue mas cercano al cliente
	var disMin = map[int]int{}
	// #Ruta por albergue
	initial=make(map[int]int)

	// Se recorren las demandas para encontrar los Albergues
	for i := range demand{
		if demand[i] <= 0{
			// Se crea una ruta por cada Albergue
			route = make([]Path,0,0)
			// Se inicializa la ruta con el Albergue
			path = Path{0,i,0}
			route = append(route,path)
			R:= Route{R:route,Residue:demand[i]}
			// Se guarda la ruta
			routes = append(routes,R)
			//Se guarda la relacion #Albergue -> Ruta
			initial[i]=len(routes)-1

		}else{
			// Para los Clientes
			// Se obtiene el albergue mas cercano
			disMin[i]=getDisMin(distances[i],demand,make([]bool,len(distances),len(distances)))
		}
	}

	// Se recorre los clientes para agregarlos a las ruta de su Albergue mas cercano
	for i,j := range disMin{
			//
			stop := false
			veto:=make([]bool,len(distances),len(distances))
			for stop == false{
			find:=getBase(&routes,demand[i],distances[i][j],initial,i,j)
			if find ==true{
				stop = true
			}else{
				veto[j]=true
				disMin[i] = getDisMin(distances[i],demand,veto)
				j=disMin[i]
				if disMin[i] == -1 {
					stop = true
				}
			}
		}
	}
		return
}

func getBase(routes *RouteList,demand int, distances float64,initial map[int]int,i int, j int)  (find bool){
	path := Path{i,j,distances}
	if((*routes)[initial[j]].Residue <= -demand){
		(*routes)[initial[j]].R=append((*routes)[initial[j]].R,path)
		(*routes)[initial[j]].Residue += demand
		(*routes)[initial[j]].Time+= distances
		find=true
	}else{
		find=false
	}
	return
}

func getDisMin(distances []float64, demand []int,veto []bool) (pos int)  {
	var min float64
	min= -1
	pos= -1

	for i :=range distances{
		distance:=distances[i]

			if( (min == -1 && demand[i]<0 && veto[i]==false) || (distance > 0 && distance < min && demand[i]<0 && veto[i]==false)){
				min = distance
				pos=i
			}
	}

	return
}

/**
	@multipleSolutions
	Funcion para encontrar diferentes variantes de las rutas
	Parameters
	(Route)([][]float64,[]int)
	return
	(SliceRouteList)
**/
/*
func (routes RouteList) multipleSolutions(distances [][]float64,demand []int)  (multipleRoutes SliceRouteList){
	//TAG:="multipleSolutions"
	//logPrintln(TAG)
	//agrega las rutas existentes a la lista de rutas
	multipleRoutes=append(multipleRoutes,routes)

	//se recorre todas las sub rutas
	for i :=range routes{
		route:=routes[i]

		//	cambia el cliente mas lejano a la base por uno aleatorio
		if len(route.R)>1 {
			newRoute:=route.changeMaxHeavy(routes,i,distances,demand)
			// Se agregan las nuevas rutas a la lista
			multipleRoutes=append(multipleRoutes,newRoute)
		}
	}
	return
}

/**
	@changeMaxHeavy
	cambia el cliente mas lejano a la base por uno aleatorio
	Parameters
	(Proute)(RouteList,int.[][]float64,[]int)
	Return
	RouteList
**/
/*
func (route Route) changeMaxHeavy(routes RouteList, pos int,distances [][]float64, demand []int) (newRoutes RouteList) {
	//TAG:="changeMaxHeavy"
	//logPrintln("start",TAG)
	var(
		max, l int
		maxTime float64
		stop bool
	)
	if(len(route.R)>0){
		//Crea la nueva ruta
		newRoutes=make([]Route,len(routes),len(routes))

		for pos := range routes {
			l:=len(routes[pos].R)
			newRoutes[pos].R=make([]Path,l,l)
			for posR := range routes[pos].R {
				newRoutes[pos].R[posR] = Path{I:routes[pos].R[posR].I,J:routes[pos].R[posR].J,Time:routes[pos].R[posR].Time}

			}
			newRoutes[pos].Time=routes[pos].Time
			newRoutes[pos].Residue=routes[pos].Residue
		}

		//logPrintln("newRoutes",newRoutes);
		// Busca el cliente mas lejano
		max=-1
		l=len(routes)
		for i :=range route.R{
			path:=route.R[i]
			//logPrintln("if maxTime",path.Time,maxTime,path.Time>maxTime);
			if(path.Time>maxTime){
				maxTime=path.Time
				max=i
			}
		}

		// Crea un numero aleatorio que sera el # sub ruta con que se va a hacer el cambio
		t := time.Now()
		r := rand.New(rand.NewSource(int64(t.Nanosecond())))
		// Ruta con que se va a intercabiar
		change := max
		for  max ==change  {
			change = r.Intn(l)
		}

		l2:=len(routes[change].R)

		con:=0
		for stop == false{
			//	Posicion con con la cual se va a intercabiar
			p:=r.Intn(l2)
			route1:=routes[pos]
			i1:=route1.R[max].I
			demand1:=demand[i1]
			residue1:=route1.Residue
			route2:=routes[change]
			i2:=route2.R[p].I
			demand2:=demand[i2]
			residue2:=route2.Residue

			// Se verefica si las rutas soportan el cambio
			if(-residue2+demand2>=demand1 && -residue1+demand1>=demand2){
				// Se realiza los cambios entre lsas dos Rutas
				//logPrintln("After Change",route1.R[max],route2.R[p])

				j1:=route1.R[max].J
				j2:=route2.R[p].J
				var timeT float64
				/*for t := range newRoutes[change].R {
					timeT+=newRoutes[change].R[t].Time
				}*/
/*
				// Los dos nuevos pasos de las sub rutas
				var path1=Path{I:i1,J:j2,Time:distances[i1][j2]}
				var path2=Path{I:i2,J:j1,Time:distances[i2][j1]}

				newRoutes[pos].R[max]=path2
				newRoutes[pos].Time=newRoutes[pos].Time-distances[i1][j1]+distances[i1][j2]
				newRoutes[pos].Residue=newRoutes[pos].Residue+demand1-demand2
				newRoutes[change].R[p]=path1
				newRoutes[change].Time=newRoutes[change].Time-distances[i2][j2]+distances[i2][j1]
				//timeT=0
				/* for t := range newRoutes[change].R {
					timeT+=newRoutes[change].R[t].Time
				} */

				/*
				newRoutes[change].Residue=newRoutes[change].Residue+demand2-demand1
				stop=true
			}
			if con ==1000 {
				stop=true
			}
			con++

		}
	}else{
		newRoutes=nil
	}
	return

}

/**
	@getFirstRoutes
	Ruteo de lista de rutas
	Parameters
	(*SliceRouteList)([][]float64)
**/
/*
func (multipleRoutes *SliceRouteList) getFirstRoutes(distances [][]float64)  {
	//TAG:= "getFirstRoutes"
	//logPrintln(TAG)
	var  jInitial, pi, pj int
	var Icomp []int
	// Se recorre la lista de rutas
	for pos1 := range (*multipleRoutes) {

		routes:=&(*multipleRoutes)[pos1]
		// Se recorren las sub rutas
		for pos2 := range (*routes) {
			// Se inicializan las variables de la sub-rutas
			route:=&(*routes)[pos2]
			R:=&(*route).R
			l:=len(*R)
			// vector de los nodos de la sub-ruta
			Icomp=make([]int,l,l)
			for i := range (*R) {
				if(i==0){
					Icomp[i]=(*R)[i].J
				}else{
					Icomp[i]=(*R)[i].I
				}
			}

			// vector nodos visitados
			veto:= make([]bool,len(distances[0]),len(distances[0]))
			conChange:=0
			veto[0]=true

			//	Se recorre la sub-ruta enrutando cada nodo con el mas cercano
			//	Se condiciona a que el nodo ya no haya sido agregado a la ruta
			for i := 1; i < l; i++ {
				jInitial=1
				path:=&(*R)[i]
				pathI:=&(*R)[i-1]
				pi=(*pathI).J
				min:=distances[pi][pj]
				posMin:=-1
				for j := jInitial; j < l; j++ {
					pj=Icomp[j]
					d:=distances[pi][pj]
					if( (veto[pj]==false) && (posMin==-1 || min>d)){
							min=d
							posMin=pj
							conChange++
						}

				}
				if(posMin==-1){
					posMin=path.J
					//logPrintln("posMin no found")
				}

				route.Time = route.Time-path.Time+min
				path.I=(*pathI).J
				path.J=posMin
				veto[posMin]=true
				path.Time=min
			}
		}
	}

}

/**
	@fit
	evaluacion de cada resultado
	Parameters
	(SliceRouteList)
	Return
	[]float64
**/
/*
func (mulRoutes SliceRouteList) fit() (fit []float64)  {
	//TAG:="fit"
	//logPrintln(TAG)

	l:=len(mulRoutes)	//Número de soluciones
	var fitTime,fitResidue []float64										//	Vectores de evaluacion
	fitTime=make([]float64,l,l)													// Vector resultado evaluacion del tiempo de cada respuesta
	fitResidue=make([]float64,l,l)											//	Vector resultado evaluacion del residuo de cada respuesta
	fit=make([]float64,l,l)															// Vector resultado founcion evaluatoria (CT*T+CR*R) en cada respuesta


	for i := range mulRoutes[0] {												//	Se recorre cada sub-ruta
			for n := 0; n < l; n++ {												//	Se guarda encuentran los totales de tiempos y residuos de cada solucion recorriendo las sub-rutas en cada solucion
				subRoute:=mulRoutes[n][i]											//	Se selecciona la sub-ruta actual en el recorrido
				fitTime[n]+=subRoute.Time											//	Se va acumulando el tiempo de cada sub-ruta por cada respuesta
				fitResidue[n]+=float64(-subRoute.Residue)			//	Se va acumulando el residuo de cada sub-ruta por cada respuesta
			}
	}
	for n := range fitTime {
		fit[n]= (CfitTime*fitTime[n])+(CfitResidue*fitResidue[n])		//	Evaluacion de cada solucion con la funcion evaluatoria (CT*T+CR*R)
	}
	return
}

func  orderMaxMin(fit []float64)(order []int)  {
	//TAG:="orderMaxMin"
	//logPrintln(TAG)
	l:=len(fit)
	order=make([]int,l,l)
	for i := 0; i < l; i++ {
		for j:=0;j<l;j++ {
			if fit[i]<fit[j] {
				order[i]++
			}
		}
	}
	return
}

func (multipleSolutions *SliceRouteList)runGenetics(fit []float64,order []int,distances [][]float64, demand []int)  {
	//TAG:="runGenetics"
	var(change1,change2,maxRepeat int)
	//logPrintln(TAG)
	l:=len(*multipleSolutions)
	div1:=int(Round(float64(l)*CselectS,1,0))
	div2:=l-div1

	sizePairs:=int(Round(float64(div1)/2,1,0))
	numSubRoutes:=len((*multipleSolutions)[0])
	for pair := 0; pair < sizePairs ; pair++ {

		t:=time.Now()
		r:=rand.New(rand.NewSource(int64(t.Nanosecond())))

		for change2,change1:=r.Intn(numSubRoutes),r.Intn(numSubRoutes);  (change2 == change1 || order[change1]>=div1 || order[change2]>=div1) && maxRepeat <=1000 ;{
			change2 = r.Intn(numSubRoutes)
			maxRepeat++
		}
		if maxRepeat!=1000 {
			route1:=(*multipleSolutions)[pair].duplicate()
			route2:=(*multipleSolutions)[pair+sizePairs].duplicate()
			subRoute1:=route1[change1]
			subRoute2:=route2[change2]
			route1.deletePaths(subRoute2,change1,distances,demand)
			route2.deletePaths(subRoute1,change2,distances,demand)

			route1.insertPaths(subRoute2,distances,demand)
			route2.insertPaths(subRoute1,distances,demand)

			*multipleSolutions=append(*multipleSolutions,route1)
			*multipleSolutions=append(*multipleSolutions,route2)
	}

	}

	//logPrintln("div1",div1,"div2",div2)
}

func (route *RouteList) deletePaths(subRoute2 Route,change int,distances [][]float64, demand []int)  {
	l:=len(subRoute2.R)
	for k := 1; k < l; k++ {
		for i := range (*route) {
			if i != change {
				subRoute:=&(*route)[i]
				R:=&subRoute.R
				Time:=&subRoute.Time
				Residue:=&subRoute.Residue
				lR:=len(*R)
				//logPrintln("R",R)
				for iS := 1; iS < lR; iS++ {
					//logPrintln("equals",(*R)[iS].J,subRoute2.R[k].J,(*R)[iS].J==subRoute2.R[k].J)
					if (*R)[iS].J==subRoute2.R[k].J {
						//logPrintln("Change ",(*R)[iS])
						*Time+= -(*R)[iS].Time
						*Residue+= -demand[(*R)[iS].J]
						var newR []Path
						if(iS==lR-1){
							//logPrintln("end found")
							newR=(*R)[:iS]
							//logPrintln("continue",newR)
						}else{
							newR=append((*R)[:iS],(*R)[iS+1:]...)
						}
						*R=newR
						if lR-1!=iS {
							(*R)[iS].I=(*R)[iS-1].J
							tChange:=(*R)[iS].Time
							(*R)[iS].Time=distances[(*R)[iS].I][(*R)[iS].J]
							*Time+=(*R)[iS].Time-tChange
						}
						iS=lR
					}
				}
				//logPrintln("stop")
			}
		}
	}
}
func (route *RouteList) insertPaths(subRoute2 Route, distances [][]float64, demand []int)  {
	l:=len(subRoute2.R)
	lR:=len(*route)
	for k := 1; k < l; k++ {
		t:=time.Now()
		r:=rand.New(rand.NewSource(int64(t.Nanosecond())))
		J:=subRoute2.R[k].J
		stop:=false
		for change := r.Intn(lR); stop==true; {
			Residue:=&(*route)[change].Residue
			if(*Residue<=demand[J]){
				subRoute:=&(*route)[change]
				Time:=&subRoute.Time
				lSR:=len((*route)[change].R)
				changePath:=r.Intn(lSR-1)+1
				I:=subRoute.R[changePath-1].J
				I2:=J
				J2:=subRoute.R[changePath].J
				tChange1:=distances[I][J]
				tChange2:=distances[I2][J2]
				*Time+= tChange1+tChange2-subRoute.R[changePath].Time
				*Residue+=demand[J]
				newPath:=make([]Path,1,1)
				newPath[0]=Path{I:I,J:J,Time:distances[I][J]}
				newR:=append(subRoute.R[:changePath],newPath...)
				newR=append(newR,subRoute.R[changePath:]...)
				subRoute.R=newR
				stop=true
			}
		}
	}
}

func (route RouteList)duplicate()(newRoute RouteList)  {
	//Crea la nueva ruta
	newRoute=make([]Route,len(route),len(route))

	for pos := range route {
		l:=len(route[pos].R)
		newRoute[pos].R=make([]Path,l,l)
		for posR := range route[pos].R {
			newRoute[pos].R[posR] = Path{I:route[pos].R[posR].I,J:route[pos].R[posR].J,Time:route[pos].R[posR].Time}
		}
		newRoute[pos].Time=route[pos].Time
		newRoute[pos].Residue=route[pos].Residue
	}
	return
}


func Round(val float64, roundOn float64,places int) (newVal float64)  {
	var round float64
	pow:=math.Pow(10,float64(places))
	digit:=pow*val
	_, div:=math.Modf(digit)
	if div >= roundOn{
		round=math.Ceil(digit)
	}else{
		round=math.Floor(digit)
	}
	newVal=round/pow
	return
}
*/
