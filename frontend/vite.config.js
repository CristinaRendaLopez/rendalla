export default {
    build: {
      rollupOptions: {
        output: {
          entryFileNames: 'assets/main.js',
          chunkFileNames: 'assets/[name].js',
          assetFileNames: 'assets/[name].[ext]'
        }
      }
    }
  };
  