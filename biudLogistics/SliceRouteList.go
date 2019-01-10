package biudLogistics

import(
  "sort"
  "log"
)

func (solutions SliceRouteList)Len() int {  return len(solutions)}

func (a SliceRouteList)Swap(i,j int)  {a[i],a[j]=a[j],a[i]}

func (a SliceRouteList)Less(i,j int) bool  {return a[i].Pos<a[j].Pos}

func (s ByTimeSliceRouteList)Less(i,j int) bool  {return s.SliceRouteList[i].Comparison.Time<s.SliceRouteList[j].Comparison.Time}

func (s ByReliabilitySliceRouteList)Less(i,j int) bool {return s.SliceRouteList[i].Comparison.Reliability>s.SliceRouteList[j].Comparison.Reliability}

func (s ByDistanceCrowding)Swap(i,j int)  {s.SliceRouteList[i],s.SliceRouteList[j]=s.SliceRouteList[j],s.SliceRouteList[i]}

func (s ByDistanceCrowding)Less(i,j int) bool { return s.SliceRouteList[i].Comparison.Reliability+s.SliceRouteList[i].Comparison.Time<s.SliceRouteList[j].Comparison.Reliability+s.SliceRouteList[j].Comparison.Time}

func (solutions *SliceRouteList)GetInitialPoblation(vehicles *Vehicles,deposits *Deposits,clients *Clients)  (penalties [2*N][]bool, err error){
  TAG:="(solutions *SliceRouteList)GetInitialPoblation(vehicles *Vehicles,deposits *Deposits,clients *Clients)  (penalties [2*N][]bool)"

  n:=N/2
  for i := 0; i  < n; i++ {
    deposits.Penalize(P)  //  Penaliza los depositos aleatoriamente con una probabilidad del 30% de penalizacion por deposito

    penalties[i]=deposits.GetPenalties() //  Se guardan el vector depositos con sus penalizaciones

    (*solutions)[i].GetNeighbors("nearest",vehicles, deposits, clients) //  Crea la ruta con la logica del vecino mas cercano
    deposits.Penalize(P)  //  Hace una nueva penalizacion aleatoria de depositos
    penalties[n+i]=deposits.GetPenalties() //  Se guarda el vector depositos con sus penalizaciones

    (*solutions)[n+i].GetNeighbors("nearest",vehicles, deposits, clients) //  Crea la ruta con la logica del vecino mas cercano
  }
  err=solutions.ValidateClientsRepeat(TAG)
  return
}

func (solutions *SliceRouteList)GetSons(penalties *[2*N][]bool,vehicles *Vehicles,deposits *Deposits,clients *Clients)(err error){
  TAG:="(solutions *SliceRouteList)GetSons(penalties *[2*N][]bool,vehicles *Vehicles,deposits *Deposits,clients *Clients)"



  n:=N/2
  for i := 0; i < n; i++ {

    (*penalties)[N+i],(*penalties)[N+n+i]= deposits.CombinedPenalty((*penalties)[i] , (*penalties)[n+i])  //  Se combinan las penalizaciones de la ruta

    deposits.SetPenalties((*penalties)[N+i])
    (*solutions)[N+i]=RouteList{}
    (*solutions)[N+i].GetNeighbors("nearest",vehicles, deposits, clients)
    deposits.SetPenalties((*penalties)[N+n+i])
    (*solutions)[N+n+i]=RouteList{}
    (*solutions)[N+n+i].GetNeighbors("nearest",vehicles, deposits, clients)
    (*solutions)[N+i].LoadPaths()
    (*solutions)[N+n+i].LoadPaths()
    (*solutions)[N+i].Swap()  //  Se aplica la logica Swap a la ruta Hijo
    (*solutions)[N+i].InsertionLastClientLongerRoute()  //  Se inserta el ultimo usuario de la ruta mas larga en una poscion anterior
    (*solutions)[N+n+i].Swap()  //  Se aplica la logica Swap a la ruta Hijo
    (*solutions)[N+n+i].InsertionLastClientLongerRoute()  //  Se inserta el ultimo usuario de la ruta mas larga en una poscion anterior
  }
  err=solutions.ValidateClientsRepeat(TAG)
  return
}

func (solutions *SliceRouteList)GetNextGeneration(fronts map[int][]int) (err error) {
  TAG:="func (solutions *SliceRouteList)GetNextGeneration(fronts map[int][]int)"

  var newGeneration SliceRouteList
  newGeneration=make(SliceRouteList, 0)
  n:=0
  for _,front := range fronts {
    frontNeatSolutions:=solutions.GetSolutions(front)

    //log.Println("solutions",solutions)
    for _,solution := range frontNeatSolutions {
      newGeneration=append(newGeneration,solution)
      n++
      if n==N-1 { break }
    }
    if n==N-1 { break }
  }
  solutions=&newGeneration
  err=solutions.ValidateClientsRepeat(TAG)
  return
}

func (solutions *SliceRouteList)Dominance() (fronts map[int][]int, err error) {
  TAG:="(solutions *SliceRouteList)Dominance(fronts map[int][]int)"


  individuals,fronts,e:=solutions.FirstFront();
  if e!=nil{
    err=e
    return
  }

  solutions.OthersFronts(&individuals,&fronts)
  sort.Sort(solutions)
  solutions.StackingDistance(&individuals,&fronts)  //  Distancia de apilamiento
  err=solutions.ValidateClientsRepeat(TAG)
  return
}
func (solutions *SliceRouteList)FirstFront() (individuals Individuals, fronts map[int][]int,err error){
  TAG:="(solutions *SliceRouteList)FirstFront() (individuals Individuals, fronts map[int][]int)"

  ifront:=1
  fronts=make(map[int][]int)
  lenSolutions:=len(*solutions)
  individuals=make(Individuals, lenSolutions,lenSolutions)
  for i,solutionI := range (*solutions) {
    for j,solutionJ := range (*solutions) {
      if len(solutionI.Routes)==0 {log.Println("stop vacio",i,solutionI);break;}
      if len(solutionJ.Routes)==0  {log.Println("stop vacio",j,solutionJ);break;}
      solutionI.FirstFront(&solutionJ,j,&individuals[i])
    }


    if individuals[i].n==0 {
      solutionI.Pos=1
      fronts[ifront]=append(fronts[ifront],i)
    }
  }
  err=solutions.ValidateClientsRepeat(TAG)
  //log.Println("front 1",fronts[ifront])
  //log.Println("fronts ",fronts)
  return
}
func (solutions *SliceRouteList)OthersFronts(individuals *Individuals,fronts *map[int][]int) (err error) {
  TAG:="(solutions *SliceRouteList)OthersFronts(individuals *Individuals,fronts *map[int][]int)"


  newFront:=make([]int,0)
  ifront:=1

  for len((*fronts)[ifront])!=0{

    for _,posRoutesFront := range (*fronts)[ifront] {
      if len((*individuals)[posRoutesFront].p)>0 {

        for _,posRouteDoninated := range (*individuals)[posRoutesFront].p {
          (*individuals)[posRouteDoninated].n--
          if (*individuals)[posRouteDoninated].n==0 {
            (*solutions)[posRouteDoninated].Pos=ifront+1
            newFront=append(newFront,posRouteDoninated)
          }
        }
      }
    }
    ifront++

    (*fronts)[ifront]=newFront
    newFront=make([]int,0)
  }
  err=solutions.ValidateClientsRepeat(TAG)
  return
}
func (solutions *SliceRouteList)StackingDistance(individuals *Individuals, fronts *map[int][]int) (err error) {
  TAG:="(solutions *SliceRouteList)StackingDistance(individuals *Individuals, fronts *map[int][]int)"

  currentIndex:=0
  newSolutions:=new(SliceRouteList)
  for _,front := range *fronts {
    //istance:=0
    //lenSolutions:=len(*solutions)
    //lenFront:=len(front)

    if len(front)>0{

      sortedBasedOnObjectives:=new(SliceRouteList)

      for _,posSolutionFront := range front {
        //logPrintln(posSolutionFront)
        (*sortedBasedOnObjectives)=append((*sortedBasedOnObjectives),(*solutions)[posSolutionFront])
      }

      currentIndex=currentIndex+1
      lenObjetives:=(*solutions)[0].GetLenObjetives()
      for i := 0; i < lenObjetives; i++ {
        switch i {
        case 1:
          sort.Sort(ByTimeSliceRouteList{*sortedBasedOnObjectives})
          break;
        case 2:
          sort.Sort(ByReliabilitySliceRouteList{*sortedBasedOnObjectives})
          break;
        }
        iMax,iMin:=len(*sortedBasedOnObjectives)-1,0
        fMax:=(*sortedBasedOnObjectives)[iMax].GetValueIndexField("Comparison", i)
        (*sortedBasedOnObjectives)[iMax].SetValueIndexField("StackingDistance",i,100000)
        fMin:=(*sortedBasedOnObjectives)[iMin].GetValueIndexField("Comparison", i)
        (*sortedBasedOnObjectives)[iMin].SetValueIndexField("StackingDistance",i,100000)

        for j := 2; j < len((*sortedBasedOnObjectives))-1; j++ {  // por verificar limites
          vNex:=(*sortedBasedOnObjectives)[j+1].GetValueIndexField("Comparison", i)
            vPrev:=(*sortedBasedOnObjectives)[j-1].GetValueIndexField("Comparison", i)

            if fMax-fMin==0 {
                (*sortedBasedOnObjectives)[j].SetValueIndexField("StackingDistance",i,100000)
              }else{
                (*sortedBasedOnObjectives)[j].SetValueIndexField("StackingDistance",i,(vNex-vPrev)/(fMax-fMin))
              }
        }
      }
      *newSolutions=append(*newSolutions,(*sortedBasedOnObjectives)[0:]...)
    }
  }
  *solutions=*newSolutions
  err=solutions.ValidateClientsRepeat(TAG)
  return
}
func (solutions *SliceRouteList)GetSolutions(front []int)  (frontNeatSolutions SliceRouteList){
  //TAG:="(solutions *SliceRouteList)GetSolutions(front []int)  (frontNeatSolutions SliceRouteList)"
  //logPrintln(TAG)

  frontNeatSolutions=make(SliceRouteList, len(front))
  for i,id := range front {
    frontNeatSolutions[i]=(*solutions)[id]
  }

  sort.Sort(ByDistanceCrowding{frontNeatSolutions})
  return
}
func (solutions *SliceRouteList)LoadPaths()  {
  for pos := range *solutions {
    solution:=&(*solutions)[pos]
    solution.LoadPaths()
  }
}
func (solutions *SliceRouteList)ValidateClientsRepeat(tag string)  (err error){
  //TAG:="(solutions *SliceRouteList)ValidateClientsRepeat()"
  //log.Println(TAG)
  for _,solution := range *solutions {
    if err=solution.ValidateClientsRepeat(tag); err!=nil{
      break
    }
  }
  return
}
