<style scoped>
  h3 {
    font-size: 18px;
    font-weight: bold;
  }

  .container {
    margin: 0 auto;
    text-align: center;
    font-family: "Avenir Next LT Pro", Helvetica, Arial, sans-serif;
    text-transform: uppercase;
  }

  .login-form {
    margin-top: 100px;
    display: inline-block;
  }

  .login-form .password,
  .login-form .login-btn {
    font-size: 16px;
    text-transform: uppercase;
    text-align: center;
    width: 400px;
    height: 50px;
    border: none;
    margin-top: 20px;
  }

  .login-form .password {
    background-color: #edf0f4;
    font-style: italic;
  }

  .login-form .login-btn {
    color: white;
    padding: 1px;
    background-color: #0060f2;
  }

  .login-form .login-btn:hover {
    background-color: black;
  }
</style>

<template lang="pug">
    div.container
      div.login-form
        h3 login bitmark node web
        div
          input.password(type="password", v-model="password" placeholder="PASSWORD")
        div
          button.login-btn(type="button", @click="this.login") login
</template>

<script>
  const axios = require("axios");
  export default {
    methods: {
      login() {
        console.log("login", this.password)
        axios.post('/api/login', {
            password: this.password
          })
          .then(function(response) {
            console.log(response);
          })
          .catch(function(error) {
            console.log(error);
          });
      }
    },

    data() {
      return {
        password: ""
      }
    }
  }
</script>
