var React = require('react');

var biudBar = React.createClass({

  render : function() {
    return (
      <div>
        <h1>SUB COMPONENT!</h1>
        <h2>Message : {this.props.message}</h2>
      </div>
    );
  }

});

module.exports = biudBar;