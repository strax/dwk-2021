const { cpus } = require("os")
const Path = require("path")

const HTMLPlugin = require("html-webpack-plugin")
const ReactRefreshPlugin = require("@pmmmwh/react-refresh-webpack-plugin")
const SRIPlugin = require("webpack-subresource-integrity")
const MiniCSSExtractPlugin = require("mini-css-extract-plugin")
const { BundleAnalyzerPlugin } = require("webpack-bundle-analyzer")

/**
 * @returns {import("webpack").Configuration}
 */
module.exports = env => {
  const isProduction = env.production ?? false
  const isProfiling = env.profiling ?? false
  return {
    mode: isProduction ? "production" : "development",
    devtool: "source-map",
    entry: ["./src/index.tsx", "./src/index.css"],
    output: {
      path: Path.join(__dirname, "assets"),
      crossOriginLoading: "anonymous",
    },
    parallelism: cpus().length,
    module: {
      rules: [
        {
          test: /\.(ts|js)x?$/,
          loader: "babel-loader",
          options: {
            envName: isProduction ? "production" : "development",
          },
          exclude: /node_modules/,
        },
        {
          test: /\.css$/,
          use: [
            MiniCSSExtractPlugin.loader,
            {
              loader: "css-loader",
              options: {
                url: url => !url.startsWith("/api/"),
              },
            },
          ],
        },
      ],
    },
    plugins: [
      new HTMLPlugin(),
      new SRIPlugin({
        hashFuncNames: ["sha512"],
        enabled: isProduction,
      }),
      !isProduction && new ReactRefreshPlugin(),
      new MiniCSSExtractPlugin(),
      isProfiling && new BundleAnalyzerPlugin(),
    ].filter(Boolean),
    resolve: {
      extensions: [".ts", ".tsx", ".js", ".jsx", ".json", ".css"],
    },
    devServer: {
      proxy: {
        "/api": "http://localhost:8081",
      },
    },
    stats: {
      version: true,
    },
  }
}
