// Copyright (c) 2014-2016 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

'use strict';

/**
 * @ngdoc function
 * @name bitmarkWebguiApp.controller:ChainCtrl
 * @description
 * # ChainCtrl
 * Controller of the bitmarkWebguiApp
 */
angular.module('bitmarkWebguiApp')
  .controller('ChainCtrl', function ($scope, $location, $q, httpService, BitmarkdNetwork, utils) {
    $scope.init = function () {
      // check bitmarkd status
      httpService.send("statusBitmarkd").then(function (result) {
        if (result == "started") {
          // go to main view
          $scope.$emit('Authenticated', true);
          $location.path("/main");
        }
      });
    };

    $scope.request = {
      chain: "testing",
      running: false
    };

    $scope.error = {
      show: false,
      msg: ""
    };

    $scope.setErrorMsg = function (show, msg) {
      utils.setErrorMsg($scope.error, show, msg);
    };

    $scope.startNode = function () {
      $scope.request.running = true;
      $scope.setErrorMsg(false, '');
      // setup bitmarkd config

      var setupBitmarkdPromise = httpService.send('setupBitmarkd', {
        network: BitmarkdNetwork[$scope.request.chain]
      }).catch(function (error) {
        return {
          error: error
        }
      })
      var setupProoferdPromise = httpService.send('setupProoferd', {
        network: BitmarkdNetwork[$scope.request.chain]
      }).catch(function (error) {
        return {
          error: error
        }
      })

      $q.all([setupBitmarkdPromise, setupProoferdPromise])
        .then(function (results) {
          var otherErrors = []
          for (var i = results.length - 1; i >= 0; i--) {
            if (results[i] && results[i].result) {
              otherErrors.push(results[i].result)
            }
          }
          if (otherErrors.length > 0) {
            $scope.setErrorMsg(true, otherErrors);
          } else {
            $location.path("/edit");
            $scope.$emit('Authenticated', true);
          }
        })
        .catch(function (error) {
          if (error instanceof Array) {
            $scope.setErrorMsg(true, error);
          } else if (error instanceof Error) {
            $scope.setErrorMsg(true, error.message);
          } else {
            $scope.setErrorMsg(true, "Unexcepted error");
            console.log("Unexcepted error:", err)
          }
          $scope.request.running = false;
        });
    };
  });
