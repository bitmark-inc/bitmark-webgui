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
    .controller('LoginProcessCtrl', function ($scope, $q, $interval, $location, httpService, configuration, BitmarkPayConfig, BitmarkCliConfig, BitmarkCliSetupConfig, BitmarkdConfig) {
        if(configuration.getConfiguration().bitmarkCliConfigFile.length != 0){
            $location.path('/login');
        }
        $scope.bitmarkd = {
            isRunning: false
        };

        // checkout bitmarkd status first
        httpService.send('statusBitmarkd').then(
            function(result){
                if(result.search("stop") >= 0) {
                    $scope.bitmarkd.isRunning = false;
                }else{
                    httpService.send('getBitmarkdInfo').then(function(info){
                        $scope.bitmarkd.isRunning = true;
                        $scope.generateConfig.chain = info.chain;
                        configuration.setChain(info.chain);
                    }, function(infoErr){
                    });
                }
            }, function(errorMsg){
            });


        $scope.stopBitmarkd = function(){
            httpService.send("stopBitmarkd").then(
                function(result){
                    $scope.bitmarkd.isRunning = false;
                }, function(errorMsg){
                });
        };

        $scope.panelConfig = {
            showPart: 1
        };

        $scope.generateConfig = {
            chain: configuration.getConfiguration().chain,
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

        // will create keypair and then encrypt bitmark wallet
        var runGenerate = function(chain){
            // create key pair
            $scope.generateConfig.msg.push("generating bitmark keypair...");
            httpService.send("generateBitmarkKeyPair").then(function(keyPair){
                $scope.privateKey = keyPair.private_key;
                $scope.generateConfig.running = false;
                $scope.panelConfig.showPart = 2;
                // encrypt bitmark wallet
                // encryptWallet(chain, undefined);

            }, function(keyPairErr){
                $scope.generateConfig.error.show = true;
                $scope.generateConfig.error.msg = "failed to generate bitmark keypair: "+keyPairErr;
                $scope.generateConfig.running = false;
            });
        };

        $scope.generate = function(){
            $scope.generateConfig.running = true;
            // save chain to bitmark
            configuration.setChain($scope.generateConfig.chain);

            // Check bitcoin if user choose loca chain
            $scope.generateConfig.msg.push("Checking bitcoind...");
            runGenerate($scope.generateConfig.chain);
        };

        $scope.killPayProcess = function(kill){
            if(kill){
                $scope.generateConfig.encryptAlert.show = false;
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


        $scope.disableDoneBtn = false;
        $scope.doneErr = {
            show: false,
            msg: ""
        };

        var restorePromise;
        var restoreWaitingTime = 60; // 60s
        var pollRestoreCount = 0;
        var restoreJobHash;

        var restoreWallet = function(chain){
            var restoreFinish = $q.defer();
            var net = chain;

            httpService.send('restoreBitmarkPay', {
                net: net,
                config: BitmarkPayConfig[$scope.generateConfig.chain],
                seed: $scope.privateKey
            }).then(function(restoreJobHash){
                $interval.cancel(restorePromise);
                restoreJobHash = restoreJobHash;
                restorePromise = $interval(function(){
                    httpService.send("getBitmarkPayStatus", {
                        job_hash: restoreJobHash
                    }).then(function(payStatusResult){
                        switch(payStatusResult){
                        case "success":
                            $interval.cancel(restorePromise);
                            pollRestoreCount = 0;
                            restoreFinish.resolve();
                            break;
                        case "running":
                            pollRestoreCount++;
                            if(pollRestoreCount*3 > restoreWaitingTime){
                                $scope.generateConfig.encryptAlert.msg = "Bitmark-pay has been running for "
                                    +pollRestoreCount*3+
                                    " seconds, normally it could cost 7 mins, would you want to stop the process?";
                                $scope.generateConfig.encryptAlert.show = true;
                            }
                            break;
                        case "fail":
                            $interval.cancel(restorePromise);
                            $scope.doneErr.show = true;
                            $scope.doneErr.msg = "failed to restore wallet";
                            $scope.diableDoneBtn = false;
                            restoreFinish.reject();
                            break;
                        case "stopped":
                            $interval.cancel(restorePromise);
                            $scope.doneErr.show = true;
                            $scope.doneErr.msg = "restore stopped";
                            $scope.diableDoneBtn = false;
                            restoreFinish.reject();
                            break;
                        }
                    }, function(payStatusError){
                        $interval.cancel(restorePromise);
                        $scope.doneErr.show = true;
                        $scope.doneErr.msg = payStatusError;
                        $scope.diableDoneBtn = false;
                        restoreFinish.reject();
                    });
                }, 3*1000);
            }, function(restoreErr){
                $scope.doneErr.show = true;
                $scope.doneErr.msg = "failed to restore bitmarPay: "+ restoreErr;
                $scope.diableDoneBtn = false;
                restoreFinish.reject();
            });

            return restoreFinish.promise;
        };

        var encryptPromise;
        var encryptWaitingTime = 60; // 60s
        var pollEncryptCount = 0;
        var encryptJobHash;

        var encryptWallet = function(chain){
            var encryptFinish = $q.defer();

            $scope.generateConfig.msg.push("encrypting wallet...");
            var net = chain;

            httpService.send("setupBitmarkPay", {
                net: net,
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
                            // $scope.generateConfig.running = false;
                            encryptFinish.resolve();
                            break;
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
                            $scope.doneErr.show = true;
                            $scope.doneErr.msg = "failed to encrypt wallet, please check your bitcoin status";
                            $scope.diableDoneBtn = false;
                            encryptFinish.reject();
                            break;
                        case "stopped":
                            $interval.cancel(encryptPromise);
                            $scope.doneErr.show = true;
                            $scope.doneErr.msg = "wallet was encrypted before, please decrypt your bitmark wallet first";
                            $scope.diableDoneBtn = false;
                            encryptFinish.reject();
                            break;
                        }
                    }, function(payStatusError){
                        $interval.cancel(encryptPromise);
                        $scope.doneErr.show = true;
                        $scope.doneErr.msg = payStatusError;
                        $scope.diableDoneBtn = false;
                        encryptFinish.reject();
                    });
                }, 3*1000);
            }, function(ecryptErr){
                $scope.doneErr.show = true;
                $scope.doneErr.msg = "failed to encrypt bitmarPay: "+ecryptErr;
                encryptFinish.reject();
            });
            return encryptFinish.promise;
        };


        var setupBitmarkCli = function(){
            var config = angular.copy(BitmarkCliSetupConfig);
            config.config = BitmarkCliConfig[$scope.generateConfig.chain];
            config.password = $scope.password;
            config.network = $scope.generateConfig.chain;
            config.private_key =  $scope.privateKey;

            httpService.send('setupBitmarkCli', config).then(function(setupCliResult){
                // setup bitmarkConfig file in server
                if(!$scope.bitmarkd.isRunning){
                    httpService.send('setupBitmarkd', {
                        config_file: BitmarkdConfig[$scope.generateConfig.chain]
                    }).then(function(setupBitmarkdResult){
                        configuration.setBitmarkCliConfigFile(BitmarkCliConfig[$scope.generateConfig.chain]);
                        $scope.$emit('Authenticated', true);
                        $location.path("/main");
                    }, function(setupBitmarkdErr){
                        $scope.doneErr.msg = setupBitmarkdErr;
                        $scope.doneErr.show = true;
                        $scope.diableDoneBtn = false;
                    });
                } else {
                    configuration.setBitmarkCliConfigFile(BitmarkCliConfig[$scope.generateConfig.chain]);
                    $scope.$emit('Authenticated', true);
                    $location.path("/main");
                }
            }, function(setupCliErr){
                $scope.doneErr.msg = setupCliErr;
                $scope.doneErr.show = true;
                $scope.diableDoneBtn = false;
            });
        };


        $scope.done = function(){
            //set up bitmark-cli
            if (!$scope.verifyPassword) {
                return;
            }

            $scope.diableDoneBtn = true;
            // restore encrypt bitmarkPay wallet and setup bitmark-cli
            restoreWallet($scope.generateConfig.chain).then(function(){
                encryptWallet($scope.generateConfig.chain).then(function(){
                    setupBitmarkCli();
                });
            });
        };

        $scope.$on("$destroy", function(){
            $interval.cancel(encryptPromise);
            $interval.cancel(restorePromise);
            $interval.cancel(killPromise);
        });
  });
