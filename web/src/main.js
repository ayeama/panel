import { createApp } from 'vue'
import App from './App.vue'
import router from './router'

import '@xterm/xterm/css/xterm.css'
import '@xterm/xterm/lib/xterm.js'

import tooltip from './directives/tooltip.js'

const app = createApp(App)

app.directive('tooltip', tooltip)

app.use(router)

app.mount('#app')
