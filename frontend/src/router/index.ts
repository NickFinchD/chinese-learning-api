import { createRouter, createWebHistory } from 'vue-router'

import AuthLayout from '@/layouts/AuthLayout.vue'
import DefaultLayout from '@/layouts/DefaultLayout.vue'

import LoginPage from '@/pages/LoginPage.vue'
import HomePage from '@/pages/HomePage.vue'

import { useAuthStore } from '@/stores/auth'

import CoursesPage from '@/pages/CoursesPage.vue'
import VocabularyPage from '@/pages/VocabularyPage.vue'
import SettingsPage from '@/pages/SettingsPage.vue'

import TestsPage from '@/pages/TestsPage.vue'
import GrammarTestPage from '@/pages/GrammarTestPage.vue'
import WordTrainingPage from '@/pages/WordTrainingPage.vue'

import TextsPage from '@/pages/TextsPage.vue'
import TextPage from '@/pages/TextPage.vue'

import AchievementsPage from '@/pages/AchievementsPage.vue'

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
    path: 'vocabulary',
    name: 'vocabulary',
    component: VocabularyPage,
  },
  {
    path: 'texts',
    name: 'texts',
    component: TextsPage,
  },
  {
    path: 'texts/:id',
    name: 'text',
    component: TextPage,
  },
  {
    path: 'tests',
    name: 'tests',
    component: TestsPage,
  },
  {
    path: 'tests/grammar',
    name: 'grammar-test',
    component: GrammarTestPage,
  },
  {
    path: 'tests/words',
    name: 'word-training',
    component: WordTrainingPage,
  },
  {
  path: 'lessons/:id',
  name: 'lesson',
  component: LessonPage,
},
  {
    path: 'achievements',
    name: 'achievements',
    component: AchievementsPage,
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