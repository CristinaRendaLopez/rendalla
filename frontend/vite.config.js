import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";

export default defineConfig({
  plugins: [vue()],  // Agrega el plugin de Vue
  build: {
    rollupOptions: {
      output: {
        entryFileNames: "assets/main.js",
        chunkFileNames: "assets/[name].js",
        assetFileNames: "assets/[name].[ext]",
      },
    },
  },
});
