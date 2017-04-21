// Copyright (c) 2014-2016 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

'use strict';

/**
 * @ngdoc overview
 * @name bitmarkWebguiApp
 * @description
 * # bitmarkWebguiApp
 *
 * Main module of the application.
 */
var app = angular
  .module('bitmarkWebguiApp', [
    'ngCookies',
    'ngResource',
    'ngRoute',
    'ui.bootstrap'
  ])
        .constant("BitmarkProxyURL", {
            "testing": "https://spoon.live.bitmark.com:17011/rpc-call",
            "bitmark": "https://spoon.test.bitmark.com:17555/rpc-call"
        })
        .constant("ProxyTemp", {
            Username: "No-need-username",
            Password: "No-need-password",
            URL: "",
            Fee: "0.0002",
            Address: ""
        })
        .constant("BitmarkdNetwork", {
            "testing": "testing",
            "bitmark": "bitmark"
        })
        .constant("BitmarkdConfig", {
            "testing": "/etc/bitmarkd-TESTING.conf",
            "bitmark": "/etc/bitmarkd-BITMARK.conf"
        })
        .config(function ($routeProvider, $httpProvider) {
            $httpProvider.defaults.withCredentials = true;
            delete $httpProvider.defaults.headers.common["X-Requested-With"];

            $routeProvider
                .when('/chain', {
                    templateUrl: 'views/chain.html',
                    controller: 'ChainCtrl'
                })
                .when('/main', {
                    templateUrl: 'views/main.html',
                    controller: 'MainCtrl'
                })
                .when('/edit', {
                    templateUrl: 'views/edit.html',
                    controller: 'EditCtrl'
                })
                .when('/console', {
                    templateUrl: 'views/console.html',
                    controller: 'ConsoleCtrl'
                })
                .when('/network', {
                    templateUrl: 'views/network.html',
                    controller: 'NetworkCtrl'
                })
                .when('/', {
                    templateUrl: 'views/login.html',
                    controller: 'LoginCtrl'
                })
                .when('/login', {
                    templateUrl: 'views/login.html',
                    controller: 'LoginCtrl'
                })
                .otherwise({
                    redirectTo: '/login'
                });
        });
