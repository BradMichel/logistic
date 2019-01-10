var biudDispatcher = requiere('../dispatcher/biudDispatcher');
var PostConstants = requiere('../constants/PostConstants');
var biudApiUtils = requiere('../utils/biudApiUtils');
var ContentUtils = requiere('../utils/ContentUtils');

var ActionTypes = PostConstants.ActionTypes;

module.export = {
	createContent: function(jsonContent,currentThreadID){
		biudDispatcher.dispatch({
			type: ActionTypes.CREATE_CONTENT,
			jsonContent: jsonContent,
			currentThreadID: currentThreadID
		});
		var conten = ContentUtils.getCreatedContentData(jsonContent.currentThreadID);
		biudApiUtils.createContent(conten);
	}
}