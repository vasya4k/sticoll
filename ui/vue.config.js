module.exports = {
  lintOnSave: false,
  runtimeCompiler: true,
  chainWebpack: config => {
    config
      .plugin('html')
      .tap(args => {
        args[0].chunksSortMode = 'none'

        return args
      })
  }
}
