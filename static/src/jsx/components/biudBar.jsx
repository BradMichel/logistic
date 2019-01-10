var React = require('react/addons');
var ActionCreator = require('../actions/biudBarActionCreators');
var mui = require('material-ui'),
  ThemeManager = new mui.Styles.ThemeManager();
  var {AppBar,IconButton}=mui;
var SearchIcon = require('react-material-icons/icons/action/search');
var emitter =  require('../stores/InterfaceStore');
class BiudBar extends React.Component{

	constructor(){
		super();
	}
	getChildContext() {
	    return {
	      muiTheme: ThemeManager.getCurrentTheme()
	    };
  	}

	handleTouchTap(task){
		ActionCreator.openNav();
	}

	render(){
		return(
			<AppBar title="Bee UIS"
				onLeftIconButtonTouchTap={this.handleTouchTap}
				iconElementRight={<IconButton><SearchIcon color="white" /></IconButton>}
				/>

			)
	}
}

BiudBar.childContextTypes = {
  muiTheme: React.PropTypes.object
};



module.exports = BiudBar;
