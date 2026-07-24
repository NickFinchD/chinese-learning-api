import { createApp } from 'vue'
import { createPinia } from 'pinia'

import './style.css'

import App from './App.vue'
import router from './router'
import { useThemeStore } from './stores/theme'
import { useFontsStore } from './stores/fonts'
import { usePaletteStore } from './stores/palette'

const app = createApp(App)

app.use(createPinia())
app.use(router)

useThemeStore().init()
useFontsStore().init()
usePaletteStore().init()

app.mount('#app')