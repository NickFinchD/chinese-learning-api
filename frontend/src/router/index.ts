import { createRouter, createWebHistory } from 'vue-router'

import AuthLayout from '@/layouts/AuthLayout.vue'
import DefaultLayout from '@/layouts/DefaultLayout.vue'

import LoginPage from '@/pages/LoginPage.vue'
import HomePage from '@/pages/HomePage.vue'

import { useAuthStore } from '@/stores/auth'

import CoursesPage from '@/pages/CoursesPage.vue'
import ReviewPage from '@/pages/ReviewPage.vue'
import SavedWordsPage from '@/pages/SavedWordsPage.vue'
import SettingsPage from '@/pages/SettingsPage.vue'

import CoursePage from '@/pages/CoursePage.vue'
import LessonPage from '@/pages/LessonPage.vue'

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
  {
    path: 'courses',
    name: 'courses',
    component: CoursesPage,
  },
  {
    path: 'courses/:id',
    name: 'course',
    component: CoursePage,
},
  {
    path: 'review',
    name: 'review',
    component: ReviewPage,
  },
  {
    path: 'saved',
    name: 'saved-words',
    component: SavedWordsPage,
  },
  {
  path: 'lessons/:id',
  name: 'lesson',
  component: LessonPage,
},
  {
    path: 'settings',
    name: 'settings',
    component: SettingsPage,
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