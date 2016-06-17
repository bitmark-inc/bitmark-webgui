// Copyright (c) 2014-2016 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

'use strict';

/**
 * @ngdoc function
 * @name bitmarkWebguiApp.controller:EditCtrl
 * @description
 * # EditCtrl
 * Controller of the bitmarkWebguiApp
 */
angular.module('bitmarkWebguiApp')
    .controller('EditCtrl', function ($scope, $location, httpService, BitmarkProxyURL, ProxyTemp, BitmarkCliSetupConfig, configuration) {
        if(configuration.getConfiguration().bitmarkCliConfigFile.length == 0){
            $location.path('/login');
        }

        $scope.BitmarkCliSetupConfig = BitmarkCliSetupConfig;

        $scope.error = {
            show: false,
            msg: ""
        };

        // Check bitamrkd is not running, if it is running, stop it first
        httpService.send('statusBitmarkd').then(
            function(result){
                if(result.search("stop") >= 0) {
                    getAndSetBitmarkConfig();
                }else{
                    httpService.send("stopBitmarkd").then(
                        function(result){
                            $scope.error.show = true;
                            $scope.error.msg = "Bitmarkd has been stopped.";
                            getAndSetBitmarkConfig();
                        }, function(errorMsg){
                            $scope.error.show = true;
                            $scope.error.msg = errorMsg;
                        });
                }
            }, function(errorMsg){
                $scope.error.show = true;
                $scope.error.msg = errorMsg;
            }
        );

        // setup proxy temp
        var proxyType = {
            "local": "local",
            "other": "other",
            "testing": "testing",
            "bitmark": "bitmark"
        };

        $scope.bitcoinUseProxy = proxyType.local;

        $scope.otherProxyTemp = angular.copy(ProxyTemp);

        $scope.bitmarkTestNetProxyTemp = angular.copy(ProxyTemp);
        $scope.bitmarkTestNetProxyTemp.URL = BitmarkProxyURL.testing;

        $scope.bitmarkProxyTemp = angular.copy(ProxyTemp);
        $scope.bitmarkProxyTemp.URL = BitmarkProxyURL.bitmark;

        $scope.localBitcoin = {
            Username: "",
            Password: "",
            URL: "",
            Fee: "",
            Address: ""
        };

        $scope.setBitcoinProxy = function(){
            switch($scope.bitmarkConfig.Chain){
            case 'local':
                $scope.bitcoinUseProxy = proxyType.local;
                break;
            default:
                if($scope.bitmarkConfig.Bitcoin.URL == BitmarkProxyURL.testing || $scope.bitmarkConfig.Bitcoin.URL == BitmarkProxyURL.bitmark){
                    $scope.bitcoinUseProxy = proxyType[$scope.bitmarkConfig.Chain];
                }else if($scope.bitmarkConfig.Bitcoin.Username == ProxyTemp.Username){
                    $scope.bitcoinUseProxy = proxyType.other;
                }else{
                    $scope.bitcoinUseProxy = proxyType.local;
                }
            }
        };

        $scope.verifyPassowrd = "";
        $scope.bitcoinPasswordEqual = true;
        $scope.publicKeyPattern = /^(\w|\d|\.|\-|:|\+|=|\^|!|\/|\*|\?|&|<|>|\(|\)|\[|\]|\{|\}|@|%|\$|#)+$/;

        // check bitcoin password
        $scope.$watchGroup(['localBitcoin.Password','verifyPassword'], function(){
            if($scope.bitmarkConfig != null && !passwordVerified($scope.localBitcoin.Password, $scope.verifyPassword)){
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
            $scope.error.show = false;
            saveConfig(function(){
                $scope.goUrl('/main');
            });
        };


        $scope.saveConfigAndStart = function(){
            $scope.error.show = false;
            // send config post api and start bitmark then return to main page
            saveConfig(function(){
                httpService.send("startBitmarkd").then(
                    function(result){
                        $scope.goUrl('/main');
                    }, function(errorMsg){
                        $scope.error.show = true;
                        $scope.error.msg = errorMsg;
                    });
            });
        };

        $scope.goUrl = function(path){
              $location.path(path);
        };

        var saveConfig = function(callBackFunc){
            // set bitcoin object
            switch($scope.bitcoinUseProxy){
            case 'local':
                $scope.bitmarkConfig.Bitcoin = $scope.localBitcoin;
                break;
            case 'other':
                $scope.bitmarkConfig.Bitcoin = $scope.otherProxyTemp;
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
                $scope.error.show = true;
                $scope.error.msg = result.error;
                return;
            }

            var bitmarkConfig = result.bitmarkConfig;
            httpService.send('updateBitmarkConfig', bitmarkConfig).then(
                function(result){
                    if (callBackFunc != undefined){
                        callBackFunc();
                    }
                }, function(errorMsg){
                    $scope.error.show = true;
                    $scope.error.msg = errorMsg;
                });
        };

        var getAndSetBitmarkConfig = function(){
            httpService.send('getBitmarkConfig').then(
                function(result){
                    $scope.bitmarkConfig = initBitmarkConfig(result);

                    // setup bitmark proxy
                    switch($scope.bitmarkConfig.Bitcoin.URL){
                    case $scope.bitmarkTestNetProxyTemp.URL:
                        $scope.bitcoinUseProxy = proxyType.testing;
                        angular.extend($scope.bitmarkTestNetProxyTemp, $scope.bitmarkConfig.Bitcoin);
                        break;
                    case $scope.bitmarkProxyTemp.URL:
                        $scope.bitcoinUseProxy = proxyType.bitmark;
                        angular.extend($scope.bitmarkProxyTemp, $scope.bitmarkConfig.Bitcoin);
                        break;
                    default:
                        if($scope.bitmarkConfig.Bitcoin.Username == ProxyTemp.Username){
                            $scope.bitcoinUseProxy = proxyType.other;
                            angular.extend($scope.otherProxyTemp, $scope.bitmarkConfig.Bitcoin);
                        }else{
                            $scope.bitcoinUseProxy = proxyType.local;
                            angular.extend($scope.localBitcoin, $scope.bitmarkConfig.Bitcoin);
                        }
                    }
                }, function(errorMsg){
                    $scope.error.show = true;
                    $scope.error.msg = errorMsg;
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
            if($scope.bitcoinUseProxy == 'local' && !passwordVerified(bitmarkConfig.Bitcoin.Password, $scope.verifyPassword)){
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
  });
