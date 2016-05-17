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
    .controller('SystemCtrl', ['$scope', '$http', '$location', 'httpService', function ($scope, $http, $location, httpService) {
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
                    $scope.goUrl('/');
                }, function(errorMsg){
                    $scope.errorMsg = "Set password failed";
                });
        };

        $scope.goUrl = function(path){
            $location.path(path);
        };
  }]);
