// Copyright (c) 2014-2016 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

'use strict';

/**
 * @ngdoc function
 * @name bitmarkMgmtApp.controller:EditCtrl
 * @description
 * # EditCtrl
 * Controller of the bitmarkMgmtApp
 */
angular.module('bitmarkMgmtApp')
    .controller('EditCtrl', ['$scope', '$location', 'httpService', function ($scope, $location, httpService) {
        // Check bitamrkd is not running, if it is running, stop it first
        httpService.send('statusBitmarkd').then(
            function(result){
                if(result.search("stop") >= 0) {
                    getAndSetBitmarkConfig();
                }else{
                    httpService.send("stopBitmarkd").then(
                        function(result){
                            $scope.errorMsg = "Bitmarkd has been stopped.";
                            getAndSetBitmarkConfig();
                        }, function(errorMsg){
                            $scope.errorMsg = errorMsg;
                        });
                }
            }, function(errorMsg){
                $scope.errorMsg = errorMsg;
            }
        );

        $scope.errorMsg = "";
        $scope.bitmarkProxy = false;

        $scope.verifyPassowrd = "";
        $scope.bitcoinPasswordEqual = true;
        $scope.publicKeyPattern = /^(\w|\d|\.|\-|:|\+|=|\^|!|\/|\*|\?|&|<|>|\(|\)|\[|\]|\{|\}|@|%|\$|#)+$/;

        // check bitcoin password
        $scope.$watchGroup(['localProxyTemp.Password','verifyPassword'], function(){
            if($scope.bitmarkConfig != null && !passwordVerified($scope.localProxyTemp.Password, $scope.verifyPassword)){
                $scope.bitcoinPasswordEqual = false;
            }else{
                $scope.bitcoinPasswordEqual = true;
            }
        });

        $scope.deleteItem = function(list, index){
            list.splice(index, 1);
        };

        $scope.addItem = function(list, limit){
            if (list.length < limit) {
                list.splice(list.length, 0, "");
            }
        };

        $scope.saveConfig = function(){
            saveConfig(function(){
                $scope.goUrl('/');
            });
        };


        $scope.saveConfigAndStart = function(){
            // send config post api and start bitmark then return to main page
            saveConfig(function(){
                httpService.send("startBitmarkd").then(
                    function(result){
                        $scope.goUrl('/');
                    }, function(errorMsg){
                        $scope.errorMsg = errorMsg;
                    });
            });
        };

        $scope.goUrl = function(path){
              $location.path(path);
        };


        $scope.bitmarkTestNetProxyTemp = {
            Username: "No-need-username",
            Password: "No-need-password",
            URL: "Testnet proxy not implement yet",
            Fee: "0.0002",
            Address: ""
        };
        $scope.bitmarkProxyTemp = {
            Username: "No-need-username",
            Password: "No-need-password",
            URL: "Bitmark proxy not implement yet",
            Fee: "0.0002",
            Address: ""
        };
        $scope.localProxyTemp = {
            Username: "",
            Password: "",
            URL: "",
            Fee: "",
            Address: ""
        };

        $scope.setBitmarkProxy = function(chainType){
            if(chainType == null) {
                chainType = angular.copy($scope.bitmarkConfig.Chain);

            }
            switch(chainType){
            case 'local':
                $scope.bitmarkProxy = false;
                break;
            case 'testing':
                $scope.bitmarkProxy = true;
                break;
            case 'bitmark':
                $scope.bitmarkProxy = true;
                break;
            }
        };

        var saveConfig = function(callBackFunc){
            // set bitcoin proxy
            switch($scope.bitmarkConfig.Chain){
            case 'local':
                $scope.bitmarkConfig.Bitcoin = $scope.localProxyTemp;
                break;
            case 'testing':
                $scope.bitmarkConfig.Bitcoin = $scope.bitmarkTestNetProxyTemp;
                break;
            case 'bitmark':
                $scope.bitmarkConfig.Bitcoin = $scope.bitmarkProxyTemp;
                break;
            }

            var result = checkBitmarkConfig($scope.bitmarkConfig);
            if (result.error !== "") {
                $scope.errorMsg = result.error;
                return;
            }

            var bitmarkConfig = result.bitmarkConfig;
            httpService.send('updateBitmarkConfig', bitmarkConfig).then(
                function(result){
                    if (callBackFunc != undefined){
                        callBackFunc();
                    }
                }, function(errorMsg){
                    $scope.errorMsg = errorMsg;
                });
        };

        var getAndSetBitmarkConfig = function(){
            httpService.send('getBitmarkConfig').then(
                function(result){
                    $scope.bitmarkConfig = initBitmarkConfig(result);

                    // setup bitmark proxy
                    switch($scope.bitmarkConfig.Bitcoin.URL){
                    case $scope.bitmarkTestNetProxyTemp.URL:
                        $scope.bitmarkProxy = true;
                        angular.extend($scope.bitmarkTestNetProxyTemp, $scope.bitmarkConfig.Bitcoin);
                        break;
                    case $scope.bitmarkProxyTemp.URL:
                        $scope.bitmarkProxy = true;
                        angular.extend($scope.bitmarkProxyTemp, $scope.bitmarkConfig.Bitcoin);
                        break;
                    default:
                        $scope.bitmarkProxy = false;
                        angular.extend($scope.localProxyTemp, $scope.bitmarkConfig.Bitcoin);
                    }
                }, function(errorMsg){
                    $scope.errorMsg = errorMsg;
                });
        };

        var initBitmarkConfig = function (bitmarkConfig){
            // give empty array for null field
            var checkItems = ["ClientRPC", "Peering", "Mining"];
            var checkFields = ["Listen", "Announce", "Connect"];

            for(var i=0; i<checkItems.length; i++){
                var checkItem = checkItems[i];
                for(var j=0; j<checkFields.length; j++){
                    var checkField = checkFields[j];
                    if(bitmarkConfig[checkItem][checkField] !== undefined && bitmarkConfig[checkItem][checkField] == null){
                        bitmarkConfig[checkItem][checkField] = [];
                    }
                }
            }
            return bitmarkConfig;
        };

        var passwordVerified = function(password, verifyPassword){
            if(password != "" && password != verifyPassword){
              return false;
            }
            return true;
        };

        // return {bitmarkConfig:{}, error:""}
        var checkBitmarkConfig = function(bitmarkConfig){
            var result = {
                bitmarkConfig:{},
                error:""
            };

            // check bitcoin password
            if(!$scope.bitmarkProxy && !passwordVerified(bitmarkConfig.Bitcoin.Password, $scope.verifyPassword)){
                result.error = "ErrPasswordNotEqual";
                return result;
            }

            // check publicKey
            for(var i=0; i<bitmarkConfig.Peering.Connect.length; i++) {
                if($scope.bitmarkForm["peerConnectPublicKey"+i].$invalid){
                    result.error = "Bitmark Peer Connect PublicKey invalid: #"+ (i+1);
                    return result;
                }
            }

            // delete empty element
            var checkItems = ["ClientRPC", "Peering", "Mining"];
            var checkFields = ["Listen", "Announce", "Connect"];

            for(var i=0; i<checkItems.length; i++){
                var checkItem = checkItems[i];
                for(var j=0; j<checkFields.length; j++){
                    var checkField = checkFields[j];
                    if(bitmarkConfig[checkItem][checkField] != undefined){
                        var fields = bitmarkConfig[checkItem][checkField];
                        for(var k=fields.length-1; k>0 ;k--){
                            if (fields[k] == ""){
                                fields.splice(k, 1);
                            }
                        }
                    }
                }
            }


            result.bitmarkConfig = bitmarkConfig;
            return result;
        };
  }]);
