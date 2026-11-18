import { defineConfig } from "vite"
import vue from "@vitejs/plugin-vue"
import tailwindcss from "@tailwindcss/vite"

export default defineConfig({
  base: "./",
  plugins: [vue(), tailwindcss()],
  resolve: {
    alias: [{ find: /^vue$/, replacement: "vue/dist/vue.runtime.esm-bundler.js" }],
  },
  server: {
    port: 51345,
  },
})
