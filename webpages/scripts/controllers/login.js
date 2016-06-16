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
    .controller('LoginCtrl', ['$scope', '$location', '$cookies', 'httpService', function ($scope, $location, $cookies, httpService) {
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
                    $scope.showWelcome = true;
                }, function(result){
                    if (result == "Already logged in") {
                        $scope.$emit('Authenticated', true);
                        $location.path('/');
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
