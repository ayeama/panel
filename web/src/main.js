import { createApp } from 'vue'
import App from './App.vue'
import router from './router'

import '@xterm/xterm/css/xterm.css'
import '@xterm/xterm/lib/xterm.js'

const app = createApp(App)

app.use(router)

app.mount('#app')
