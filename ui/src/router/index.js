import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import DashboardView from '../views/DashboardView.vue'
import ServerCreateView from '../views/ServerCreateView.vue'
import ServersView from '../views/ServersView.vue'
import ServerView from '../views/ServerView.vue'
import NodesView from '../views/NodesView.vue'
import UsersView from '../views/UsersView.vue'
import UserEditView from '../views/UserEditView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView,
    },
    {
      path: '/about',
      name: 'about',
      // route level code-splitting
      // this generates a separate chunk (About.[hash].js) for this route
      // which is lazy-loaded when the route is visited.
      component: () => import('../views/AboutView.vue'),
    },
    {
      path: '/dashboard',
      name: 'dashboard',
      component: DashboardView,
    },
    {
      path: '/servers/create',
      name: 'serverCreate',
      component: ServerCreateView,
    },
    {
      path: '/servers',
      name: 'servers',
      component: ServersView,
    },
    {
      path: '/servers/:id',
      name: 'server',
      component: ServerView,
    },
    {
      path: '/nodes',
      name: 'nodes',
      component: NodesView,
    },
    {
      path: '/users',
      name: 'users',
      component: UsersView,
    },
    {
      path: '/users/:id/edit',
      name: 'editUser',
      component: UserEditView,
    },
  ],
})

export default router
