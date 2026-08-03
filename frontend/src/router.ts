import { createRouter, createWebHistory } from 'vue-router'
import { isSignedIn, waitForClerk } from './auth'
import FeedView from './views/FeedView.vue'
import LoginView from './views/LoginView.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', name: 'feed', component: FeedView },
    { path: '/login', name: 'login', component: LoginView },
  ],
})

router.beforeEach(async (to) => {
  await waitForClerk()
  if (to.name !== 'login' && !isSignedIn()) {
    return { name: 'login', query: { next: to.fullPath } }
  }
  if (to.name === 'login' && isSignedIn()) {
    return { name: 'feed' }
  }
  return true
})

export default router
