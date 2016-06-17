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
    .controller('LoginCtrl', ['$scope', '$location', '$cookies', 'httpService', '$log', 'configuration', function ($scope, $location, $cookies, httpService, $log, configuration) {

        $scope.request = {
            Password: ""
        };
        $scope.errorMsg = "";

        $scope.init = function() {
            httpService.send("checkAuthenticate").then(
                function(result){ // already login
                    $scope.showWelcome = true;
                    if(result.bitmark_cli_config_file.length != 0) {
                        $scope.showWarning = true;
                    }
                },function(){
                    $scope.$emit('Authenticated', false);
                    $scope.showWelcome = false;
            });
        };

        $scope.login = function(){
            if($scope.request.Password == ""){
                $scope.errorMsg = "Please enter password";
                return;
            }
            // Clean cookie first
            $cookies.remove('bitmark-webgui');
            httpService.send('login', $scope.request).then(
                function(result){
                    configuration.setChain(result.chain);
                    configuration.setBitmarkCliConfigFile(result.bitmark_cli_config_file);

                    $scope.showWelcome = true;
                    if(result.bitmark_cli_config_file.length != 0) {
                        $scope.showWarning = true;
                    }
                }, function(result){
                    if (result == "Already logged in") {
                        $log.error("Already login!");
                        $scope.errorMsg = "Already login! should not see this page";
                    }else{
                        $scope.$emit('Authenticated', false);
                        $scope.errorMsg = "Login failed";
                    }
                });
        };


        $scope.goUrl = function(type){
            switch(type){
            case "create":
                $location.path("/login/create");
                break;
            case "access":
                $location.path("/login/access");
                break;
            case "logout":
                $location.path("/logout");
                break;
            case "main":
                $scope.$emit('Authenticated', true);
                $location.path("/main");
                break;
            };
        };
  }]);
