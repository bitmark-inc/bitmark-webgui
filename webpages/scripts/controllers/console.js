// Copyright (c) 2014-2016 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

'use strict';

/**
 * @ngdoc function
 * @name bitmarkWebguiApp.controller:ConsoleCtrl
 * @description
 * # ConsoleCtrl
 * Controller of the bitmarkWebguiApp
 */
angular.module('bitmarkWebguiApp')
  .controller('ConsoleCtrl', function ($scope, $interval, httpService) {
    $scope.bitmarkNodeIsRunning = false;

    var getInfoPromise;
    var intervalTime = 6 * 1000;


    $scope.init = function () {
      // check bitmarkd status
      httpService.send("statusBitmarkd").then(function (result) {
        if (result == "started") {
          $scope.bitmarkNodeIsRunning = true;
          getInfoPromise = $interval(getBitmarkInfo, intervalTime);
        } else {
          $scope.bitmarkNodeIsRunning = false;
        }
      });
    };

    var getBitmarkInfo = function () {
      httpService.send("getBitmarkdInfo").then(
        function (result) {
          $scope.bitmarkInfo = result;
        },
        function (errorMsg) {
          if (errorMsg != "Failed to connect to bitmarkd") {
            $interval.cancel(getInfoPromise);
          }
        });
    };

    $scope.$on('$destroy', function () {
      $interval.cancel(getInfoPromise);
    });


  });
