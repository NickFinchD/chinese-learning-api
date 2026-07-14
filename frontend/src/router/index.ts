import { createRouter, createWebHistory } from 'vue-router'

import AuthLayout from '@/layouts/AuthLayout.vue'
import DefaultLayout from '@/layouts/DefaultLayout.vue'

import LoginPage from '@/pages/LoginPage.vue'
import HomePage from '@/pages/HomePage.vue'

import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(),

  routes: [
    {
      path: '/',
      component: AuthLayout,

      children: [
        {
          path: 'login',
          name: 'login',
          component: LoginPage,
          meta: {
            guest: true,
          },
        },
      ],
    },

    {
      path: '/app',
      component: DefaultLayout,

      meta: {
        requiresAuth: true,
      },

      children: [
        {
          path: '',
          name: 'home',
          component: HomePage,
        },
      ],
    },
  ],
})


router.beforeEach(async (to) => {
  const auth = useAuthStore()

  if (auth.user === null) {
    await auth.loadUser()
  }

  const requiresAuth = to.matched.some(route => route.meta.requiresAuth)
  const guestOnly = to.matched.some(route => route.meta.guest)

  if (requiresAuth && !auth.isAuthenticated) {
    return { name: 'login' }
  }

  if (guestOnly && auth.isAuthenticated) {
    return { name: 'home' }
  }
})
export default router