package biudLogistics
/*
import(
	"math/rand"
	"time"
	"runtime"
	"net/http"
	"github.com/gorilla/context"
)


func PostBee(w http.ResponseWriter, r *http.Request) {

	// Obtiene los datos enviados del Usuario
	content :=context.Get(r,"body").(*Content)

	// Llama la funcion vecino mas cercano
	// Recibe las rutas y el deposito
	routes,base := nearestNeighbor(content.Distances,content.Demand,content.Capacity)

	// Recorre todas las rutas
	for i := range routes {
		// Obtiene los tiempos de cada ruta
		routes[i].getWeigts(content.Distances,base)
	}

	// min =  Tiempo de la ruta con menor tiempo
	// acu = Tiempo total de las rutas
	//var acu, acu2 float64


	// Recorre las rutas
	for i := range routes {
		// Se imprime en consola el numero de la ruta
		////logPrintln("Ruta ", i)
		// Se imprime en consola la ruta del vecino mas cercano
		////logPrintln("Vecino mas cercano: ",routes[i])
		// Se acumula los tiempos de las rutas
		//acu += routes[i].Time
		// Se llama la permutacion interna para cada ruta
		routes[i].internalPermutations(content.Distances,base)
		// Se imprime en consola la ruta resultante de la permutaciones internas
		////logPrintln("Permutaciones internas: ",routes[i])
		////logPrintln()
		//acu2 += routes[i].Time
		// Se compara el tiempo de la ruta con el tiempo minimo actual

	}

	// Se imprime en consola el tiempo total de las rutas para el vecino mas cercano
	////logPrintln("Timenpo total Vecino mas cercano: ", acu)
	// Se imprime el nuevo tiempo total de las rutas luego de las permutaciones internas
	////logPrintln("Timenpo total permutaciones internas: ", acu2)

	// Si hay mas de dos rutas se hace permutacion entre rutas
	if len(routes)>2 {

		////logPrintln("Premutaciones entre rutas")
		// Se llama la funcion para intercambio entre rutas
		// posMin -> el numero de la ruta que se hereda
		routes.permutationRoutesBetween(content.Distances, content.Demand,base)

	}

	// Tiempo final total de las rutas
	var Time float64

	// Se recorre las rutas para acumular el nuevo tiempo total
	for i := range routes {
		Time += routes[i].Time
	}

	// Se crea una variable rta de tipo Rta para enviar la respuesta
	rta := Rta{Routes:routes,Time:Time}

	// Se formatea la rta en formato json para enviar la respuesta al usuario
	responseJson(w,201,rta)

}

// Funcion vecino mas cercano
// Recibe distancias, demanda, capacidad
// Devuelve las rutas como Routelist y la base como int
func nearestNeighbor(distances [][]float64, demand []int ,capacity []int) (routes RouteList, initial int){

	// Se crea la variable route que es un vector de trayecto
	var route []Path
	// Se crea la variable visited que es un vector booleano para registrar los clientes visitados
	var visited = make([]bool,len(distances),len(distances))

	var acu, con, pos, posVisited, end, detener int

	// Se recorren las demandas para encontrar el deposito
	for i := range demand{
		if demand[i] == 0{
			pos = i;
		}
	}
	// Se selecciona el deposito como punto de partida de la ruta
	initial = pos
	// Se registra el deposito como visitado
	visited[pos]=true;
	posVisited++

	// Se comienza a crear las rutas
	for {

		var min float64

		var posMin int = -1

		// Se busca el proximo cliente mas cercano
		for j := range distances[pos]{
			// Se verifica que el proximo cliente no exceda la capacidad del vehiculo
			// y que tenga una distancia menor al minimo actual y que no haya sido visitado

			if ( pos != j && acu+demand[j] <= capacity[con] && (distances[pos][j] < min || min == 0) && !visited[j]){
				if !visited[j]{
					// Se selecciona el nuevo cliente a visitar

					min = distances[pos][j]
					posMin = j


				}
			}

		}


		// Se verifica si hubo un proximo cliente  o
		// Si la cantidad de rutas actuales supera la cantidad de vehiculos y
		// la capacidad del vehiculo es superada por la demanda actual de la ruta
		if(posMin == -1 || (len(routes)<len(capacity) && acu+demand[posMin] > capacity[con])){

			// Se Guarda la ruta en la lista de rutas
			// Se crea la ruta R
			R:= Route{R:route, Residue: capacity[con]-acu}
			// Se agrega la nueva ruta R a la lista de rutas routes
			routes = append(routes,R)
			// Inicializa la ruta en ceros
			route = make([]Path,0,0)
			acu = 0
			con += 1
			pos = initial
		}

		// Se verifica si hay un proximo cliente
		if (posMin>=0){

			// Se registra el nuevo cliente como visitado
			visited[posMin] = true
			posVisited++
			// Se crea el nuevo trayecto P
			P:=Path{I:pos,J:posMin,Time:min}
			// Agrega el nuevo trayecto a la ruta
			route = append(route,P)
			// Agrega la demanda del cliente visitado a la acumulada de la ruta
			acu += demand[posMin]
			// Se suma una nueva ruta a la cantidad de rutas actuales
			end += 1
			// Se verifica si las rutas actuales son iguales a la cantidad de vehiculos
			if (end == len(demand)-1)  {

				// Se detienee la creacion de rutas
				break
			}

			// Se selecciona la posicion visitada como de partida para el proximo trayecto
			pos = posMin
		}else{

			// Se detiene el enrutador cuando no encuentra en mas de 1000 oportunidades un proximo cliente
			if (detener == 1000){

				break
			}
			detener += 1
		}

		if con==len(capacity) {
			break
		}


	}


	// Se verifica si hay trayectos en la ruta  y si hay vehiculo que puedan tomar esa ruta para agregarla a la lista de rutas
	if (len(route)>0 && len(routes)<len(capacity)){
		R:= Route{R:route, Residue: capacity[con]-acu}
		routes = append(routes,R)
	}

		return
}


// funcion para obtener el tiempo total de una ruta
func (route *Route) getWeigts(distances [][]float64,base int) {
	// Se inicializa el tiempo total de la ruta en 0
	route.Time = 0
	// Se recorren todas las trayectorieas de la ruta
	for i:= range route.R{
		// Se agrega el tiempo de cada trayectoria al tiempo total de la ruta
		route.Time += route.R[i].Time
	}
	// Se agrega el tiempo desde ultimo cliente a la base
	route.Time += distances[route.R[len(route.R)-1].J][base]
}


// Funcion para permutaciones internas de la ruta
func (route *Route ) internalPermutations(distances [][]float64,base int) {

	// Numero de trayectos de la ruta
	l:=len(route.R)
	// Numero de permutaciones internas
	f:=factorial(l)
	// Numero de CPUs del equipo
	numCPU := runtime.NumCPU()
	// Crea un vector de pociciones
	positions:= make([]int,l)
	// inicializa el vector de posciciones
	for i:= range route.R{
		positions[i] = i
	}
	// Genera un Permutador de las posiciones
	p,err:=NewPerm(positions,nil)
	if err!=nil {

	}

	var lim,acu int



	if f>100000 {
		lim=1000000
		// lim=int(a)
	}else{
		lim=f
	}

	for acu < f {

		if acu+lim>f && acu+lim < f+lim {
			lim=f-acu
		}

		// break
		nPer:= p.NextN(lim)

		perm:= nPer.([][]int)
	    // Se crea un canal de comunicacion con la rutina

		c:=make(chan Route,numCPU)
		// se crean rutinas de acuerdo al numero de CPUs del equipo
		for j := 0; j < numCPU; j++ {
			/*Se ejecuta la funcion internalPermutation como rutina dividiendo
			la cantidad de permutaciones segun el numero de CPUs del equipo*/
			/*
			go route.internalPermutation(perm,(j*lim/numCPU),((j+1)*lim/numCPU),distances,base,c)
	   }

	   	//  Se reciben las respuestas de las rutinas por medio del canal
	    for j := 0; j < numCPU; j++ {
	    	// Se recibe la respuesta de cada rutina con la mejor ruta de sus permutaciones por el canal c y se guarda en r
	    	r:=<-c
	    	// Se compara la nueva ruta con la actual
	        if r.Time<route.Time {
	        	// Se guarda la nueva ruta como la mejor
	              	*route=r
	              }
	    }
	    break
		acu+=lim
	}
}

// Funcion para paralelizar las permutacion
func (route Route) internalPermutation(perm [][]int,a int,b int,distances [][]float64, base int, c chan Route) {
	// Recorre el tramo de las permutaciones a hasta b
	for i := a; i < b; i++ {
		// Se crea la variable el vector newPositions con las nuevas posiciones de los clientes en la ruta
		newPositions:=perm[i]
		newRoute:= Route{R:make([]Path, len(route.R),cap(route.R))}
		for j := range newPositions{
			newRoute.R[newPositions[j]] = route.R[j]
		}
		for j:= range newRoute.R{
			if j == 0 {
				newRoute.R[0].I = route.R[0].I
			}
			if(j>0){
				newRoute.R[j].I =newRoute.R[j-1].J
			}
			newRoute.R[j].Time =distances[newRoute.R[j].I][newRoute.R[j].J]
		}

		// Se obtiene el tiempo total de la nueva ruta
	    newRoute.getWeigts(distances,base)

	    // Se verifica si la ruta nueva es mejor que la mejor actualmente
	    if newRoute.Time<route.Time {
	       	copy(route.R,newRoute.R)
	    	route.Time = newRoute.Time
	    }

	}

	// Al terminar las permutaciones se envia por el canal la mejor ruta
    c <- route
}

// Funcion de intercambio entre rutas
func (routes *RouteList) permutationRoutesBetween(distances [][]float64, demand []int, base int) {

	var(
		newRoutes RouteList
		min float64
		posMin, con int
	)
	min=0
	for i := range (*routes) {
		if (*routes)[i].Time<min {
			// Se selecciona una nueva ruta con tiempo minimo
			min = (*routes)[i].Time
			posMin = i
		}
	}


	// Nuemero de rutas
	lRoutes := len(*routes)
	// Se crea nueva lista de rutas para comparar
	routesCom := make([]*Route,len(*routes)-1,cap(*routes)-1)

	// Se inicializa las nuevas listas de rutas
	for i := 0; i < lRoutes; i++ {

		if i != posMin {
			// Se crea una nueva ruta a
			a := make([]Path,len((*routes)[i].R),cap((*routes)[i].R))
			// Se copia la ruta (*routes)[i].R en la nueva ruta a
			copy(a,(*routes)[i].R)
			// Se guarda la nueva ruta en las lista de rutas
			newRoutes = append(newRoutes,Route{R:a,Time:(*routes)[i].Time,Residue:(*routes)[i].Residue})
			// Se apunta la nueva ruta a comparar a la poscicion de memoria de la ruta
			routesCom[con] = &(*routes)[i]
			con++
		}
	}

	// Cantidad de nuevas rutas
	l := len(newRoutes)
	var change, cant int
	var perm []int
	// Cantidad de intercambios
	cant = 2000

	// Ciclo de intercambios
	for i := 0; i < cant; i++ {
		// Se crea el randomico
		t := time.Now()
		r := rand.New(rand.NewSource(int64(t.Nanosecond())))
		// Posicion a intercabiar
		change = r.Intn(l)

		for {
			// Vector randomico para elegir la segunda ruta a intercambiar
			perm = r.Perm(l)
			// Se verifica que las dos rutas a intercambiar no sean las mismas
			if change != perm[change] {
				break
			}
		}
		// Número de trayectos de la primera ruta a intercabiar
		l1 := len(newRoutes[change].R)
		t1 := time.Now()
		r1 := rand.New(rand.NewSource(int64(t1.Nanosecond())))
		// Posicion del trayecto a intercabiar en la primera ruta
		change1 := r1.Intn(l1)

		// Número de trayectos de la segunda ruta a intercambiar
		l2 := len(newRoutes[perm[change]].R)
		t2 := time.Now()
		r2 := rand.New(rand.NewSource(int64(t2.Nanosecond())))
		// Posicion del trayecto a intercambiar en la segunda ruta
		change2 := r2.Intn(l2)

		// Nuevas rutas a intercambiar
		route1:= &newRoutes[change]
		route2:= &newRoutes[perm[change]]
		// Rutas para comparar
		routeCom1:= routesCom[change]
		routeCom2:= routesCom[perm[change]]

		/*
			Se compara si las nueva rutas no superan las capacidades
		*/
	/*
		if((demand[routeCom2.R[change2].J]<=demand[routeCom1.R[change1].J]+routeCom1.Residue)&&(demand[routeCom1.R[change1].J]<=demand[routeCom2.R[change2].J]+routeCom2.Residue)) {

			// Se coloca el cliente a cambier de la ruta 2 en la ruta 1
			route1.R[change1].J = routeCom2.R[change2].J
			// Se coloca el nuevo tiempo de la trayectoria cambiada en la ruta 1
			route1.R[change1].Time =distances[route1.R[change1].I][route1.R[change1].J]
			// Se modifica el resuduo de la ruta 1
			route1.Residue = demand[routeCom1.R[change1].J]-demand[routeCom2.R[change2].J]+routeCom1.Residue

			// Se verifica si la trayectoria cambiada no es la ultima para cambiar la proxima trayectoria en la ruta 1
			if (change1<len(route1.R)-1 ){
				// Se cambia el inicio de la proxima trayectoria en la ruta 1
				route1.R[change1+1].I = routeCom2.R[change2].J
				// Se cambia el tiempo de la proxima trayectoria en la ruta 1
				route1.R[change1+1].Time =distances[route1.R[change1+1].I][route1.R[change1+1].J]
			}


			route2.R[change2].J = routeCom1.R[change1].J
			route2.R[change2].Time=distances[route2.R[change2].I][route2.R[change2].J]
			route2.Residue = demand[routeCom2.R[change2].J]-demand[routeCom1.R[change1].J]+routeCom2.Residue

			if (change2<len(route2.R)-1){
				route2.R[change2+1].I=routeCom1.R[change1].J
				route2.R[change2+1].Time=distances[route2.R[change2+1].I][route2.R[change2+1].J]
			}

			// Se calcula los nuevos tiempos de las rutas intercambiadas
			route1.getWeigts(distances,base)
			route2.getWeigts(distances,base)

			// Se hace permutaciones internas en las nuevas rutas
			route1.internalPermutations(distances,base)
			route2.internalPermutations(distances,base)

			// Se sacan los tiempos totales del antes y despues
			timeCom := routeCom1.Time+routeCom2.Time
			timeRta2 := route1.Time+route2.Time

			// Se compara los tiempos para verificar si la nueva ruta es mejor
			if(timeRta2<timeCom){



				// Se selecciona las nuevas rutas como la nuevas mejores
				*routeCom1 = *route1
				*routeCom2 = *route2


			// break

			}


			// Se verifica si no se a terminado con el numero de intercambios
			if i<cant-1 {

				min=0
				for i := range (*routes) {
					if (*routes)[i].Time<min {
						// Se selecciona una nueva ruta con tiempo minimo
						min = (*routes)[i].Time
						posMin = i
					}
				}


				// Se inicializa una nueva lista de nuevas rutas
				newRoutes = make(RouteList,0,0)
				for j := 0; j < lRoutes; j++ {
					if j != posMin {
						a := make([]Path,len((*routes)[j].R),cap((*routes)[j].R))
						copy(a,(*routes)[j].R)
						newRoutes = append(newRoutes,Route{R:a,Time:(*routes)[j].Time,Residue:(*routes)[j].Residue})
					}
				}
			}

		}else{
			 // //logPrintf("Sin cambios")
		}
	}

	// Timepo total final
	var TimeT float64
	for f := 0; f < lRoutes; f++ {
		//TimeT += (*routes)[f].Time
		}


}
*/
