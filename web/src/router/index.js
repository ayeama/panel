import { createRouter, createWebHistory } from 'vue-router'

import HomeView from '@/views/HomeView.vue'
import InstanceCreateView from '@/views/InstanceCreateView.vue'
import InstanceView from '@/views/InstanceView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView,
    },
    {
      path: '/instances/create',
      name: 'instanceCreate',
      component: InstanceCreateView,
    },
    {
      path: '/instances/:id',
      name: 'instance',
      component: InstanceView,
    },
    {
      path: '/:pathMatch(.*)*',
      redirect: { name: 'home' },
    },
  ],
})

export default router
