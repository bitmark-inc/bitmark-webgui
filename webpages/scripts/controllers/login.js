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
    .controller('LoginCtrl', ['$scope', '$location', '$cookies', 'httpService', function ($scope, $location, $cookies, httpService) {
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
                    $scope.$emit('Authenticated', true);
                    $scope.goUrl('/');
                }, function(result){
                    if (result == "Already logged in") {
                        $scope.$emit('Authenticated', true);
                        $scope.goUrl('/');
                    }else{
                        $scope.$emit('Authenticated', false);
                        $scope.errorMsg = "Login failed";
                    }

                });
        };

        $scope.goUrl = function(path){
            $location.path(path);
        };

  }]);
