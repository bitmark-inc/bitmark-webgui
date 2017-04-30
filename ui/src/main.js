var Vue = require('vue')
var VueRouter = require('vue-router')
Vue.use(VueRouter)

var Main = require('./app/main.vue')
var Login = require('./app/login.vue')
var Node = require('./app/node.vue')
var Console = require('./app/console.vue')

import {getCookie} from "./utils"

var routes = [
  {path: '/', component: Main, redirect: '/node'},
  {path: '/login', component: Login},
  {path: '/node', component: Node},
  {path: '/console', component: Console}
]

var router = new VueRouter({routes, linkActiveClass: "active"})

router.beforeEach((to, from, next) => {
  if (to.path === "/login") {
    if (getCookie("bitmark-webgui") !== "") {
      next({path: "/"})
    } else {
      next()
    }
  } else if (getCookie("bitmark-webgui") === "") {
    next({path: "/login", query: {redirect: to.fullPath}})
  } else {
    next()
  }
})

var app = new Vue({
  router,
  el: '#main',
  render: h => h(Main)
})
