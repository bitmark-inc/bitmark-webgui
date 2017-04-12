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
  .controller('MainCtrl', function ($scope, $location, $uibModal, httpService, ProxyTemp, $interval, utils) {

    $scope.disableStart = true;
    $scope.disableStop = true;

    $scope.error = {
      show: false,
      msg: ""
    };

    $scope.setErrorMsg = function (show, msg) {
      utils.setErrorMsg($scope.error, show, msg);
    };

    var getInfoPromise;
    var intervalTime = 6 * 1000;
    var bitmarkdisRunning = false;
    httpService.send('statusBitmarkd').then(
      function (result) {
        // check status and set disable button
        if (result.search("stop") >= 0) {
          disableStartBitmarkBtn(false);
        } else { // bitmarkd is running
          getBitmarkInfo();
          getInfoPromise = $interval(getBitmarkInfo, intervalTime);
        }
      },
      function (errorMsg) {
        $scope.setErrorMsg(true, errorMsg);
        $interval.cancel(getInfoPromise);
      });


    httpService.send('getBitmarkConfig').then(
      function (result) {
        $scope.bitmarkConfig = result;

        $scope.showOptionBitcoinItems = true;
        if (result.Bitcoin.Username == ProxyTemp.Username) {
          $scope.showOptionBitcoinItems = false;
        }

      },
      function (errorMsg) {
        $scope.setErrorMsg(true, errorMsg);
      });

    $scope.startBitmark = function () {
      allBitmarkdDisable();
      $scope.setErrorMsg(false, "");
      httpService.send("startBitmarkd").then(
        function (result) {
          // disable bitmark start button
          if (result.search("start running bitmarkd") >= 0) {
            getBitmarkInfo();
            getInfoPromise = $interval(getBitmarkInfo, intervalTime);
          } else {
            disableStartBitmarkBtn(false);
          }
        },
        function (errorMsg) {
          disableStartBitmarkBtn(false);
          $scope.setErrorMsg(true, errorMsg);
          $interval.cancel(getInfoPromise);
        });
    };
    $scope.stopBitmark = function () {
      allBitmarkdDisable();
      $scope.setErrorMsg(false, "");
      $interval.cancel(getInfoPromise);

      // check bitmark mode
      //  mode: resynchronize: show modal to alert

      $scope.bitmarkInfo = undefined;
      httpService.send("stopBitmarkd").then(
        function (result) {
          if (result.search("stop running bitmarkd") >= 0) {
            disableStartBitmarkBtn(false);
          } else {
            disableStartBitmarkBtn(true);
          }
        },
        function (errorMsg) {
          disableStartBitmarkBtn(true);
          $scope.setErrorMsg(true, errorMsg);
        });
    };
    $scope.configBitmark = function () {
      // check bitmark mode
      //  mode Resynchronise: show modal to alert user bitmarkd will be stopped and data will be removed
      //  mode normal: show modal to alert user bitmarkd will be stopped
      if (bitmarkdisRunning) {
        showStopBitmarkdModal('config');
      } else {
        $location.path('/edit');
      }
    };


    var showStopBitmarkdModal = function (type) {
      var modalInstance = $uibModal.open({
        templateUrl: 'views/stopBitmarkdModal.html',
        controller: 'StopBitmarkdModalCtrl',
        resolve: {
          type: function () {
            return type;
          }
        }
      });

      modalInstance.result.then(function () {
        // kill the stop
        allBitmarkdDisable();
        $scope.setErrorMsg(false, "");
        $interval.cancel(getInfoPromise);
        $scope.bitmarkInfo = undefined;
        httpService.send("stopBitmarkd").then(
          function (result) {
            if (result.search("stop running bitmarkd") >= 0) {
              disableStartBitmarkBtn(false);
              $location.path('/edit');
            } else {
              disableStartBitmarkBtn(true);
            }
          },
          function (errorMsg) {
            disableStartBitmarkBtn(true);
            $scope.setErrorMsg(true, errorMsg);
          });
      });
    };

    var disableStartBitmarkBtn = function (startDisableBool) {
      bitmarkdisRunning = startDisableBool;
      $scope.disableStart = startDisableBool;
      $scope.disableStop = !startDisableBool;
    };

    var allBitmarkdDisable = function () {
      bitmarkdisRunning = true;
      $scope.disableStart = true;
      $scope.disableStop = true;
    };

    var getBitmarkInfo = function () {
      httpService.send("getBitmarkdInfo").then(
        function (result) {
          $scope.setErrorMsg(false, "");
          $scope.bitmarkInfo = result;
          if (result.mode == undefined) {
            allBitmarkdDisable();
          } else {
            disableStartBitmarkBtn(true);
          }
        },
        function (errorMsg) {
          if (errorMsg != "Failed to connect to bitmarkd") {
            $scope.setErrorMsg(true, errorMsg);
            $interval.cancel(getInfoPromise);
          }
        });
    };

    $scope.$on('$destroy', function () {
      $interval.cancel(getInfoPromise);
    });
  });
