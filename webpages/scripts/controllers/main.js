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
    .controller('MainCtrl', ['$scope', '$location', 'httpService', 'ProxyTemp', '$interval', function ($scope, $location, httpService, ProxyTemp, $interval) {

        var bitmarkStatusObj = {
            "run": "Running",
            "stop": "Stopped",
            "error": "Error"
        };

        $scope.disableStart = true;
        $scope.disableStop = true;
        $scope.bitmarkStatus = "Running";
        $scope.error = {
            show: false,
            msg: ""
        };

        var getInfoPromise;
        var intervalTime = 6 * 1000;
        $scope.$on('$destroy', function(){
            console.log("cancel promise");
            $interval.cancel(getInfoPromise);
        });

        httpService.send('statusBitmarkd').then(
            function(result){
                // TODO: check status and set disable button
                if(result.search("stop") >= 0) {
                    setBitmarkdDisable(false);
                    $scope.bitmarkStatus = bitmarkStatusObj.stop;
                }else{
                    $scope.bitmarkStatus = bitmarkStatusObj.run;
                    getBitmarkInfo();
                    getInfoPromise = $interval(getBitmarkInfo, intervalTime);
                }
            }, function(errorMsg){
                $scope.error.show = true;
                $scope.error.msg = errorMsg;
                $interval.cancel(getInfoPromise);
            });


        httpService.send('getBitmarkConfig').then(
            function(result){
                $scope.bitmarkConfig = result;

                $scope.showOptionBitcoinItems = true;
                if(result.Bitcoin.Username == ProxyTemp.Username){
                        $scope.showOptionBitcoinItems = false;
                }

            },function(errorMsg){
                $scope.error.show = true;
                $scope.error.msg = errorMsg;
            });

        $scope.startBitmark = function(){
            allBitmarkdDisable();
            $scope.error.show = false;
            httpService.send("startBitmarkd").then(
                function(result){
                    // disable bitmark start button
                    if(result.search("start running bitmarkd")>= 0){
                        $scope.bitmarkStatus = bitmarkStatusObj.run;
                        getBitmarkInfo();
                        getInfoPromise = $interval(getBitmarkInfo, intervalTime);
                    }else{
                        setBitmarkdDisable(false);
                    }
                }, function(errorMsg){
                    setBitmarkdDisable(false);
                    $scope.bitmarkStatus = bitmarkStatusObj.error;
                    $scope.error.show = true;
                    $scope.error.msg = errorMsg;
                    $interval.cancel(getInfoPromise);
                });
        };
        $scope.stopBitmark = function(){
            allBitmarkdDisable();
            $scope.error.show = false;
            $interval.cancel(getInfoPromise);
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
                    $scope.error.show = true;
                    $scope.error.msg = errorMsg;
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
                    console.log("get bitmark info");
                    $scope.error.show = false;
                    $scope.bitmarkInfo = result;
                    if(result.mode == undefined || result.mode !== 'Normal'){
                        allBitmarkdDisable();
                    }else{
                        setBitmarkdDisable(true);
                    }
                },
                function(errorMsg){
                    if(errorMsg != "Failed to connect to bitmarkd") {
                        $scope.error.show = true;
                        $scope.error.msg = errorMsg;
                    }
                });
        };

  }]);
