// Copyright (c) 2014-2016 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

'use strict';

/**
 * @ngdoc function
 * @name bitmarkWebguiApp.controller:ChainCtrl
 * @description
 * # ChainCtrl
 * Controller of the bitmarkWebguiApp
 */
angular.module('bitmarkWebguiApp')
    .controller('ChainCtrl', function ($scope, $location, httpService, BitmarkdConfig, utils) {

        $scope.init = function(){
            // check bitmarkd status
            httpService.send("statusBitmarkd").then(function(result){
                if(result == "started") {
                    // go to main view
                    $scope.$emit('Authenticated', true);
                    $location.path("/main");
                }
            });
        };

        $scope.request = {
            chain: "testing",
            running: false
        };

        $scope.error = {
            show: false,
            msg: ""
        };

        $scope.setErrorMsg = function(show, msg) {
            utils.setErrorMsg($scope.error, show, msg);
        };

        $scope.startNode = function(){
            $scope.request.running = true;
            $scope.setErrorMsg(false, '');
            // setup bitmarkd config
            httpService.send('setupBitmarkd', {
                config_file: BitmarkdConfig[$scope.request.chain]
            }).then(function(setupBitmarkdResult){
                // start bitmarkd
                httpService.send("startBitmarkd").then(
                    function(result){
                        $location.path("/main");
                        $scope.$emit('Authenticated', true);
                    }, function(errorMsg){
                        $scope.setErrorMsg(true, errorMsg);
                        $scope.request.running = false;
                    });
            }, function(setupBitmarkdErr){
                $scope.setErrorMsg(true, setupBitmarkdErr);
                $scope.request.running = false;
            });
        };
  });
