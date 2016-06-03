// Copyright (c) 2014-2016 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

'use strict';

/**
 * @ngdoc function
 * @name bitmarkWebguiApp.controller:LoginCreateCtrl
 * @description
 * # LoginCtrl
 * Controller of the bitmarkWebguiApp
 */
angular.module('bitmarkWebguiApp')
    .controller('LoginCreateCtrl', ['$scope', '$timeout', '$location', '$cookies', 'httpService', function ($scope, $timeout, $location, $cookies, httpService) {
        $scope.panelConfig = {
            showPart: 1
        };

        $scope.generateConfig = {
            running: false,
            msg: [
                "checking bitcoin ...",
                "generating key pair ...",
            ]
        };
        $scope.generate = function(){
            $scope.generateConfig.running = true;
            // TODO:
            // save chain to bitmark
            //check bitcoin
            // create key pair
            // encrypt bitmark wallet

            $timeout(function(){
                // mock when generate is done
                $scope.generateConfig.running = false;
                $scope.panelConfig.showPart = 2;
            }, 2*1000);

        };


        $scope.setPassCode = function(){
            $scope.panelConfig.showPart = 3;
        };

        $scope.done = function(){
            $location.path("/main");
        };

  }]);
