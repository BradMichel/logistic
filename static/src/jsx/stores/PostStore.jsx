var injectTapEventPlugin = require("react-tap-event-plugin");
injectTapEventPlugin();
var React = require('react/addons');
var ActionCreator = require('../actions/biudBarActionCreators');
var mui = require('material-ui'),
AppBar = mui.AppBar;

module.export = {
	class default extends React.Component{
		handleTouchTap(task){
			ActionCreator.openNav();
		}
		render(){
			return(
				<AppBar title="biud"
					onLeftIconButtonTouchTap={this.handleTouchTap}
					/>

				)
		}
	}
}
