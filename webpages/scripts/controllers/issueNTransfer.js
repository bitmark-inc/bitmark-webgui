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
    .controller('IssueNTransferCtrl', ['$scope', '$timeout', 'httpService', "BitmarkCliConfig", "BitmarkPayConfig", "BitmarkChain", function ($scope, $timeout, httpService, BitmarkCliConfig, BitmarkPayConfig, BitmarkChain) {
        $scope.showSetup = false;
        $scope.bitmarkChain = BitmarkChain;

        // get config file by chan type
        // var bitmarkCliConfigFile = BitmarkCliConfig[BitmarkChain];
        // var bitmarkPayConfigFile = BitmarkPayConfig[BitmarkChain];
        var bitmarkCliConfigFile = "/home/yuntai/testWebgui/config/bitmark-cli/bitmark-cli-testing.config";
        var bitmarkPayConfigFile = "/home/yuntai/testWebgui/config/bitmark-pay/bitmark-pay-TESTING.xml";

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

        getInfo();

        $scope.clearErrAlert = function(type) {
            switch(type) {
            case "status":
                $scope.bitmarkPayError = null;
            case "issue":
                $scope.issueResult = null;
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
        $scope.bitmarkCliIssueConfig = {
            Config: bitmarkCliConfigFile,
            Identity:"",
            Password:"",
            Asset:"",
            Description:"",
            Fingerprint:"",
            Quantity:1
        };

        $scope.bitmarkPayIssueConfig = {
            Config: bitmarkPayConfigFile,
	    Password: "",
	    Net:  BitmarkChain,
            Txid: "",
	    Address: ""
        };

        $scope.issueResult = {
            type:"danger",
            msg: "",
            cliResult: null
        };



        // $scope.clearIssueResult = function() {
        //     $scope.issueResult = {
        //         type:"danger",
        //         msg: "",
        //         cliResult: null
        //     };
        // };

        $scope.submitIssue = function(){

            $scope.clearIssueResult();
            httpService.send("issueBitmark", $scope.bitmarkCliIssueConfig).then(
                function(cliResult){
                    // issue success, pay the tx
                    // check bitmarkPay net config
                    if($scope.bitmarkPayIssueConfig.Net == "local" ){
                        $scope.bitmarkPayIssueConfig.Net = "local_bitcoin_reg";
                    }

                    // TODO: set pay config txid, address
                    httpService.send("payBitmark", $scope.bitmarkPayIssueConfig).then(function(payResult){
                        // pay sucess
                        $scope.issueResult.type = "success";
                        $scope.issueResult.msg = "Pay success!";
                        $scope.issueResult.cliResult = cliResult;

                        // clean bitmarkCliIssueConfig
                        $scope.bitmarkCliIssueConfig = {
                            Identity:"",
                            Password:"",
                            Asset:"",
                            Description:"",
                            Fingerprint:"",
                            Quantity:1
                        };
                    }, function(payResultErr){
                        $scope.issueResult.type = "danger";
                        $scope.issueResult.msg = "Issue bitmark success but Payment Error!";
                    });
                }, function(cliResultErr){
                    $scope.issueResult.type = "danger";
                    $scope.issueResult.msg = "Issue bitmark Error!";
                });
        };

        // transfer config
        $scope.bitmarkCliTransferConfig = {
            Config: bitmarkCliConfigFile,
            Identity: "",
            Password: "",
            Txid: "",
            Receiver: ""
        };

        $scope.bitmarkPayTransferConfig = {
	    Net: BitmarkChain,
	    Config: bitmarkPayConfigFile,
	    Password: "",
	    Txid: "",
	    Addresses: null
        };


  }]);
