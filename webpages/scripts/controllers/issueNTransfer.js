// Copyright (c) 2014-2016 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

'use strict';

/**
 * @ngdoc function
 * @name bitmarkWebguiApp.controller:MainCtrl
 * @description
 * # MainCtrl
 * Controller of the bitmarkWebguiApp
 */
angular.module('bitmarkWebguiApp')
    .controller('IssueNTransferCtrl', ['$scope', '$timeout', '$interval', 'httpService', "BitmarkCliConfig", "BitmarkPayConfig", function ($scope, $timeout, $interval, httpService, BitmarkCliConfig, BitmarkPayConfig) {
        // var bitmarkCliConfigFile = "/home/yuntai/testWebgui/config/bitmark-cli/bitmark-cli-local.config";
        // var bitmarkPayConfigFile = "/home/yuntai/testWebgui/config/bitmark-pay/bitmark-pay-LOCAL.xml";

        var bitmarkCliConfigFile = "";
        var bitmarkPayConfigFile = "";

        $scope.init = function(){
            httpService.send('getBitmarkConfig').then(
                function(result){
                    localInit(result.Chain);
                }, function(errorMsg){
                });
        };


        var localInit = function(bitmarkChain){
            $scope.showSetup = false;
            $scope.bitmarkChain = bitmarkChain;
            // get config file by chan type
            bitmarkCliConfigFile = BitmarkCliConfig[$scope.bitmarkChain];
            bitmarkPayConfigFile = BitmarkPayConfig[$scope.bitmarkChain];

            $scope.showWaiting = false;

            // default setup config
            $scope.bitmarkCliInfoSuccess = false;
            $scope.payPasswordResult = true;
            $scope.cliPasswordResult = true;
            $scope.setupErr = {
                show: false,
                msg: ""
            };
            $scope.setupAlert = {
                show: false,
                msg: ""
            };
            $scope.setupConfig = {
                network:  $scope.bitmarkChain,
                cli_config: bitmarkCliConfigFile,
                pay_config: bitmarkPayConfigFile,
                connect: "",
	        identity: "",
                description: "",
	        cli_password: "",
                pay_password: ""
            };


            // default info config
            $scope.infoErr = {
                show: false,
                msg: ""
            };
            $scope.infoAlert = {
                show: false,
                msg: ""
            };

            // default issue config
            $scope.issueConfig = {
                network:  $scope.bitmarkChain,
                cli_config: bitmarkCliConfigFile,
                pay_config: bitmarkPayConfigFile,
                identity:"",
                asset:"",
                description:"",
                fingerprint:"",
                quantity:1,
                cli_password:"",
                pay_password:""
            };

            // transfer config
            $scope.transferConfig = {
                network:  $scope.bitmarkChain,
                cli_config: bitmarkCliConfigFile,
                pay_config: bitmarkPayConfigFile,
                identity:"",
                txid:"",
                cli_password:"",
                pay_password:""
            };

            getInfo();
        };

        var infoPromise;
        var infoJobHash = "";
        var infoWaitingTime = 10; // 10s
        var pollInfoCount = 0;
        var getBitmarkPayInfoInterval = function(){
            return $interval(function(){
                httpService.send("getBitmarkPayStatus", {
                    job_hash: infoJobHash
                }).then(
                    function(statusResult){
                        switch(statusResult){
                        case "success":
                            $interval.cancel(infoPromise);
                            pollInfoCount = 0;
                            $scope.showWaiting = false;
                            httpService.send("getBitmarkPayResult", {"job_hash":infoJobHash}).then(function(payResult){
                                $scope.onestepStatusResult.pay_result = payResult;
                                $scope.bitmarkCliInfoSuccess = true;
                            },function(payErr){
                                $scope.infoErr.msg = payErr;
                                $scope.infoErr.show = true;
                            });
                            break;
                        case "running":
                            pollInfoCount++;
                            if(pollInfoCount*3 > infoWaitingTime && !$scope.infoAlert.show){
                                $scope.infoAlert.msg = "The bitmark-pay seems running for a long time, please check your bitcoin and bitmark-pay configuration. Would you like to stop the process?";
                                $scope.showWaiting = false;
                                $scope.infoAlert.show = true;
                            }
                            break;
                        case "fail":
                            $interval.cancel(infoPromise);
                            $scope.infoErr.msg = "bitmark-pay error: "+statusResult;
                            $scope.infoErr.show = true;
                            $scope.showWaiting = false;
                            break;
                        }
                    });
            }, 3*1000);
        };
        var getInfo = function(){
            $scope.showWaiting = true;
            $scope.infoErr.show = false;
            $scope.infoAlert.show = false;

            httpService.send("onestepStatus",{
                cli_config: bitmarkCliConfigFile,
                network: $scope.bitmarkChain,
                pay_config: bitmarkPayConfigFile
            }).then(function(infoResult){
                $interval.cancel(infoPromise);
                infoJobHash = infoResult.job_hash;
                $scope.bitmarkCliInfoSuccess = false;
                $scope.onestepStatusResult = infoResult;
                $scope.showSetup = false;
                infoPromise = getBitmarkPayInfoInterval();

            }, function(infoErr){
                if( infoErr == "Failed to get bitmark-cli info") {
                    // bitmark-cli never setup, show setup view
                    $scope.showWaiting = false;
                    $scope.showSetup = true;
                } else {
                    httpService.send('getBitmarkPayJob').then(function(jobHash){
                        infoJobHash = jobHash;
                        $scope.showWaiting = false;
                        if(jobHash != "") {
                            // bitmark-pay error
                            infoPromise = getBitmarkPayInfoInterval();
                            $scope.infoAlert.msg = "The previous bitmark-pay is running. Would you like to stop the process?";
                            $scope.infoAlert.show = true;
                        } else {
                            $scope.infoErr.msg = infoErr;
                            $scope.infoErr.show = true;
                        }
                    });
                }
            });
        };

        $scope.clearErrAlert = function(type) {
            switch(type) {
            case "issue":
                $scope.issueResult = null;
            case "transfer":
                $scope.transferResult = null;
            }
        };

        var setupPromise;
        var setupWaitingTime = 60; // 60s
        var pollSetupCount = 0;
        var setupJobHash;
        $scope.submitSetup = function(){
            $scope.showWaiting = true;
            $scope.setupErr.show = false;
            $scope.setupAlert.show = false;
            var net = $scope.setupConfig.network;
            if(net == "local") {
                net = "local_bitcoin_reg";
            }

            // setup bitmark-pay first, setup bitmark-cli while success
            httpService.send("setupBitmarkPay", {
                net: net,
                config: $scope.setupConfig.pay_config,
                password: $scope.setupConfig.pay_password
            }).then(function(setupPayJobHash){
                $interval.cancel(setupPromise);
                setupJobHash = setupPayJobHash;
                setupPromise = $interval(function(){
                    httpService.send("getBitmarkPayStatus", {
                        job_hash: setupPayJobHash
                    }).then(function(payStatusResult){
                        switch(payStatusResult){
                        case "success":
                            // do bitmark-cli setup
                            $interval.cancel(setupPromise);
                            pollSetupCount = 0;
                            $scope.showWaiting = false;
                            httpService.send('setupBitmarkCli', {
                                config: $scope.setupConfig.cli_config,
                                identity: $scope.setupConfig.identity,
                                password: $scope.setupConfig.cli_password,
                                network: $scope.setupConfig.network,
                                connect: $scope.setupConfig.connect,
                                description: $scope.setupConfig.description
                            }).then(function(setupCliResult){
                                $scope.showSetup = false;
                                getInfo();
                            }, function(setupCliErr){
                                $scope.setupErr.msg = setupCliErr;
                                $scope.setupErr.show = true;
                            });
                            break;
                        case "running":
                            pollSetupCount++;
                            if(pollSetupCount*3 > setupWaitingTime){

                                $scope.setupAlert.msg = "Bitmark-pay has been running for "
                                    +pollSetupCount*3+
                                    " seconds, normally it could cost 7 mins, would you want to stop the process?";
                                $scope.showWaiting = false;
                                $scope.setupAlert.show = true;
                            }
                            break;
                        case "fail":
                            $interval.cancel(setupPromise);
                            $scope.setupErr.msg = "bitmark-pay error: "+payStatusResult;
                            $scope.setupErr.show = true;
                            $scope.showWaiting = false;
                            break;
                        }
                    });
                }, 3*1000);
            }, function(setupBitmarkPayErr){
                $scope.showWaiting = false;
                $scope.setupErr.msg = setupBitmarkPayErr;
                $scope.setupErr.show = true;
            });
        };

        var issuePromise;
        $scope.submitIssue = function(){
            $scope.clearErrAlert('issue');
            $scope.issueResult = {
                type:"danger",
                msg: "",
                failStart: null,
                cliResult: null
            };
            $scope.issueConfig.identity = $scope.onestepStatusResult.cli_result.identities[0].name;
            httpService.send("onestepIssue", $scope.issueConfig).then(
                function(result){
                    issuePromise = $interval(function(){
                        httpService.send("getBitmarkPayStatus", {
                            job_hash: result.job_hash
                        }).then(function(payStatus){
                            if(payStatus == "success"){
                                $interval.cancel(issuePromise);
                                httpService.send("getBitmarkPayResult", {
                                    "job_hash": result.job_hash
                                }).then(function(payResult){
                                    $scope.issueResult.type = "success";
                                    $scope.issueResult.msg = "Pay success!";
                                    $scope.issueResult.cliResult = result.cli_result;
                                }, function(payErr){
                                    $scope.issueResult.type = "danger";
                                    if(payErr.cli_result != null) {
                                        $scope.issueResult.msg = "Pay failed";
                                        $scope.issueResult.failStart = payErr.fail_start;
                                        $scope.issueResult.cliResult = payErr.cli_result;
                                    } else{
                                        $scope.issueResult.msg = payErr;
                                    }
                                });
                            }else{
                            // TODO: see if bitmark-pay is still running
                            }
                        });
                    }, 3*1000);
                },
                function(errResult){
                    $scope.issueResult.type = "danger";
                    if(errResult.cli_result != null) {
                        $scope.issueResult.msg = "Pay failed";
                        $scope.issueResult.failStart = errResult.fail_start;
                        $scope.issueResult.cliResult = errResult.cli_result;
                    } else{
                        $scope.issueResult.msg = errResult;
                    }
                });
        };


        var transferPromise;
        $scope.submitTransfer = function(){
            $scope.clearErrAlert('transfer');
            $scope.transferResult = {
                type:"danger",
                msg: "",
                cliResult: null
            };
            $scope.transferConfig.identity = $scope.onestepStatusResult.cli_result.identities[0].name;

            httpService.send("onestepTransfer", $scope.transferConfig).then(
                function(result){
                    transferPromise = $interval(function(){
                        httpService.send("getBitmarkPayStatus", {
                            job_hash: result.job_hash
                        }).then(function(payStatus){
                            if(payStatus == "success"){
                                $interval.cancel(transferPromise);
                                httpService.send("getBitmarkPayResult", {"job_hash": result.job_hash}).then(function(payResult){
                                    $scope.transferResult.type = "success";
                                    $scope.transferResult.msg = "Pay success!";
                                    $scope.transferResult.cliResult = result.cli_result;
                                },function(payErr){
                                    // TODO: pay error
                                });
                            } else {
                                // TODO: see if bitmark-pay is still running
                            }
                        });
                    }, 3*1000);


                }, function(errResult){
                    $scope.transferResult.type = "danger";
                    if(errResult.cli_result != null) {
                        $scope.transferResult.msg = "Pay failed";
                        $scope.transferResult.cliResult = errResult.cli_result;
                    } else{
                        $scope.transferResult.msg = errResult;
                    }
                });
        };

        var killPromise;
        var killBitmarkPayStatusProcess = function(jobHash, alertObj){
            $scope.showWaiting = true;
            httpService.send('stopBitmarkPayProcess', {"job_hash": jobHash}).then(function(result){
                $interval.cancel(killPromise);
                killPromise = $interval(function(){
                    httpService.send("getBitmarkPayStatus", {
                        job_hash: jobHash
                    }).then(function(payStatus){
                        if(payStatus == "stopped"){
                            $interval.cancel(killPromise);
                            $scope.showWaiting = false;
                            alertObj.show = false;
                        }
                    });
                }, 3*1000);
            }, function(err){
                alertObj.show = true;
                alertObj.msg = err;
                $scope.showWaiting = false;
            });
        };

        $scope.killPayProcess = function(type, kill){
            switch(type){
            case "setup":
                $interval.cancel(setupPromise);
                pollSetupCount = 0;
                killBitmarkPayStatusProcess(setupJobHash, $scope.setupAlert);
                break;
            case "info":
                if(kill){
                    $interval.cancel(infoPromise);
                    pollInfoCount = 0;
                    if(infoJobHash == "" || infoJobHash == null) {
                        httpService.send('getBitmarkPayJob').then(function(jobHash){
                            infoJobHash = jobHash;
                            killBitmarkPayStatusProcess(infoJobHash, $scope.infoAlert);
                        });
                    }else{
                        killBitmarkPayStatusProcess(infoJobHash, $scope.infoAlert);
                    }
                }else{
                    $scope.infoAlert.show = false;
                    $scope.showWaiting = true;
                    pollInfoCount = 0;
                }
                break;
            }
        };

        $scope.$on("$destroy", function(){
            $interval.cancel(setupPromise);
            $interval.cancel(infoPromise);
            $interval.cancel(issuePromise);
            $interval.cancel(transferPromise);
            $interval.cancel(killPromise);
        });

  }]);
