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

        $scope.disableStart = true;
        $scope.disableStop = true;
        $scope.bitmarkStatus = "Running";


        httpService.send('statusBitmarkd').then(
            function(result){
                // TODO: check status and set disable button
                if(result.search("stop") >= 0) {
                    setBitmarkdDisable(false);
                    $scope.bitmarkStatus = "Stopped";
                }else{
                    setBitmarkdDisable(true);
                    $scope.bitmarkStatus = "Running";
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
                        $scope.bitmarkStatus = "Running";
                        setBitmarkdDisable(true);
                        getBitmarkInfo();
                    }else{
                        setBitmarkdDisable(false);
                    }
                }, function(errorMsg){
                    $scope.bitmarkStatus = "Error";
                    $scope.errorMsg = errorMsg;
                });
        };
        $scope.stopBitmark = function(){
            allBitmarkdDisable();
            httpService.send("stopBitmarkd").then(
                function(result){
                    if(result.search("stop running bitmarkd")>=0) {
                        $scope.bitmarkStatus = "Stopped";
                        setBitmarkdDisable(false);
                    }else{
                        setBitmarkdDisable(true);
                    }
                }, function(errorMsg){
                    $scope.bitmarkStatus = "Error";
                    $scope.errorMsg = errorMsg;
                });
        };
        $scope.goUrl = function(path){
            $location.path(path);
        };

        $scope.logout = function(){
            httpService.send("logout").then(
                function(){
                    $scope.goUrl('/login');
                }, function(errorMsg){
                    $scope.errorMsg = errorMsg;
                });
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
