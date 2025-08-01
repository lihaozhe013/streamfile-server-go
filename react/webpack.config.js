const path = require('path');

module.exports = {
  entry: './src/MarkdownViewer.tsx',
  output: {
    path: path.resolve(__dirname, 'dist'),
    filename: 'markdown-viewer.js',
    library: 'MarkdownViewer',
    libraryTarget: 'umd',
    globalObject: 'this'
  },
  module: {
    rules: [
      {
        test: /\.tsx?$/,
        exclude: /node_modules/,
        use: {
          loader: 'ts-loader',
          options: {
            transpileOnly: true
          }
        }
      },
      {
        test: /\.css$/,
        use: ['style-loader', 'css-loader']
      }
    ]
  },
  externals: {
    'react': 'React',
    'react-dom': 'ReactDOM'
  },
  resolve: {
    extensions: ['.tsx', '.ts', '.js', '.jsx']
  },
  devtool: 'source-map'
};
