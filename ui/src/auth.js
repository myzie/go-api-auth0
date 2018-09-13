import auth0 from 'auth0-js'
import Vue from 'vue'
import { AUTH_CONFIG } from './auth0-variables'

let webAuth = new auth0.WebAuth({
  // we will use the api/v2/ to access the user information as payload
  // audience: 'https://' + 'your_auth0_domain' + '/api/v2/', 
  domain: AUTH_CONFIG.domain,
  clientID: AUTH_CONFIG.clientId,
  redirectUri: AUTH_CONFIG.callbackUrl,
  responseType: 'token id_token',
  scope: 'openid profile email'
})

let auth = new Vue({
  computed: {
    token: {
      get: function() {
        return localStorage.getItem('id_token')
      },
      set: function(id_token) {
        localStorage.setItem('id_token', id_token)
      }
    },
    accessToken: {
      get: function() {
        return localStorage.getItem('access_token')
      },
      set: function(accessToken) {
        localStorage.setItem('access_token', accessToken)
      }
    },
    expiresAt: {
      get: function() {
        // eslint-disable-next-line
        console.log('get expiresAt', localStorage.getItem('expires_at'))
        return localStorage.getItem('expires_at')
      },
      set: function(expiresIn) {
        let expiresAt = JSON.stringify(expiresIn * 1000 + new Date().getTime())
        localStorage.setItem('expires_at', expiresAt)
        // eslint-disable-next-line
        console.log('set expiresAt', expiresAt)
      }
    },
    user: {
      get: function() {
        let user = JSON.parse(localStorage.getItem('user'))
        // eslint-disable-next-line
        console.log('get user', user)
        return user
      },
      set: function(user) {
        localStorage.setItem('user', JSON.stringify(user))
      }
    }
  },
  methods: {
    login() {
      webAuth.authorize()
    },
    logout() {
      return new Promise(() => { 
        localStorage.removeItem('access_token')
        localStorage.removeItem('id_token')
        localStorage.removeItem('expires_at')
        localStorage.removeItem('user')
        webAuth.authorize()
        // this.$router.push({ name: 'home' })
      })
    },
    isAuthenticated() {
      return new Date().getTime() < this.expiresAt
    },
    handleAuthentication() {
      return new Promise((resolve, reject) => {  
        webAuth.parseHash((err, authResult) => {
          // eslint-disable-next-line
          console.log('webAuth hash', err, authResult)
          if (authResult && authResult.accessToken && authResult.idToken) {
            this.expiresAt = authResult.expiresIn
            this.accessToken = authResult.accessToken
            this.token = authResult.idToken
            this.user = authResult.idTokenPayload
            // eslint-disable-next-line
            console.log('GOOD!', this.expiresAt, this.accessToken, this.token, this.user)
            resolve()
          } else if (err) {
            this.logout()
            reject(err)
          }
        })
      })
    }
  }
})

export default {
  install: function(Vue) {
    Vue.prototype.$auth = auth
  }
}
