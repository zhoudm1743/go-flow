import vue from '@vitejs/plugin-vue'
import { defineConfig } from 'vite'
import AutoImport from 'unplugin-auto-import/vite'
import Components from 'unplugin-vue-components/vite'
import { ElementPlusResolver } from 'unplugin-vue-components/resolvers'
import monacoEditorPlugin from 'vite-plugin-monaco-editor'
import {resolve} from 'path'

const pathSrc = resolve(__dirname, 'src')

export default defineConfig({
  resolve: {
    alias: {
      '~/': `${pathSrc}/`,
      '@/': `${pathSrc}/` ,
      }
  },

  server: {
    proxy: {
      '/admin': {
        target:  'http://127.0.0.1:8080',
        changeOrigin: true,
      }
    }
  },
 
  // ...
  plugins: [
    vue(),
    // ...
    AutoImport({
      resolvers: [ElementPlusResolver()],
    }),
    Components({
      resolvers: [ElementPlusResolver()],
    }),
    monacoEditorPlugin({
      languageWorkers:['editorWorkerService', 'json']
    })
  ],
  base: './'
})

