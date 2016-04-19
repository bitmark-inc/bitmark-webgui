// Copyright (c) 2014-2016 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

'use strict';

/**
 * @ngdoc function
 * @name bitmarkMgmtApp.controller:MainCtrl
 * @description
 * # MainCtrl
 * Controller of the bitmarkMgmtApp
 */
angular.module('bitmarkMgmtApp')
    .controller('MainCtrl', ['$scope', '$location', 'httpService', function ($scope, $location, httpService) {

        var bitmarkStatusObj = {
            "run": "Running",
            "stop": "Stopped",
            "error": "Error"
        };

        $scope.disableStart = true;
        $scope.disableStop = true;
        $scope.bitmarkStatus = "Running";


        httpService.send('statusBitmarkd').then(
            function(result){
                // TODO: check status and set disable button
                if(result.search("stop") >= 0) {
                    setBitmarkdDisable(false);
                    $scope.bitmarkStatus = bitmarkStatusObj.stop;
                }else{
                    setBitmarkdDisable(true);
                    $scope.bitmarkStatus = bitmarkStatusObj.run;
                    getBitmarkInfo();
                }
            }, function(errorMsg){
                $scope.errorMsg = errorMsg;
            });


        httpService.send('getBitmarkConfig').then(
            function(result){
                $scope.bitmarkConfig = result;
            },function(errorMsg){
                $scope.errorMsg = errorMsg;
            });

        $scope.startBitmark = function(){
            allBitmarkdDisable();
            httpService.send("startBitmarkd").then(
                function(result){
                    // disable bitmark start button
                    if(result.search("start running bitmarkd")>= 0){
                        $scope.bitmarkStatus = bitmarkStatusObj.run;
                        setBitmarkdDisable(true);
                        getBitmarkInfo();
                    }else{
                        setBitmarkdDisable(false);
                    }
                }, function(errorMsg){
                    setBitmarkdDisable(false);
                    $scope.bitmarkStatus = bitmarkStatusObj.error;
                    $scope.errorMsg = errorMsg;
                });
        };
        $scope.stopBitmark = function(){
            allBitmarkdDisable();
            httpService.send("stopBitmarkd").then(
                function(result){
                    if(result.search("stop running bitmarkd")>=0) {
                        $scope.bitmarkStatus = bitmarkStatusObj.stop;
                        setBitmarkdDisable(false);
                    }else{
                        setBitmarkdDisable(true);
                    }
                }, function(errorMsg){
                    setBitmarkdDisable(true);
                    $scope.bitmarkStatus = bitmarkStatusObj.error;
                    $scope.errorMsg = errorMsg;
                });
        };
        $scope.goUrl = function(path){
            $location.path(path);
        };

        var setBitmarkdDisable = function(startDisableBool) {
            $scope.disableStart = startDisableBool;
            $scope.disableStop = !startDisableBool;
        };

        var allBitmarkdDisable = function() {
            $scope.disableStart = true;
            $scope.disableStop = true;
        };

        var getBitmarkInfo = function(){
            httpService.send("getBitmarkdInfo").then(
                function(result){
                    $scope.bitmarkInfo = result;
                },
                function(errorMsg){
                    $scope.errorMsg = errorMsg;
                });
        };
  }]);
