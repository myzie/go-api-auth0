import Vue from 'vue'
import Router from 'vue-router'
import Home from '@/components/Home'
import Callback from '@/components/Callback'
import Profile from '@/components/Profile'

Vue.use(Router)

// https://router.vuejs.org/guide/advanced/meta.html

const router = new Router({
  mode: 'history',
  routes: [
    {
      path: '/home',
      name: 'home',
      component: Home
    },
    {
      path: '/profile',
      name: 'profile',
      component: Profile,
      meta: { requiresAuth: true }
    },
    {
      path: '/auth_callback',
      name: 'callback',
      component: Callback
    },
    {
      path: '*',
      redirect: '/home'
    }
  ]
})

router.beforeEach((to, from, next) => {

  let loggedIn = router.app.$auth.isAuthenticated()

  if (to.matched.some(record => record.meta.requiresAuth)) {
    // this route requires auth, check if logged in
    // if not, redirect to login page.
    if (!loggedIn) {
      router.app.$auth.login()
    } else {
      next()
    }
  } else {
    next()
  }

})

export default router
