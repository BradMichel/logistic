var biudDispatcher=require('../dispatcher/biudDispatcher');
var MappConstants=require('../constants/MappConstants');
var EventEmitter=require('events').EventEmitter;
var assing=require('object-assign');

var ActionTypes=MappConstants.ActionTypes;
var CHANGE_EVENT = 'change';

var MappStore=assing({},EventEmitter.prototype,{
  emitChange:function(){
    this.emit(CHANGE_EVENT);
  },
  addChangeListener:function(callback){
    this.on(CHANGE_EVENT,callback);
  },
  removeChangeListener:function(callback){
    this.removeListener(CHANGE_EVENT,callback);
  }
})
module.exports=MappStore;
