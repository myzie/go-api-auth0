<template>
  <div class="dashboard">
    <nav class="navbar navbar-dark bg-dark">
      <a class="navbar-brand" href="#">
        <img src="https://a.storyblok.com/f/39898/1024x1024/dea4e1b62d/vue-js_logo-svg.png" width="40" height="40">
      </a>
      <div v-if="$auth.user">
        <img :src="$auth.user.picture" width="30" height="30">
        <span class="text-muted font-weight-light px-2">{{$auth.user.name}}</span>
        <button type="button" class="btn btn-outline-secondary btn-sm" @click="$auth.logout()">Logout</button>
      </div>
    </nav>
  
    <div class="jumbotron" v-if="$auth.user">
      <div class="container">
        <h1 class="display-4">Hello, {{$auth.user.name}}!</h1>
        <pre>{{$auth.user}}</pre>
        <pre>{{message}}</pre>
      </div>
    </div>
    <div v-else class="jumbotron">
      <div class="container">
        <h1>You are not logged in</h1>
        <a href="/profile">Login</a>
      </div>
    </div>

    <div class="container">
    </div>
  </div>
</template>

<script>
import axios from 'axios'

export default {
  data () {
    return {
      message: null
    }
  },
  mounted() {
    axios.get('http://localhost:8080/')
      .then(resp => {
        this.message = resp.data
      })
  }
}
</script>

<style scoped>
@import url('https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0/css/bootstrap.min.css');

.btn-primary {
  background: #468f65;
  border: 1px solid #468f65;
}
.card {
  text-decoration: none;
  color: #000;
}
</style>