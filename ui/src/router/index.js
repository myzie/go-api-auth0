import Vue from 'vue'
import Router from 'vue-router'
import Home from '@/components/Home'
import Callback from '@/components/Callback'
import Test from '@/components/Test'

Vue.use(Router)

const router = new Router({
  mode: 'history',
  routes: [
    {
      path: '/home',
      name: 'home',
      component: Home
    },
    {
      path: '/auth_callback',
      name: 'callback',
      component: Callback
    },
    {
      path: '/test',
      name: 'test',
      component: Test
    },
    {
      path: '*',
      redirect: '/home'
    }
  ]
})

router.beforeEach((to, from, next) => {

  let isAuth = router.app.$auth.isAuthenticated()

  // eslint-disable-next-line
  console.log('isAuth', isAuth)

  if (to.name == 'callback') {
    next()
  } else if (to.name == 'test') {
    next()
  } else if (router.app.$auth.isAuthenticated()) {
    next()
  } else { // trigger auth0 login
    router.app.$auth.login()
  }
})

export default router
