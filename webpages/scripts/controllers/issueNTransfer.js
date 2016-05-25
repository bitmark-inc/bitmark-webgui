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

            // default setup config
            $scope.showWaiting = false;
            $scope.bitmarkCliInfoSuccess = false;
            $scope.payPasswordResult = true;
            $scope.cliPasswordResult = true;
            $scope.bitmarkSetupConfig = {
                network:  $scope.bitmarkChain,
                cli_config: bitmarkCliConfigFile,
                pay_config: bitmarkPayConfigFile,
                connect: "",
	        identity: "",
                description: "",
	        cli_password: "",
                pay_password: ""
            };

            // default issue config
            $scope.bitmarkIssueConfig = {
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
            $scope.bitmarkTransferConfig = {
                network:  $scope.bitmarkChain,
                cli_config: bitmarkCliConfigFile,
                pay_config: bitmarkPayConfigFile,
                identity:"",
                txid:"",
                cli_password:"",
                pay_password:""
            };

            $scope.bitmarkPayStatusAlert = {
                show: false,
                msg: ""
            };
            getInfo();
        };

        var infoPayPromise;
        var payJobHash = "";
        var waitingTime = 10; // 10s
        var pollBitmarkPayCount = 0;
        var getBitmarkPayInfoInterval = function(jobHash){
            return $interval(function(){
                pollBitmarkPayCount++;
                if(pollBitmarkPayCount*3 > waitingTime && !$scope.bitmarkPayStatusAlert.show){
                    $scope.bitmarkPayStatusAlert.msg = "The bitmark-pay seems running for a long time, please check your bitcoin and bitmark-pay configuration. Would you like to stop the process?";
                    $scope.showWaiting = false;
                    $scope.bitmarkPayStatusAlert.show = true;
                }
                httpService.send("getBitmarkPayStatus").then(
                    function(statusResult){
                        if(statusResult == "success"){
                            $interval.cancel(infoPayPromise);
                            infoPayPromise = null;
                            pollBitmarkPayCount = 0;
                            httpService.send("getBitmarkPayResult", {"job_hash":jobHash}).then(function(payResult){
                                $scope.onestepStatusResult.pay_result = payResult;
                                $scope.bitmarkCliInfoSuccess = true;
                            },function(payErr){
                                $scope.bitmarkPayError = payErr;
                                $interval.cancel(infoPayPromise);
                                infoPayPromise = null;
                            });
                        } else {
                            // TODO: see if statusResult is running and disable buttons
                        }
                    });
            }, 3*1000);
        };

        var getInfo = function(){
            httpService.send("onestepStatus",{
                cli_config: bitmarkCliConfigFile,
                network: $scope.bitmarkChain,
                pay_config: bitmarkPayConfigFile
            }).then(function(infoResult){
                // send onestep status success
                payJobHash = infoResult.job_hash;

                // using the job hash to query bitmark pay status, until the result is success then get the result
                $interval.cancel(infoPayPromise);
                infoPayPromise = getBitmarkPayInfoInterval(payJobHash);

                $scope.bitmarkCliInfoSuccess = false;
                $scope.onestepStatusResult = infoResult;
                $scope.showSetup = false;
            }, function(infoErr){
                if( infoErr == "Failed to get bitmark-cli info") {
                    // bitmark-cli never setup, show setup view
                    $scope.showSetup = true;
                } else {
                    httpService.send('getBitmarkPayJob').then(function(jobHash){
                        payJobHash = jobHash;
                        if(jobHash != "") {
                            // bitmark-pay error
                            infoPayPromise = getBitmarkPayInfoInterval(payJobHash);
                            $scope.bitmarkPayStatusAlert.msg = "The previous bitmark-pay is running. Would you like to stop the process?";
                            $scope.showWaiting = false;
                            $scope.bitmarkPayStatusAlert.show = true;
                        } else {
                            $scope.bitmarkPayError = infoErr;
                        }
                    });
                }
            });
        };

        var killPromise;
        var killBitmarkPayStatusProcess = function(jobHash){
            $scope.showWaiting = true;
            httpService.send('stopBitmarkPayProcess', {"job_hash": jobHash}).then(function(result){
                $interval.cancel(killPromise);
                killPromise = $interval(function(){
                    httpService.send("getBitmarkPayStatus").then(function(payStatus){
                        if(payStatus == "stopped"){
                            $interval.cancel(killPromise);
                            $scope.showWaiting = false;
                            $scope.bitmarkPayStatusAlert.show = false;
                        }else {
                            // TODO: disable buttons
                        }
                    });
                }, 3*1000);

            }, function(err){
                $scope.bitmarkPayStatusAlert.show = true;
                $scope.bitmarkPayStatusAlert.msg = err;
            });
        };

        $scope.killBitmarkPayStatusProcess = function(kill) {
            if(kill){
                $interval.cancel(infoPayPromise);
                infoPayPromise = null;
                pollBitmarkPayCount = 0;

                // TODO: kill bitmark-pay process .. hide the alert
                if(payJobHash == "" || payJobHash == null) {
                    httpService.send('getBitmarkPayJob').then(function(jobHash){
                        payJobHash = jobHash;
                        killBitmarkPayStatusProcess(jobHash);
                    });
                }else{
                     killBitmarkPayStatusProcess(payJobHash);
                }

            }else{
                $scope.bitmarkPayStatusAlert.show = false;
                pollBitmarkPayCount = 0;
            }
        };

        $scope.clearErrAlert = function(type) {
            switch(type) {
            case "status":
                $scope.bitmarkPayError = null;
            case "issue":
                $scope.issueResult = null;
            case "transfer":
                $scope.transferResult = null;
            default:
            }
        };

        $scope.submitSetup = function(){
            $scope.setupError = '';

            httpService.send("onestepSetup", $scope.bitmarkSetupConfig).then(function(result){
                //wait for 10 seconds to sync the bitcoin
                $timeout(function(){
                    $scope.showSetup = false;
                    getInfo();
                }, 10*1000);
            }, function(error){
                $scope.setupError = cliResult;
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
            $scope.bitmarkIssueConfig.identity = $scope.onestepStatusResult.cli_result.identities[0].name;
            httpService.send("onestepIssue", $scope.bitmarkIssueConfig).then(
                function(result){
                    issuePromise = $interval(function(){
                        httpService.send("getBitmarkPayStatus").then(function(payStatus){
                            if(payStatus == "success"){
                                $interval.cancel(issuePromise);
                                httpService.send("getBitmarkPayResult", {"job_hash": result.job_hash}).then(function(payResult){
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
            $scope.bitmarkTransferConfig.identity = $scope.onestepStatusResult.cli_result.identities[0].name;

            httpService.send("onestepTransfer", $scope.bitmarkTransferConfig).then(
                function(result){
                    transferPromise = $interval(function(){
                        httpService.send("getBitmarkPayStatus").then(function(payStatus){
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

        $scope.$on("$destroy", function(){
            $interval.cancel(infoPayPromise);
            $interval.cancel(issuePromise);
            $interval.cancel(transferPromise);
            $interval.cancel(killPromise);
        });

  }]);
