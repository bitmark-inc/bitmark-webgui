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
    .controller('MainCtrl', function ($scope, $location, httpService, ProxyTemp, $interval, configuration, utils) {

        // var bitmarkStatusObj = {
        //     "run": "● RUNNING",
        //     "stop": "● STOPPED",
        //     "error": "● ERROR",
        //     "resync": "RESYNCHRONIZING.."
        // };

        $scope.disableStart = true;
        $scope.disableStop = true;
        // $scope.bitmarkStatus = bitmarkStatusObj.resync;
        // $scope.bitmarkStatusStyle = {"color":"#4A90E2"};
        $scope.error = {
            show: false,
            msg: ""
        };

        $scope.setErrorMsg = function(show, msg){
            utils.setErrorMsg($scope.error, show, msg);
        };

        var getInfoPromise;
        var intervalTime = 6 * 1000;

        httpService.send('statusBitmarkd').then(
            function(result){
                // TODO: check status and set disable button
                if(result.search("stop") >= 0) {
                    disableStartBitmarkBtn(false);
                    // $scope.bitmarkStatus = bitmarkStatusObj.stop;
                    // $scope.bitmarkStatusStyle.color = "#FF0000";
                }else{ // bitmarkd is running
                    getBitmarkInfo();
                    getInfoPromise = $interval(getBitmarkInfo, intervalTime);
                }
            }, function(errorMsg){
                $scope.setErrorMsg(true, errorMsg);
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
                $scope.setErrorMsg(true, errorMsg);
            });

        $scope.startBitmark = function(){
            allBitmarkdDisable();
            $scope.setErrorMsg(false, "");
            httpService.send("startBitmarkd").then(
                function(result){
                    // disable bitmark start button
                    if(result.search("start running bitmarkd")>= 0){
                        getBitmarkInfo();
                        getInfoPromise = $interval(getBitmarkInfo, intervalTime);
                    }else{
                        disableStartBitmarkBtn(false);
                    }
                }, function(errorMsg){
                    disableStartBitmarkBtn(false);
                    // $scope.bitmarkStatus = bitmarkStatusObj.error;
                    // $scope.bitmarkStatusStyle.color = "#FF0000";
                    $scope.setErrorMsg(true, errorMsg);
                    $interval.cancel(getInfoPromise);
                });
        };
        $scope.stopBitmark = function(){
            allBitmarkdDisable();
            $scope.setErrorMsg(false, "");
            $interval.cancel(getInfoPromise);
            $scope.bitmarkInfo = undefined;
            httpService.send("stopBitmarkd").then(
                function(result){
                    if(result.search("stop running bitmarkd")>=0) {
                        // $scope.bitmarkStatus = bitmarkStatusObj.stop;
                        // $scope.bitmarkStatusStyle.color = "#FF0000";
                        disableStartBitmarkBtn(false);
                    }else{
                        disableStartBitmarkBtn(true);
                    }
                }, function(errorMsg){
                    disableStartBitmarkBtn(true);
                    // $scope.bitmarkStatus = bitmarkStatusObj.error;
                    // $scope.bitmarkStatusStyle.color = "#FF0000";
                    $scope.setErrorMsg(true, errorMsg);
                });
        };
        $scope.goUrl = function(path){
            $location.path(path);
        };

        var disableStartBitmarkBtn = function(startDisableBool) {
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
                    $scope.setErrorMsg(false, "");
                    $scope.bitmarkInfo = result;
                    if(result.mode == undefined){
                        allBitmarkdDisable();
                        // $scope.bitmarkStatus = bitmarkStatusObj.resync;
                        // $scope.bitmarkStatusStyle.color = "#4A90E2";
                    }else{
                        disableStartBitmarkBtn(true);
                        // $scope.bitmarkStatus = bitmarkStatusObj.run;
                        // $scope.bitmarkStatusStyle.color = "#7ED321";
                    }
                },
                function(errorMsg){
                    if(errorMsg != "Failed to connect to bitmarkd") {
                        $scope.setErrorMsg(true, errorMsg);
                        $interval.cancel(getInfoPromise);
                    }
                });
        };

        $scope.$on('$destroy', function(){
            $interval.cancel(getInfoPromise);
        });
  });
