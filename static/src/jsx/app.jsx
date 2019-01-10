(function () {
let React = require('react/addons');
let injectTapEventPlugin = require("react-tap-event-plugin");
let BiudBar = require('./components/biudBar');
let LeftNav = require('./components/leftNav');
let MapBee = require('./components/mapBee');


//setTimeout(function () { window.location.reload(); }, 120000);

window.React = React;
injectTapEventPlugin();
React.render(<div><BiudBar/><LeftNav/><MapBee/></div>, document.body);
})();
