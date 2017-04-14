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

      function isNotFoundError(err) {
        return err.search("not found") >= 0
      }

      $q.all([setupBitmarkdPromise, setupProoferdPromise])
        .then(function (results) {
          var setupBitmarkdResult = results[0],
            setupProoferdResult = results[1];
          if (setupBitmarkdResult || setupProoferdResult) {
            var otherErrors = []
            if ((setupBitmarkdResult && !isNotFoundError(setupBitmarkdResult.error))) {
              otherErrors.push(setupBitmarkdResult.error)
            }
            if ((setupProoferdResult && !isNotFoundError(setupProoferdResult.error))) {
              otherErrors.push(setupProoferdResult.error)
            }
            if (otherErrors.length > 0) {
              $scope.setErrorMsg(true, otherErrors);
            } else {
              $location.path("/edit");
              $scope.$emit('Authenticated', true);
            }
          } else {
            $location.path("/main");
            $scope.$emit('Authenticated', true);
          }
        })
        .catch(function (errors) {
          $scope.setErrorMsg(true, errors);
          $scope.request.running = false;
        });
    };
  });
