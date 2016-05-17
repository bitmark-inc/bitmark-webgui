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
    .controller('IssueNTransferCtrl', ['$scope', '$timeout', 'httpService', "BitmarkCliConfig", "BitmarkPayConfig", function ($scope, $timeout, httpService, BitmarkCliConfig, BitmarkPayConfig) {
        // var bitmarkCliConfigFile = "/home/yuntai/testWebgui/config/bitmark-cli/bitmark-cli-testing.config";
        // var bitmarkPayConfigFile = "/home/yuntai/testWebgui/config/bitmark-pay/bitmark-pay-TESTING.xml";

        var BitmarkChain = "";
        var bitmarkCliConfigFile = "";
        var bitmarkPayConfigFile = "";

        $scope.init = function(){
            httpService.send('getBitmarkConfig').then(
                function(result){
                    BitmarkChain = result.Chain;
                    $scope.showSetup = false;
                    $scope.bitmarkChain = BitmarkChain;
                    // get config file by chan type
                    bitmarkCliConfigFile = BitmarkCliConfig[BitmarkChain];
                    bitmarkPayConfigFile = BitmarkPayConfig[BitmarkChain];
                    getInfo();
                }, function(errorMsg){
                });
        };


        var getInfo = function(){
            httpService.send("onestepStatus",{
                cli_config: bitmarkCliConfigFile,
                network: BitmarkChain,
                pay_config: bitmarkPayConfigFile
            }).then(function(result){
                $scope.onestepStatusResult = result;
                $scope.showSetup = false;
            }, function(result){
                if( result == "Failed to get bitmark-cli info") {
                    $scope.showSetup = true;
                } else {
                    $scope.bitmarkPayError = result;
                }

            });
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

        // default setup config
        $scope.bitmarkSetupConfig = {
            network:  BitmarkChain,
            cli_config: bitmarkCliConfigFile,
            pay_config: bitmarkPayConfigFile,
            connect: "",
	    identity: "",
            description: "",
	    cli_password: "",
            pay_password: ""
        };


        $scope.submitSetup = function(){
            $scope.setupError = '';
            httpService.send("onestepSetup", $scope.bitmarkSetupConfig).then(function(result){
                $scope.showSetup = false;
                //wait for 10 seconds to sync the bitcoin
                $timeout(function(){
                    getInfo();
                }, 10*1000);
            }, function(error){
                $scope.setupError = cliResult;
            });
        };

        // default issue config
        $scope.bitmarkIssueConfig = {
            network:  BitmarkChain,
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

        $scope.submitIssue = function(){
            $scope.clearErrAlert('issue');
            $scope.issueResult = {
                type:"danger",
                msg: "",
                failStart: null,
                cliResult: null
            };
            $scope.bitmarkIssueConfig.identity = $scope.onestepStatusResult.identities[0].name;
            httpService.send("onestepIssue", $scope.bitmarkIssueConfig).then(
                function(result){
                    $scope.issueResult.type = "success";
                    $scope.issueResult.msg = "Pay success!";
                    $scope.issueResult.cliResult = result;
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

        // transfer config
        $scope.bitmarkTransferConfig = {
            network:  BitmarkChain,
            cli_config: bitmarkCliConfigFile,
            pay_config: bitmarkPayConfigFile,
            identity:"",
            txid:"",
            cli_password:"",
            pay_password:""
        };

        $scope.submitTransfer = function(){
            $scope.clearErrAlert('transfer');
            $scope.transferResult = {
                type:"danger",
                msg: "",
                cliResult: null
            };
            $scope.bitmarkTransferConfig.identity = $scope.onestepStatusResult.identities[0].name;

            httpService.send("onestepTransfer", $scope.bitmarkTransferConfig).then(
                function(result){
                    $scope.transferResult.type = "success";
                    $scope.transferResult.msg = "Pay success!";
                    $scope.transferResult.cliResult = result;
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

  }]);
