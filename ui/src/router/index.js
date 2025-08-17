import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import ServerCreateView from '../views/ServerCreateView.vue'
import ServerView from '../views/ServerView.vue'
import SigninView from '../views/SigninView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView,
    },
    {
      path: '/servers/create',
      name: 'serverCreate',
      component: ServerCreateView,
    },
    {
      path: '/servers/:id',
      name: 'server',
      component: ServerView,
    },
    {
      path: '/signin',
      name: 'signin',
      component: SigninView,
    },
  ],
})

export default router
