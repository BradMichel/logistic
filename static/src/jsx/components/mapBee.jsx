var React=require('react/addons');
var update=require('react-addons-update');
var $=require('jquery');
google = window.google;
var biudDispatcher=require('../dispatcher/biudDispatcher')
var ActionCreator=require('../actions/MappActionCreators');
var MappConstants=require('../constants/MappConstants');
var emmiter=require('../stores/MappStore');
let SubMenu = require('./subMenu');

var mui = require('material-ui'),
SelectField= mui.SelectField,
MenuItem=mui.MenuItem,
 ThemeManager = new mui.Styles.ThemeManager;

var PATH="api/";
var indicatorUpdate=0

/*
 * This is the modify version of:
 * https://developers.google.com/maps/documentation/javascript/examples/event-arguments
 *
 * Add <script src="https://maps.googleapis.com/maps/api/js"></script> to your HTML to provide google.maps reference
 */
class MapBee extends React.Component {

  constructor(props) {
    super(props);
    var body = document.body,html = document.documentElement;
    var height = Math.max( body.scrollHeight, body.offsetHeight,html.clientHeight, html.scrollHeight, html.offsetHeight );
    var width = Math.max( body.scrollWidth, body.offsetWidth,html.clientWidth, html.scrollWidth, html.offsetWidth );    
    let posOrigins = ( localStorage.getItem('posOrigins') == undefined ) ? 0 : localStorage.getItem('posOrigins')
    let posDestination = ( localStorage.getItem('posDestination') == undefined ) ? 0 : localStorage.getItem('posDestination')
    this.state={markers:[],origins:[],destinations:[],distances:[],duration:[],posOrigins:posOrigins,posDestination:posDestination,height:height,width:width}
  }

  getChildContext() {
    return {
      muiTheme: ThemeManager.getCurrentTheme()
    };
  }

  componentDidMount (rootNode) {
    var a=this;
    this.onEventListener();

    var mapOptions = {
      center: {lat: 7.107246, lng: -73.109496},
      zoom: 12
    };
    var map = new google.maps.Map(this.getDOMNode(),mapOptions); 
    this.setState({map: map});
    window.loadDistances=(distances,duration)=>{
      a.setState({distances:distances,duration:duration});
      let newState =a.state;
        a.setState(newState);

    }

    let kml=localStorage.getItem('kml')
    map.addListener('tilesloaded',()=>{
      if (kml!=undefined) {
      a.addMarkersFromKml(kml)
      a.getLocalStore()
      google.maps.event.clearListeners(map,'tilesloaded')
    }
    })

  }

  onEventListener(){
    var a = this;

    var uploadKml = (action)=>{
      switch (action.type) {
        case MappConstants.ActionTypes.CLICK_UPLOAD_KML:
          $('#file-upload').click();
          emmiter.emitChange();
          break;
        default:
      }
    }
    biudDispatcher.register(uploadKml);
  }

  getMatrix(callback) {

    let [a, service] = [this, new google.maps.DistanceMatrixService()]
    let [origins, destinations, posOrigins, posDestination] = [this.state.origins, this.state.destinations, parseInt(this.state.posOrigins), parseInt(this.state.posDestination)]
    let limitPositions = 5;
    let newOrigins = origins.slice().splice(posOrigins,limitPositions)
    let newDestinations = destinations.slice().splice(posDestination,limitPositions)
    newOrigins = newOrigins.map((item,i)=>{return item.position})
    newDestinations = newDestinations.map((item,i) => { return item.position}) 

    if (newDestinations.length === 0 || newOrigins.length == 0 ){
      return
    }

    service.getDistanceMatrix({
        origins: newOrigins,
        destinations: newDestinations,
        travelMode:google.maps.TravelMode.DRIVING,
      },function(response,status){
        if(status!== google.maps.DistanceMatrixStatus.OK){
          console.log(response,status);
          let delay = 10000;

          window.clock = setTimeout(a.getMatrix.bind(a),delay);
          indicatorUpdate++
          if (indicatorUpdate==3) {
            location.reload()
          }
        }else {
            console.log('origins from: ', posOrigins, ' to: ', posOrigins+newOrigins.length-1)
            console.log('destinate from: ', posDestination, ' to: ', posDestination+newDestinations.length-1)

            //indicatorUpdate=0
            let distances=(a.state.distances === undefined) ? [] : a.state.distances;
            let duration = ( a.state.duration === undefined ) ? [] : a.state.duration;

            let rows=response.rows;
            let originsRTA=response.originAddresses;
            let destinationsRTA=response.destinationAddresses;
            let newState=a.state;

            for(let i in rows){
              let distance=[];
              let elements=response.rows[i].elements;
              let posi = parseInt(posOrigins)+parseInt(i)
              for(let j in elements){
                let posj = parseInt(posDestination)+parseInt(j)
                if(elements[j].status=="OK"){                  
                  if (distances[posi] === undefined) {distances.push([])}
                  if( duration[posi] === undefined ){ duration.push([]) }
                  distances[posi].push(elements[j].distance.value);
                  duration[posi].push(elements[j].duration.value);
                }else{
                    console.log(elements[j].status);
                }
              }
                }

              newState=update(newState,{distances:{$set:distances},duration:{$set:duration}});
              posDestination = posDestination + newDestinations.length
              newState = update(newState,{posDestination:{$set:posDestination}})

              if ( posDestination >= destinations.length) {
                posOrigins = posOrigins + newOrigins.length
                newState = update(newState,{posOrigins:{$set:posOrigins},posDestination:{$set:0}})
            }

            localStorage.setItem('posOrigins',posOrigins)
            localStorage.setItem('posDestination',posDestination)
            let testObject={'distances':distances,'duration':duration};
            localStorage.setItem('testObject', JSON.stringify(testObject));

            a.setState(newState);
            if( posOrigins < origins.length-1 || posDestination < destinations.length-1 ){
              let delay = 2000
              console.log('nex',delay);

              window.clock = setTimeout(a.getMatrix.bind(a),delay);
            }else {
                console.log('Fin');
                //a.showDistancias()
              
              }

            }
          }
      )
  }

  showDistancias(){
    let a = this
    let tablaDistancias={}
    tablaDistancias[0]={}
    for (let i = 0; i <= a.state.demand.length - 1; i++) {
      tablaDistancias[0][a.state.demand[i][0]]={}
      for(let j=0; j<= a.state.demand.length-1;j++){
        let p=((j)-(j%20))/20;
        if (tablaDistancias[p]==undefined) {
          tablaDistancias[p]={}
        }
        if (tablaDistancias[p][a.state.demand[i][0]]==undefined) {
          tablaDistancias[p][a.state.demand[i][0]]={}
        }

        tablaDistancias[p][a.state.demand[i][0]][a.state.demand[j][0]]=a.state.distances[i][j];
      }
    }

    let p=((a.state.demand.length)-(a.state.demand.length%20))/20;
    //console.log(tablaDistancias)
    for (let i = 0; i <= p; i++) {
    //console.log('Tabla '+(i+1))
    //console.table(tablaDistancias[i])
    }
  }
 
  addMarker(location,map,label,title){
    let marker=new google.maps.Marker({
      position:location,
      label:label,
      title:title,
      draggable:true,
      map:map,
      animation: google.maps.Animation.DROP
    });
    return marker;
  }

  getDOMNode(){
    let map = document.getElementsByClassName('map-gic')[0];
    return map;
  }


  deleteLocalStore(){
    let  TAG='deleteLocalDB'
    console.log( TAG)
    localStorage.clear();
  }



  loadKml(event,a){
    console.log('loadKml');
    var input = document.getElementById('file-upload');
    var file = input.files[0];

    var fr = new FileReader();
         fr.onload = receivedText;
         fr.readAsText(file);

     function receivedText() {
         fr = new FileReader();
         fr.onload = receivedBinary;
         fr.readAsBinaryString(file);
     }

     function receivedBinary() {
          a.addMarkersFromKml(fr.result)
          a.getMatrix()
     }
  }

  addMarkersFromKml(kml){
    console.log('addMarkersFromKml')

    let [a,pos,num]=[this,0,0]
    let state=a.state
    let map=state.map
    let markers = (a.state.markers!=undefined) ? a.state.markers.slice():[];
    localStorage.setItem('kml',kml)

    $(kml).find("Folder").each((i,v)=>{
      let typeLayer="destinate"
      $(v).children().each((i,v)=>{
        let tagName = $(v).prop("tagName")
        switch (tagName) {
          case "NAME":
            if ($(v).text()==="Centroides") {
              typeLayer="origin"
            }
            break;
          case "PLACEMARK":

            let title=$(v).find("name").text();
            let coords= $(v).find("coordinates").text().split(",");
            let LatLng=new google.maps.LatLng(parseFloat(coords[1]),parseFloat(coords[0]));
            let marker = a.addMarker(LatLng,map,''+num++,title);
            markers.push(marker);
            if(typeLayer == 'origin'){
              state=update(state,{markers:{$push:[marker]},origins:{$push:[{position:LatLng,label:title,idMarkers:markers.length}]}})
            }else{
              state=update(state,{markers:{$push:[marker]},destinations:{$push:[{position:LatLng,label:title,idMarkers:markers.length}]},origins:{$push:[{position:LatLng,label:title,idMarkers:markers.length}]}})
            }
            
            pos++

            break;
          default:

        }

      })
    })

    a.setState(state);

  }

  getLocalStore(){
    let a = this;
    var retrievedObject = JSON.parse(localStorage.getItem('testObject'));

  setTimeout(()=>{
  if(retrievedObject!=null){
    loadDistances(retrievedObject.distances,retrievedObject.duration);    
    }
    a.getMatrix(a.getMatrix);
              },3000)
  }


  render () {
      var a=this; 

      let [solutions,solutionSelected,listRouteSelected,routeSelected,itemSolutions,itemListRoutes,itemRoutes]= [a.state.solutions,a.state.solutionSelected,a.state.listRouteSelected,a.state.routeSelected,[],[],[]]

      return (<div>
        <div className='map-gic' style={{position:'relative',height:(a.state.height-64-56-16-70)+'px',width:100+'%',}}></div>

        <SubMenu deleteLocalStore={a.deleteLocalStore} send={(event)=>{a.send(event,a)}} demand={this.state.demand} reliabilitys={this.state.reliabilitys} transport={this.state.transport} timeWindows={this.state.timeWindows}></SubMenu>
        <input type='file' id='file-upload' style={{display:'none'}} onChange={(event)=>{a.loadKml(event,a)}}></input>
        </div>);
      }
    }

  MapBee.childContextTypes={
    muiTheme:React.PropTypes.object
  };

module.exports=MapBee;
