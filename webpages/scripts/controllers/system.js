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
        $scope.verifyPassword = "";
        $scope.bitmarkWebguiPasswordEqual = true;
        $scope.errorMsg = "";

        $scope.save = function(){
            if($scope.request.Origin == "" || $scope.request.New == ""){
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

        // check password equality
        $scope.$watchGroup(['request.New','verifyPassword'], function(){
            if (!passwordVerified($scope.request.New, $scope.verifyPassword)){
                $scope.bitmarkWebguiPasswordEqual = false;
            }else{
                $scope.bitmarkWebguiPasswordEqual = true;
            }
        });

        function passwordVerified(password, verifyPassword){
            if(password != "" && password != verifyPassword){
              return false;
            }
            return true;
        };
  }]);
