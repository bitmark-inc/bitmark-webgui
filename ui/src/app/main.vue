<style>
  @font-face {
    font-family: 'Avenir Next LT Pro';
    font-style: normal;
    src: url(/static/fonts/AvenirNextLTPro-Regular.otf) format('opentype'), local('Avenir Next LT Pro');
  }

  body {
    font-family: "Avenir Next LT Pro", Helvetica, Arial, sans-serif;
    padding: 0 20px;
  }

  .container {
    margin-top: 70px;
  }

  .container h3,
  .container h4 {
    text-transform: uppercase;
  }
</style>

<template lang="pug">
  div
    nav-bar(v-if="this.$route.path !== '/login'")
    div.container
      router-view
</template>

<script>
  import Navbar from "./navbar.vue"
  export default {
    components: {
      "nav-bar": Navbar
    }
  }
</script>
