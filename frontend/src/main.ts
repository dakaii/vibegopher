import { clerkPlugin } from '@clerk/vue'
import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import './style.css'

const PUBLISHABLE_KEY = import.meta.env.VITE_CLERK_PUBLISHABLE_KEY || ''

const app = createApp(App)
if (PUBLISHABLE_KEY) {
  app.use(clerkPlugin, { publishableKey: PUBLISHABLE_KEY })
} else {
  console.warn('VITE_CLERK_PUBLISHABLE_KEY is not set — Clerk sign-in will be unavailable')
}
app.use(router)
app.mount('#app')
