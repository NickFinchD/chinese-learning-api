import { createRouter, createWebHistory } from 'vue-router'

import AuthLayout from '@/layouts/AuthLayout.vue'
import DefaultLayout from '@/layouts/DefaultLayout.vue'

import LoginPage from '@/pages/LoginPage.vue'

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
        },
      ],
    },

    {
      path: '/app',

      component: DefaultLayout,

      children: [],
    },
  ],
})

export default router