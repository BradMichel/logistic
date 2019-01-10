var webpack = require('webpack');

module.exports = {
	context: __dirname,
	entry: "./src/jsx/app.jsx",

	output:{
		path: "./lib/js",
		filename: "app.js"
	},

	module: {
		loaders:[
			{
				test:/\.jsx$/,
				exclude: /(node_modules|pkg)/,
				loaders: ["babel-loader"]
			}
		]
	},
	 resolve : {
    extensions: ['', '.js', '.es6.js', '.jsx'],
    alias : {
      actions : __dirname + '/src/jsx/actions',
      constants : __dirname + '/src/jsx/constants',
      stores : __dirname + '/src/jsx/stores',
      componets : __dirname + '/src/jsx/componets',
      dispatcher : __dirname + '/src/jsx/dispatcher',
      utils : __dirname + '/src/jsx/utils'

    }
  }

}
