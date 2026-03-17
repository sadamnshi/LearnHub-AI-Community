const { defineConfig } = require('@vue/cli-service')
module.exports = defineConfig({
  transpileDependencies: true,
  devServer: {
    port: 8081, // 前端开发服务器端口（避免和后端 8080 冲突）
    proxy: {
      // 将所有 /api 开头的请求代理到后端
      '/api': {
        target: 'http://localhost:8080', // 后端 Gin 服务地址
        changeOrigin: true, // 改变请求头中的 Origin
        ws: true, // 支持 WebSocket
        pathRewrite: {
          // 不需要重写路径，因为后端就是 /api/xxx
        }
      }
    }
  }
})
