// Copyright (c) 2014-2016 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

'use strict';

/**
 * @ngdoc function
 * @name bitmarkWebguiApp.controller:LoginProcessCtrl
 * @description
 * # LoginCtrl
 * Controller of the bitmarkWebguiApp
 */
angular.module('bitmarkWebguiApp')
    .controller('LoginProcessCtrl', ['$scope', '$interval', '$location', '$cookies', 'httpService', 'configuration', 'BitmarkPayConfig', 'BitmarkCliSetupConfig', function ($scope, $interval, $location, $cookies, httpService, configuration, BitmarkPayConfig, BitmarkCliSetupConfig) {
        // var bitmarkCliConfigFile = "/home/yuntai/testWebgui/config/bitmark-cli/bitmark-cli-local.config";
        // var bitmarkPayConfigFile = "/home/yuntai/testWebgui/config/bitmark-pay/bitmark-pay-LOCAL.xml";

        $scope.panelConfig = {
            showPart: 1
        };

        $scope.generateConfig = {
            chain: "testing",
            running: false,
            msg: [],
            error: {
                msg: "",
                show: false
            },
            encryptAlert: {
                msg: "",
                show: false
            }
        };
        $scope.privateKey = "";

        var encryptPromise;
        var encryptWaitingTime = 60; // 60s
        var pollEncryptCount = 0;
        var encryptJobHash;

        var encryptWallet = function(chain, finalCallback){
            $scope.generateConfig.msg.push("encrypting wallet...");
            var net = chain;
            if(chain == "local"){
                net = "local_bitcoin_reg";
            }
            httpService.send("setupBitmarkPay", {
                net: net,
                // config: bitmarkPayConfigFile,
                config: BitmarkPayConfig[chain],
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
                            $scope.generateConfig.running = false;
                            $scope.panelConfig.showPart = 2;

                            if(null != finalCallback){
                                finalCallback();
                            }
                        case "running":
                            pollEncryptCount++;
                            if(pollEncryptCount*3 > encryptWaitingTime){
                                $scope.generateConfig.encryptAlert.msg = "Bitmark-pay has been running for "
                                    +pollEncryptCount*3+
                                    " seconds, normally it could cost 7 mins, would you want to stop the process?";
                                $scope.generateConfig.encryptAlert.show = true;
                            }
                            break;
                        case "fail":
                            $interval.cancel(encryptPromise);
                            $scope.generateConfig.error.show = true;
                            $scope.generateConfig.error.msg = "failed to encrypt wallet, please check your bitcoin status";
                            $scope.encryptErr.show = true;
                            break;
                        case "stopped":
                            $interval.cancel(encryptPromise);
                            $scope.generateConfig.error.show = true;
                            $scope.generateConfig.error.msg = "wallet was encrypted before, please decrypt your bitmark wallet first";
                            break;
                        }
                    });
                }, 3*1000);
            }, function(ecryptErr){
                $scope.generateConfig.error.show = true;
                $scope.generateConfig.error.msg = "failed to encrypt bitmarPay: "+ecryptErr;
            });
        };
        // will create keypair and then encrypt bitmark wallet
        var runGenerate = function(chain){
            // create key pair
            $scope.generateConfig.msg.push("generating bitmark keypair...");
            httpService.send("generateBitmarkKeyPair").then(function(keyPair){
                $scope.privateKey = keyPair.private_key;
                // encrypt bitmark wallet
                encryptWallet(chain, undefined);

            }, function(keyPairErr){
                $scope.generateConfig.error.show = true;
                $scope.generateConfig.error.msg = "failed to generate bitmark keypair: "+keyPairErr;
                $scope.generateConfig.running = false;
            });
        };

        $scope.generate = function(isCreateAccount){
            $scope.generateConfig.running = true;
            // save chain to bitmark
            configuration.setChain($scope.generateConfig.chain);

            // Check bitcoin if user choose loca chain
            $scope.generateConfig.msg.push("Checking bitcoind...");
            if($scope.generateConfig.chain == "local"){
                httpService.send("statusBitcoind").then(function(bitcoinStatus){
                    if(bitcoinStatus == "stopped"){ // start bitcoind for the user
                        $scope.generateConfig.msg.push("bitcoind is stopped, try to start it");

                        httpService.send("startBitcoind").then(function(startSuccess){
                            $scope.generateConfig.msg.push("bitcoind is started");
                            if(isCreateAccount){
                                runGenerate($scope.generateConfig.chain);
                            } else {
                                encryptWallet($scope.generateConfig.chain, $scope.done);
                            }
                        }, function(startErr){
                            $scope.generateConfig.error.show = true;
                            $scope.generateConfig.error.msg = "failed to start bitcoind: "+startErr;
                            $scope.generateConfig.running = false;
                        });
                    }else{
                        $scope.generateConfig.msg.push("bitcoind is started...");
                        if(isCreateAccount){
                            runGenerate($scope.generateConfig.chain);
                        } else {
                            encryptWallet($scope.generateConfig.chain, $scope.done);
                        }
                    }
                });
            } else {
                if(isCreateAccount){
                    runGenerate($scope.generateConfig.chain);
                } else {
                    encryptWallet($scope.generateConfig.chain, $scope.done);
                }
            }
        };

        $scope.killPayProcess = function(kill){
            if(kill){
                $interval.cancel(encryptPromise);
                pollEncryptCount = 0;
                if(encryptJobHash == "" || encryptJobHash == null) {
                    httpService.send('getBitmarkPayJob').then(function(jobHash){
                        encryptJobHash = jobHash;
                        killBitmarkPayStatusProcess(encryptJobHash, $scope.generateConfig.encryptAlert);
                    });
                }else{
                    killBitmarkPayStatusProcess(encryptJobHash, $scope.generateConfig.encryptAlert);
                }
            }else{
                $scope.generateConfig.encryptAlert.show = false;
                pollEncryptCount = 0;
            }
        };

        var killPromise;
        var killBitmarkPayStatusProcess = function(jobHash, alertObj){
            httpService.send('stopBitmarkPayProcess', {"job_hash": jobHash}).then(function(result){
                $interval.cancel(killPromise);
                killPromise = $interval(function(){
                    httpService.send("getBitmarkPayStatus", {
                        job_hash: jobHash
                    }).then(function(payStatus){
                        if(payStatus == "stopped"){
                            $interval.cancel(killPromise);
                            $scope.generateConfig.running = false;
                            alertObj.show = false;
                        }
                    });
                }, 3*1000);
            }, function(err){
                alertObj.show = true;
                alertObj.msg = err;
            });
        };

        $scope.toSetPassCode = function(){
                $scope.panelConfig.showPart = 3;
        };

        $scope.doneErr = {
            show: false,
            msg: ""
        };
        $scope.done = function(){
            //set up bitmark-cli
            if (!$scope.verifyPassword) {
                return;
            }
            var config = angular.copy(BitmarkCliSetupConfig);
            config.config = bitmarkCliConfigFile;
            config.password = $scope.password;
            config.network = $scope.generateConfig.chain;
            config.private_key =  $scope.privateKey;

            httpService.send('setupBitmarkCli', config).then(function(setupCliResult){
                $scope.$emit('Authenticated', true);
                $cookies.put('bitmark-chain', $scope.generateConfig.chain, {secure: true});
                $location.path("/main");
            }, function(setupCliErr){
                $scope.doneErr.msg = setupCliErr;
                $scope.doneErr.show = true;
            });
        };

        $scope.$on("$destroy", function(){
            $interval.cancel(encryptPromise);
            $interval.cancel(killPromise);
        });
  }]);
