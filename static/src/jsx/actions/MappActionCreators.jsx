var biudDispatcher = require('../dispatcher/biudDispatcher');
var PostConstants = require('../constants/MappConstants');

var ActionTypes = PostConstants.ActionTypes;

module.exports = {
	createPositions: function () {
		biudDispatcher.dispatch({
			type: ActionTypes.CLICK_NEW_POSITIONS
		});
	},
	createPolygons:function(){
		biudDispatcher.dispatch({
			type:ActionTypes.CLICK_NEW_POLYGONS
		});
	},
	uploadKml:function(){
		biudDispatcher.dispatch({
			type:ActionTypes.CLICK_UPLOAD_KML
		})
	}
};
