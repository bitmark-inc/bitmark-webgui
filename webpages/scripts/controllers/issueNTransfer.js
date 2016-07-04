// Copyright (c) 2014-2016 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

'use strict';

/**
 * @ngdoc function
 * @name bitmarkWebguiApp.controller:IssueNTransferCtrl
 * @description
 * # IssueNTransferCtrl
 * Controller of the bitmarkWebguiApp
 */
angular.module('bitmarkWebguiApp')
    .controller('IssueNTransferCtrl', function ($scope, $interval, $location, httpService, configuration, BitmarkCliConfig, BitmarkPayConfig) {
        if(configuration.getConfiguration().bitmarkCliConfigFile.length == 0){
            $location.path('/login');
        }

        if($location.path() == "/issue" ){
            $scope.showIssueView = true;
        }else if($location.path() == "/transfer"){
            $scope.showIssueView = false;
        }else {
            $location.path('/main');
        }

        var chain = configuration.getConfiguration().chain;
        var bitmarkCliConfigFile = BitmarkCliConfig[chain];
        var bitmarkPayConfigFile = BitmarkPayConfig[chain];;

        $scope.init = function(){
            localInit(chain);
        };


        var localInit = function(bitmarkChain){
            // get config file by chan type
            $scope.showWaiting = false;

            // default setup config
            $scope.bitmarkCliInfoSuccess = false;

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
                network:  chain,
                pay_config: bitmarkPayConfigFile,
                identity:"",
                asset:"",
                description:"",
                fingerprint:"",
                quantity:1,
                password:""
            };

            // transfer config
            $scope.transferConfig = {
                network:  chain,
                pay_config: bitmarkPayConfigFile,
                identity:"",
                txid:"",
                receiver:"",
                password:""
            };

            getInfo();
            checkBitmarkMode();
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
                        case "stopped":
                            $interval.cancel(infoPromise);
                            break;
                        }
                    }, function(statusErr){
                        $interval.cancel(infoPromise);
                        $scope.infoErr.msg = "bitmark-pay error: "+statusErr;
                        $scope.infoErr.show = true;
                        $scope.showWaiting = false;
                    });
            }, 3*1000);
        };
        var getInfo = function(){
            $scope.showWaiting = true;
            $scope.infoErr.show = false;
            $scope.infoAlert.show = false;

            httpService.send("onestepStatus",{
                network: chain,
                pay_config: bitmarkPayConfigFile
            }).then(function(infoResult){
                $interval.cancel(infoPromise);
                infoJobHash = infoResult.job_hash;
                $scope.bitmarkCliInfoSuccess = false;
                $scope.onestepStatusResult = infoResult;
                infoPromise = getBitmarkPayInfoInterval();

            }, function(infoErr){
                if( infoErr == "Failed to get bitmark-cli info") {
                    // bitmark-cli never setup, show setup view
                    $scope.showWaiting = false;
                    $scope.infoErr.show = true;
                    $scope.infoErr.msg = infoErr;
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


        $scope.bitmarkdAlert = {
            isNormal: false,
            msg: ""
        };
        var checkBitmarkMode = function(){
            httpService.send('getBitmarkdInfo').then(function(bitmarkInfo){
                if(bitmarkInfo.mode != 'Normal') {
                    // disable issue and transfer and show warning
                    $scope.bitmarkdAlert.msg = "Please wait until bitmarkd mode becomes Normal, please go to bitmark page to wait and check";
                } else {
                    $scope.bitmarkdAlert.isNormal = true;
                }
            }, function(bitmarkInfoErr){
                // bitmarkd is not running, disable issue and transfer and show warning
                $scope.bitmarkdAlert.msg = "bitmarkd is not running, please go to bitmark page to start and wait until the mode is Normal";
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

        $scope.submitGeneral = function (type) {
            $scope.showWaiting = true;

            var apiPath, configObj, resultObj;
            $scope.clearErrAlert(type);
            switch(type){
            case 'issue':
                $scope.issueResult = {
                    type:"danger",
                    msg: "",
                    failStart: null,
                    cliResult: null
                };
                apiPath = "onestepIssue";
                configObj = $scope.issueConfig;
                resultObj = $scope.issueResult;
                break;
            case 'transfer':
                $scope.transferResult = {
                    type:"danger",
                    msg: "",
                    cliResult: null
                };
                apiPath = "onestepTransfer";
                configObj = $scope.transferConfig;
                resultObj = $scope.transferResult;
                break;
            default:
                return;
            }

            submitGeneral(apiPath, configObj, resultObj);
        };

        var issuePromise;
        var issueWaitingTime = 10; // 10s
        var pollIssueCount = 0;
        var submitGeneral = function(api, apiParamObj, resultObj) {
            if($scope.onestepStatusResult.pay_result.available_balance < configuration.getConfiguration().mineFee) {
                resultObj.msg = "Your wallet doesn't have enough money (need at least > " + configuration.getConfiguration().mineFee + " satoshi), please send more bicoin to the address: "+ $scope.onestepStatusResult.pay_result.address;
                $scope.showWaiting = false;
                return;
            }

            apiParamObj.identity = $scope.onestepStatusResult.cli_result.identities[0].name;
            httpService.send(api, apiParamObj).then(
                function(result){
                    issuePromise = $interval(function(){
                        httpService.send("getBitmarkPayStatus", {
                            job_hash: result.job_hash
                        }).then(function(payStatus){
                            switch(payStatus){
                            case "success":
                                $interval.cancel(issuePromise);
                                httpService.send("getBitmarkPayResult", {
                                    "job_hash": result.job_hash
                                }).then(function(payResult){
                                    $scope.showWaiting = false;
                                    resultObj.type = "success";
                                    resultObj.msg = "Pay success!";
                                    resultObj.cliResult = result.cli_result;
                                }, function(payErr){
                                    $scope.showWaiting = false;
                                    resultObj.type = "danger";
                                    if(payErr.cli_result != null) {
                                        resultObj.msg = "Pay failed";
                                        resultObj.failStart = payErr.fail_start;
                                        resultObj.cliResult = payErr.cli_result;
                                    } else{
                                        resultObj.msg = payErr;
                                    }
                                });
                                break;
                            case "running":
                                pollIssueCount++;
                                if(pollIssueCount*3 > issueWaitingTime && !$scope.infoAlert.show){
                                    $scope.infoAlert.msg = "The bitmark-pay seems running for a long time, please check your bitcoin and bitmark-pay configuration. Would you like to stop the process?";
                                    $scope.showWaiting = false;
                                    $scope.infoAlert.show = true;
                                }
                                break;
                            default:
                                $interval.cancel(issuePromise);
                                $scope.showWaiting = false;
                                resultObj.type = "danger";
                                resultObj.msg = "Pay Error! Pay failed";
                                resultObj.cliResult = result.cli_result;
                                resultObj.failStart = 0;
                                break;
                            }
                        }, function(payStatusErr){
                            $interval.cancel(issuePromise);
                            $scope.showWaiting = false;
                            resultObj.type = "danger";
                            resultObj.msg = payStatusErr;
                        });
                    }, 3*1000);
                },
                function(errResult){
                    $scope.showWaiting = false;
                    resultObj.type = "danger";
                    if(errResult.cli_result != null) {
                        resultObj.msg = "Pay failed";
                        resultObj.failStart = errResult.fail_start;
                        resultObj.cliResult = errResult.cli_result;
                    } else{
                        resultObj.msg = errResult;
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

                            // after kill pay process, ask for info again
                            getInfo();
                        }
                    }, function(payStatusErr){
                        $interval.cancel(killPromise);
                        $scope.showWaiting = false;
                        $scope.infoErr.show = true;
                        $scope.infoErr.msg = payStatusErr;
                    });
                }, 3*1000);
            }, function(err){
                alertObj.show = true;
                alertObj.msg = err;
                $scope.showWaiting = false;
            });
        };

        $scope.killPayProcess = function(kill){
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
        };

        $scope.$on("$destroy", function(){
            $interval.cancel(infoPromise);
            $interval.cancel(issuePromise);
            $interval.cancel(killPromise);
        });


  });
