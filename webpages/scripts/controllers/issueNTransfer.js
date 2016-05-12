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
    .controller('IssueNTransferCtrl', ['$scope', 'httpService', "BitmarkCliConfig", "BitmarkPayConfig", "BitmarkChain", function ($scope, httpService, BitmarkCliConfig, BitmarkPayConfig, BitmarkChain) {
        $scope.showSetup = false;
        $scope.bitmarkChain = BitmarkChain;

        // get config file by chan type
        // var bitmarkCliConfigFile = BitmarkCliConfig[BitmarkChain];
        // var bitmarkPayConfigFile = BitmarkPayConfig[BitmarkChain];
        var bitmarkCliConfigFile = "/home/yuntai/testWebgui/config/bitmark-cli/bitmark-cli-TESTING.conf";
        var bitmarkPayConfigFile = "/home/yuntai/testWebgui/config/bitmark-pay/bitmark-pay-TESTING.xml";

        var getInfo = function(){
            httpService.send("getBitmarkCliInfo",{
                config:bitmarkCliConfigFile
            }).then(
                function(cliResult){
                    // get bitmarkPay info
                    httpService.send("getBitmarkPayInfo", {
                    net:BitmarkChain,
                        config: bitmarkPayConfigFile}).then(
                            function(payResult){
                                $scope.bitmarCliConfig = cliResult;
                                $scope.bitmarkPayConfig = payResult;
                                $scope.showSetup = false;
                            }, function(result){
                                $scope.bitmarkPayError = result;
                            });
                },function(result){
                    // need to setup bitmark-cli
                    $scope.showSetup = true;
                }
            );
        };

        getInfo();

        // default setup config
        $scope.bitmarkCliSetupConfig = {
            Config: bitmarkCliConfigFile,
	    Identity: "",
	    Password: "",
	    Network:  BitmarkChain,
	    Connect: "",
	    Description: ""
        };
        $scope.bitmarkPayEncryptConfig = {
            Config: bitmarkPayConfigFile,
	    Password: "",
	    Net:  BitmarkChain
        };

        $scope.submitSetup = function(){
            // check bitmarkPay net config
            if($scope.bitmarkPayEncryptConfig.Net == "local"){
                $scope.bitmarkPayEncryptConfig.Net = "local_bitcoin_reg";
            }

            $scope.setupError = '';
            httpService.send("setupBitmarkCli", $scope.bitmarkCliSetupConfig).then(
                function(cliResult){
                    // setupBitmarkCli success
                    httpService.send("setupBitmarkPay",$scope.bitmarkPayEncryptConfig).
                        then(function(payResult){
                            // setupBitmarkPay success
                            getInfo();
                        }, function(payResult){
                            // setupBitmarkPay fail
                            $scope.setupError = payResult;
                        });
                }, function(cliResult){
                    // setupBitmarkCli fail
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

        $scope.clearIssueResult = function() {
            $scope.issueResult = {
                type:"danger",
                msg: "",
                cliResult: null
            };
        };

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
