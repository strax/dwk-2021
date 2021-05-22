module.exports = function (api) {
  const isDevelopment = api.env("development")
  const plugins = []
  if (isDevelopment) {
    plugins.push("react-refresh/babel")
  }
  return {
    presets: [
      [
        "@babel/preset-react",
        {
          runtime: "automatic",
          development: isDevelopment,
        },
      ],
      [
        "@babel/preset-typescript",
        {
          allowNamespaces: true,
          allowDeclareFields: true,
          onlyRemoveTypeImports: true,
        },
      ],
    ],
    plugins,
  }
}
