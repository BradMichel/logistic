var biudDispatcher = require('../dispatcher/biudDispatcher');
var InterfaceConstants = require('../constants/InterfaceConstants');

var ActionTypes = InterfaceConstants.ActionTypes;

module.exports = {
	openNav: function () {
		biudDispatcher.dispatch({
			type: ActionTypes.OPEN_NAV
		});
	},

}
