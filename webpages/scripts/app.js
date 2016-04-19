// Copyright (c) 2014-2016 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

'use strict';

/**
 * @ngdoc overview
 * @name bitmarkMgmtApp
 * @description
 * # bitmarkMgmtApp
 *
 * Main module of the application.
 */
var app = angular
  .module('bitmarkMgmtApp', [
    'ngCookies',
    'ngResource',
    'ngRoute',
    'ui.bootstrap'
  ])
        .constant("BitmarkProxyURL", {
            "testing": "https://spoon.test.bitmark.com:17555/rpc-call",
            "bitmark": "https://spoon.live.bitmark.com:17555/rpc-call"
        })
        .constant("ProxyTemp", {
            Username: "No-need-username",
            Password: "No-need-password",
            URL: "",
            Fee: "0.0002",
            Address: ""
        })
        .config(function ($routeProvider, $httpProvider) {
    $httpProvider.defaults.withCredentials = true;
    delete $httpProvider.defaults.headers.common["X-Requested-With"];

    $routeProvider
      .when('/', {
        templateUrl: 'views/main.html',
        controller: 'MainCtrl'
      })
      .when('/edit', {
        templateUrl: 'views/edit.html',
        controller: 'EditCtrl'
      })
      .when('/system', {
        templateUrl: 'views/system.html',
        controller: 'SystemCtrl'
      })
      .when('/login', {
        templateUrl: 'views/login.html',
        controller: 'LoginCtrl'
      })
      .otherwise({
        redirectTo: '/'
      });

  });
