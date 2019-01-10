var React=require('react/addons');
var update=require('react-addons-update');
var mui=require('material-ui'),
  ThemeManager=new mui.Styles.ThemeManager;
var {TextField,FloatingActionButton,Toolbar,ToolbarGroup,FlatButton,List,ListDivider,ListItem,TextField,Toggle,TimePicker}=mui;

var MapActionCreator=require('../actions/MappActionCreators')

var LocationCity=require('react-material-icons/icons/social/location-city');
var Transport=require('react-material-icons/icons/maps/local-shipping');
var Place=require('react-material-icons/icons/maps/place');
var Send=require('react-material-icons/icons/content/send');
var Shape=require('react-material-icons/icons/editor/border-outer');
var FileFileUpload=require('react-material-icons/icons/file/file-upload');
var Delete=require('react-material-icons/icons/action/delete')

class SubMenu extends React.Component {
  constructor(props) {
    super(props);
    this.state ={showDialogDemant:'none',showDialogTransport:'none',timeWindows:this.props.timeWindows,demand:this.props.demand,reliabilitys:this.props.reliabilitys,transport:this.props.transport};
    //this.setState({showDialogDemant:false});
  }
  getChildContext(){
    return{
      muiTheme:ThemeManager.getCurrentTheme()
    };
  }

  componentWillReceiveProps(nextProps){
    this.setState({timeWindows:nextProps.timeWindows,demand:nextProps.demand,reliabilitys:nextProps.reliabilitys,transport:nextProps.transport})
  }

  _showDemand(event){
    let a = this;
    var show = this.state.showDialogDemant;
    if(show=='inherit'){
       localStorage.setItem('demand',JSON.stringify(a.state.demand))
      this.setState({showDialogDemant:'none'});
    }else{
      this.setState({showDialogDemant:'inherit'});
    }
  }
  _showTransport(event){
    var show = this.state.showDialogTransport;
    if(show=='inherit'){
      this.setState({showDialogTransport:'none'});
    }else{
      this.setState({showDialogTransport:'inherit'});
    }
  }


  createPositions(task){
    MapActionCreator.createPositions();
  }
  createPolygons(task){
    MapActionCreator.createPolygons();
  }
  uploadKml(task){
    MapActionCreator.uploadKml();
  }


  render(){
    let demands=this.state.demand;
    let reliabilitys=this.state.reliabilitys
    let timeWindows=this.state.timeWindows
    var a=this;


    let styleButton= {
      backgroundColor:"transparent",
      fontSize: "40px",
      minWidth:"37px",
      position:"relative",
      width:"50px",
      marginLeft:"0px",
      marginRight:"0px"
    }
    let styleIcon={
      fill:"red"
    }


    return(
      <div>
      <Toolbar style={{backgroundColor:"transparent"}}>
      <ToolbarGroup float="left" >
      <FlatButton onTouchTap={this.props.deleteLocalStore}  style={styleButton} ><Delete color="red" ></Delete></FlatButton></ToolbarGroup>
      <ToolbarGroup float="right">
      <FloatingActionButton onTouchTap={a.uploadKml} secondary={true} ><FileFileUpload></FileFileUpload></FloatingActionButton></ToolbarGroup>
      </Toolbar>
      </div>
    )
  }
}
SubMenu.childContextTypes={
  muiTheme:React.PropTypes.object
};
module.exports=SubMenu;
