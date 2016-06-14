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
    .controller('LoginCreateCtrl', ['$scope', '$timeout', '$location', '$cookies', 'httpService', 'configuration', function ($scope, $timeout, $location, $cookies, httpService, configuration) {
        $scope.panelConfig = {
            showPart: 1
        };

        $scope.generateConfig = {
            chain: "testing",
            running: false,
            msg: []
        };
        $scope.generate = function(){
            $scope.generateConfig.running = true;
            // save chain to bitmark
            configuration.setChain($scope.generateConfig.chain);

            // TODO: Check bitcoin if user choose loca chain
            $scope.generateConfig.msg.push("Checking bitcoind...");
            if($scope.generateConfig.chain == "local"){
                httpService.send("statusBitcoind").then(function(bitcoinStatus){
                    if(bitcoinStatus == "stopped"){ // start bitcoind for the user
                        $scope.generateConfig.msg.push("bitcoind is stopped, try to starting it");
                        httpService.send("startBitcoind").then(function(startSuccess){
                            $scope.generateConfig.msg.push("bitcoind is started");

                            // create key pair
                            // encrypt bitmark wallet
                        }, function(startFail){
                            $scope.generateConfig.msg.push("fail to start bitcoind: "+startFail);
                        });
                    }
                });
            }



            // $timeout(function(){
            //     // mock when generate is done
            //     $scope.generateConfig.running = false;
            //     $scope.panelConfig.showPart = 2;
            // }, 2*1000);

        };


        $scope.setPassCode = function(){
            $scope.panelConfig.showPart = 3;
        };

        $scope.done = function(){
            $location.path("/main");
        };

  }]);
