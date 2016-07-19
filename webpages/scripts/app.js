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
        .constant("BitmarkdConfig", {
            "testing": "/etc/bitmarkd-TESTING.conf",
            "bitmark": "/etc/bitmarkd-BITMARK.conf"
        })
        .constant("BitmarkCliSetupConfig", {
            identity: "admin",
            connect: "127.0.0.1:2130",
            description: "bitmark-webgui generated"
        })
        .constant("BitmarkCliConfig", {
            "testing": "/home/bitmark/config/bitmark-cli/bitmark-cli-TESTING.conf",
            "bitmark": "/home/bitmark/config/bitmark-cli/bitmark-cli-BITMARK.conf"
            // "testing": "/home/yuntai/testWebgui/config/bitmark-cli/bitmark-cli-testing.config"

        })
        .constant("BitmarkPayConfig", {
            "testing": "/home/bitmark/config/bitmark-pay/bitmark-pay-TESTING.xml",
            "bitmark": "/home/bitmark/config/bitmark-pay/bitmark-pay-BITMARK.xml"

            // "testing": "/home/yuntai/testWebgui/config/bitmark-pay/bitmark-pay-TESTING.xml"

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
                .when('/logout', {
                    templateUrl: 'views/logout.html',
                    controller: 'LogoutCtrl'
                })
                .otherwise({
                    redirectTo: '/login'
                });
        });
