var biudDispatcher = require('../dispatcher/biudDispatcher');
var InterfaceConstants = require('../constants/InterfaceConstants');
var EventEmitter = require('events').EventEmitter;
var assing = require('object-assign');

var ActionTypes = InterfaceConstants.ActionTypes;
var CHANGE_EVENT = 'change';

var InerfaceStore = assing({},EventEmitter.prototype,{

  emitChange: function() {
    this.emit(CHANGE_EVENT);
  },

  /**
   * @param {function} callback
   */
  addChangeListener: function(callback) {
    this.on(CHANGE_EVENT, callback);
  },

  /**
   * @param {function} callback
   */
  removeChangeListener: function(callback) {
    this.removeListener(CHANGE_EVENT, callback);
  }
})


module.exports = InerfaceStore;
