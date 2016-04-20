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
    .controller('LoginCtrl', ['$scope', '$location', 'httpService', function ($scope, $location, httpService) {
        httpService.send('checkAuthenticate').then(
            function(){
                // already logined
                $scope.goUrl('/');
            },function(){}
        );

        $scope.request = {
            Password: ""
        };
        $scope.errorMsg = "";

        $scope.login = function(){
            if($scope.request.Password == ""){
                $scope.errorMsg = "Please enter password";
                return;
            }
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
