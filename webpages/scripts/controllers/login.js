// Copyright (c) 2014-2016 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

'use strict';

/**
 * @ngdoc function
 * @name bitmarkWebguiApp.controller:LoginCtrl
 * @description
 * # LoginCtrl
 * Controller of the bitmarkWebguiApp
 */
angular.module('bitmarkWebguiApp')
    .controller('LoginCtrl', function ($scope, $location, $cookies, httpService, $log, utils) {

        $scope.request = {
            Password: ""
        };

        var checkBitmarkdStatus = function(){
            httpService.send("statusBitmarkd").then(function(result){
                if(result == "started") {
                    // go to main view
                    $scope.goUrl("main");
                } else {
                    // go to chain view
                    $scope.goUrl("chain");
                }
            });
        };

        $scope.init = function() {
            httpService.send("checkAuthenticate").then(
                function(result){ // already login
                    // check bitmarkd status
                    checkBitmarkdStatus();
                },function(){
                    $scope.$emit('Authenticated', false);
            });
        };

        $scope.error = {
            show: false,
            msg: ""
        };

        $scope.setErrorMsg = function(show, msg) {
            utils.setErrorMsg($scope.error, show, msg);
        };

        $scope.login = function(){
            if($scope.request.Password == ""){
                $scope.setErrorMsg(true, "Please enter your password.");
                return;
            }

            // Clean cookie first
            $cookies.remove('bitmark-webgui');
            httpService.send('login', $scope.request).then(
                function(result){
                    checkBitmarkdStatus();
                }, function(result){
                    if (result == "Already logged in") {
                        $scope.setErrorMsg(true, "Already login! should not see this page");
                    }else{
                        $scope.$emit('Authenticated', false);
                        $scope.setErrorMsg(true, "Login failed, please try again.");
                    }
                });
        };


        $scope.goUrl = function(type){
            switch(type){
            case "main":
                $scope.$emit('Authenticated', true);
                $location.path("/main");
                break;
            case "chain":
                $location.path("/chain");
            };
        };
  });
