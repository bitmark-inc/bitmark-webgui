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
        $scope.showWelcome = false;

        $scope.request = {
            Password: ""
        };
        $scope.errorMsg = "";

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
                    if(result.bitmark_cli_config_file.length == 0) {
                        $scope.showWelcome = true;
                    }else{ // not logout properly, use last time setting
                        $scope.$emit('Authenticated', true);
                        $location.path('/');
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
            if(type == "create" || type == "access") {
                $location.path("/login/"+type);
            }
        };
  }]);
