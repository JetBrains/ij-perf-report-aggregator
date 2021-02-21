// eslint-disable-next-line @typescript-eslint/no-var-requires
const path = require("path")

module.exports = {
  integrity: process.env.NODE_ENV === "production",
  outputDir: path.resolve(__dirname, "../cmd/frontend/kodata/v2"),
  devServer: {
    host: "localhost",
  },
  chainWebpack: (config) => {
    config.devtool("source-map")
    config.optimization.splitChunks({
      cacheGroups: {
        elementPlus: {
          name: "chunk-element-plus",
          test: /[\\/]node_modules[\\/]element-plus/,
          priority: -5,
          chunks: "initial",
          enforce: true,
        },
        vendors: {
          name: "chunk-vendors",
          test: /[\\/]node_modules[\\/]/,
          priority: -10,
          chunks: "initial",
          enforce: true,
        },
      },
    })
  },
}
