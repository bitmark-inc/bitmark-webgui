// Copyright (c) 2014-2016 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

'use strict';

/**
 * @ngdoc function
 * @name bitmarkWebguiApp.controller:LoginAccessCtrl
 * @description
 * # LoginAccessCtrl
 * Controller of the bitmarkWebguiApp
 */
angular.module('bitmarkWebguiApp')
    .controller('LoginAccessCtrl', ['$scope', '$location', '$cookies', 'httpService', function ($scope, $location, $cookies, httpService) {

        $scope.go = function(){
            $location.path("/main");
        };
  }]);
