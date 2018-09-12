<template>
  <div class="callback">Callback</div>
</template>

<script>
export default {
  name: 'callback',
  mounted() {
    this.$auth.handleAuthentication().then((data) => {
      // eslint-disable-next-line
      console.log('Callback', data)
      // this.$router.push({ name: 'home' })
    })
  }
}
</script>