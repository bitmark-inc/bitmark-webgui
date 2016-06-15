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
    .controller('LoginCreateCtrl', ['$scope', '$interval', '$location', '$cookies', 'httpService', 'configuration', 'BitmarkPayConfig', function ($scope, $interval, $location, $cookies, httpService, configuration, BitmarkPayConfig) {
        $scope.panelConfig = {
            showPart: 1
        };

        $scope.generateConfig = {
            chain: "testing",
            running: false,
            msg: []
        };
        $scope.privateKey = "";
        // var bitmarkCliConfigFile = "/home/yuntai/testWebgui/config/bitmark-cli/bitmark-cli-local.config";
        var bitmarkPayConfigFile = "/home/yuntai/testWebgui/config/bitmark-pay/bitmark-pay-LOCAL.xml";

        var encryptPromise;
        var encryptWaitingTime = 60; // 60s
        var pollEncryptCount = 0;
        var encryptJobHash;
        // will create keypair and then encrypt bitmark wallet
        var runGenerate = function(chain){
            // create key pair
            httpService.send("generateBitmarkKeyPair").then(function(keyPair){
                $scope.privateKey = keyPair.private_key;
                // TODO: encrypt bitmark wallet
                var net = chain;
                if(chain == "local"){
                    net = "local_bitcoin_reg";
                }
                httpService.send("setupBitmarkPay", {
                    net: net,
                    config: bitmarkPayConfigFile,
                    // config: BitmarkPayConfig[chain],
                    password: $scope.privateKey
                }).then(function(encryptPayJobHash){
                    $interval.cancel(encryptPromise);
                    encryptJobHash = encryptPayJobHash;
                    encryptPromise = $interval(function(){
                        httpService.send("getBitmarkPayStatus", {
                            job_hash: encryptPayJobHash
                        }).then(function(payStatusResult){
                            switch(payStatusResult){
                            case "success":
                                $interval.cancel(encryptPromise);
                                pollEncryptCount = 0;
                                $scope.showWaiting = false;
                                $scope.generateConfig.running = false;
                                $scope.panelConfig.showPart = 2;
                            break;
                            case "running":
                                pollEncryptCount++;
                                if(pollEncryptCount*3 > encryptWaitingTime){
                                    $scope.encryptAlert.msg = "Bitmark-pay has been running for "
                                        +pollEncryptCount*3+
                                        " seconds, normally it could cost 7 mins, would you want to stop the process?";
                                    $scope.showWaiting = false;
                                    $scope.encryptAlert.show = true;
                                }
                                break;
                            case "fail":
                                $interval.cancel(encryptPromise);
                                $scope.encryptErr.msg = "bitmark-pay error: "+payStatusResult;
                                $scope.encryptErr.show = true;
                                $scope.showWaiting = false;
                                break;
                            }
                        });
                    }, 3*1000);
                }, function(ecryptErr){
                    $scope.generateConfig.msg.push("failed to encrypt bitmarPay: "+ecryptErr);
                });
            }, function(keyPairErr){
                $scope.generateConfig.msg.push("failed to generate bitmark keypair: "+keyPairErr);
                $scope.generateConfig.running = false;
            });
        };

        $scope.generate = function(){
            $scope.generateConfig.running = true;
            // save chain to bitmark
            configuration.setChain($scope.generateConfig.chain);

            // Check bitcoin if user choose loca chain
            $scope.generateConfig.msg.push("Checking bitcoind...");
            if($scope.generateConfig.chain == "local"){
                httpService.send("statusBitcoind").then(function(bitcoinStatus){
                    if(bitcoinStatus == "stopped"){ // start bitcoind for the user
                        $scope.generateConfig.msg.push("bitcoind is stopped, try to starting it");

                        httpService.send("startBitcoind").then(function(startSuccess){
                            $scope.generateConfig.msg.push("bitcoind is started");
                            runGenerate($scope.generateConfig.chain);
                        }, function(startErr){
                            $scope.generateConfig.msg.push("failed to start bitcoind: "+startErr);
                            $scope.generateConfig.running = false;
                        });
                    }else{
                        $scope.generateConfig.msg.push("bitcoind is started...");
                        runGenerate($scope.generateConfig.chain);
                    }
                });
            } else {
                runGenerate($scope.generateConfig.chain);
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

        $scope.$on("$destroy", function(){
            $interval.cancel(encryptPromise);
        });
  }]);
