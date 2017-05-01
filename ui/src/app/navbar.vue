<style scoped>
  ul,
  li {
    margin: 0;
    padding: 0;
    list-style: none;
  }

  .navbar {
    top: 0;
    left: 0;
    width: calc(100% - 70px);
    position: fixed;
    margin: 0 20px;
    padding: 0 15px;
    min-height: 50px;
    border-bottom: 0.5px solid #b2c1d8;
    background-color: rgba(255, 255, 255, 0.8);
  }

  .menu {
    float: left;
    height: 46px;
  }

  .menu .trademark {
    margin: 0;
    padding: 0 16px;
    line-height: 50px;
  }

  .menu.right {
    float: right;
  }

  .menu a.active,
  .menu a:hover {
    border-bottom: 3.5px solid #b2c1d8
  }

  .menu a {
    display: inline-block;
    line-height: 0;
    box-sizing: border-box;
    padding: 23px;
    color: black;
    text-decoration: none;
  }
</style>

<template lang="pug">
div.navbar
  div.menu
    p.trademark Bitmark
  div.menu
    router-link(to="/node") Node
    router-link(to="/console") Console
  div.right.menu
    a(href="") Switch Network
    a(href="", @click="this.logout") Logout

</template>

<script>
  import {
    setCookie
  } from "../utils"
  export default {
    methods: {
      logout(e) {
        e.preventDefault();
        setCookie("bitmark-webgui", "", 0)
        this.$router.push({path: "/login"})
      }
    }
  }
</script>
