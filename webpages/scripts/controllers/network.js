// Copyright (c) 2014-2016 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

'use strict';

/**
 * @ngdoc function
 * @name bitmarkWebguiApp.controller:NetworkCtrl
 * @description
 * # MainCtrl
 * Controller of the bitmarkWebguiApp
 */
angular.module('bitmarkWebguiApp')
    .controller('NetworkCtrl', function ($scope, $http, $location, httpService, configuration) {
        if(configuration.getConfiguration().bitmarkCliConfigFile.length == 0){
            $location.path('/login');
        }

        $scope.request = {
            Origin: "",
            New: ""
        };
        $scope.verifyPassword = true;
        $scope.bitmarkWebguiPasswordEqual = true;
        $scope.errorMsg = "";

        $scope.save = function(){
            if($scope.request.Origin == "" || $scope.request.New == "" || !$scope.verifyPassword){
                $scope.errorMsg = "All fields should be filled";
                return;
            }
            httpService.send('updateBitmarkWebguiPassword', $scope.request).then(
                function(result){
                    $scope.goUrl('/main');
                }, function(errorMsg){
                    $scope.errorMsg = "Set password failed";
                });
        };

        $scope.goUrl = function(path){
            $location.path(path);
        };
  });
