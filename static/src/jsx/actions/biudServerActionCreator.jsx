var biudDispatcher = requiere('../dispatcher/biudDispatcher');
var PostConstants = requiere('../constants/PostConstants');

var ActionTypes = PostConstants.ActionTypes;

module.exports ={
	receiveAll: function(jsonContents){
		biudDispatcher.dispatch({
		type: ActionTypes.RECEIVE_CONTENTS,
		jsonContents: jsonContents
	});
	},

	receiveCreatedContent: function (createdContent) {
		biudDispatcher.dispatch({
			type: ActionTypes.RECEIVE_CREATED_CONTENT,
			jsonContent: createdContent
		});
	}
};