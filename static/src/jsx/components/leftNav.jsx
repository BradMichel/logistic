let injectTapEventPlugin = require("react-tap-event-plugin");
injectTapEventPlugin();
var biudDispatcher = require('../dispatcher/biudDispatcher');
var InterfaceConstants = require('../constants/InterfaceConstants');
var React = require('react/addons');
var mui = require('material-ui'),
  LeftNav = mui.LeftNav,
  MenuItem = mui.MenuItem,
  FlatButton=mui.FlatButton,
  ThemeManager = new mui.Styles.ThemeManager;

var MapActionCreator=require('../actions/MappActionCreators')
var emitter =  require('../stores/InterfaceStore');
class BiudNav extends React.Component{

	constructor(){
		super();
	}
	getChildContext() {
	    return {
	      muiTheme: ThemeManager.getCurrentTheme()
	    };
  	}

		componentDidMount(){
			this.onEventListener();
		}
		onEventListener(){
		var a = this;
		biudDispatcher.register(function(action){
			switch(action.type){
				case InterfaceConstants.ActionTypes.OPEN_NAV:
					a.refs.leftNav.toggle();
					emitter.emitChange();
				break;
			}
		})
		}
    createPositions(task){
      MapActionCreator.createPositions();
    }

	render(){
		this.onEventListener();
    var a=this;
    var menuItems = [
	      //{ route: 'get-started', text: 'Get Started' },
	      //{ route: 'customization', text: 'Customization' },
      { type: MenuItem.Types.SUBHEADER, text: 'Integrantes' },
      {
         type: MenuItem.Types.LINK,
         payload: 'https://www.facebook.com/andres.barreranavarro',
         text: 'Andres Barrera'
      },
      {
         type: MenuItem.Types.LINK,
         payload: 'https://www.facebook.com/angiehdez1893',
         text: 'Angie Hernandez'
      },
    ];
    //var FlatButtonPos=React.createElement(FlatButton,{style:{width:100+'%'},label:"Insertar posiciones",onTouchTap:function(event){a.createPositions();}},null)
    return(
				  //<LeftNav ref="leftNav" docked={false} header={FlatButtonPos}  menuItems={menuItems} />
          <LeftNav ref="leftNav" docked={false}  menuItems={menuItems} />
				)
	}
}

BiudNav.childContextTypes = {
  muiTheme: React.PropTypes.object
};
module.exports = BiudNav;
