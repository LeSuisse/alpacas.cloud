const path = require("path");
const { CleanWebpackPlugin } = require('clean-webpack-plugin');
const MiniCssExtractPlugin = require('mini-css-extract-plugin');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const CopyPlugin = require('copy-webpack-plugin');

module.exports = {
    entry: {
        "app": "./app.js"
    },
    output: {
        path: path.resolve(__dirname, "dist/"),
        filename: "[name]-[contenthash].js",
        publicPath: "/assets/"
    },
    module: {
        rules: [
            {
                test: /\.css$/i,
                use: [MiniCssExtractPlugin.loader, 'css-loader'],
            }
        ],
    },
    plugins: [
        new CleanWebpackPlugin(),
        new MiniCssExtractPlugin({
            filename: '[name]-[contenthash].css'
        }),
        new HtmlWebpackPlugin({
            template: "index.html"
        }),
        new CopyPlugin([
            {from: "openapi.json", to: "openapi.json"}
        ]),
    ]
};
