var Vue = require('vue')
var VueRouter = require('vue-router')
Vue.use(VueRouter)

var Main = require('./app/main.vue')
var Login = require('./app/login.vue')

var routes = [
  {path: '/', component: Main},
  {path: '/login', component: Login}
]

var router = new VueRouter({routes})

var app = new Vue({
  router,
  el: '#main',
  render: h => h(Main)
})
