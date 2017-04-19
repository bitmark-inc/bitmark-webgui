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

    $scope.disableBitmarkdStart = true;
    $scope.disableBitmarkdStop = true;
    $scope.disableProoferdStart = true;
    $scope.disableProoferdStop = true;

    $scope.error = {
      show: false,
      msg: ""
    };

    $scope.setErrorMsg = function (show, msg) {
      utils.setErrorMsg($scope.error, show, msg);
    };

    var getInfoPromise;
    var getProoferdStatusPromise;
    var intervalTime = 6 * 1000;
    var bitmarkdisRunning = false;
    var prooferdisRunning = false;
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

    getProoferdStatus()
    getProoferdStatusPromise = $interval(getProoferdStatus, intervalTime);

    httpService.send('getBitmarkConfig').then(
      function (result) {
        var errors = []
        if (result.bitmarkd.error) {
          errors.push(result.bitmarkd.error)
        } else {
          $scope.bitmarkConfig = result.bitmarkd.data;
        }

        if (result.prooferd.error) {
          errors.push(result.prooferd.error)
        } else {
          $scope.prooferdConfig = result.prooferd.data;
        }

        if (errors.length > 0) {
          $scope.setErrorMsg(true, errors);
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

    $scope.startProoferd = function () {
      allProoferdDisable();
      $scope.setErrorMsg(false, "");
      httpService.send("startProoferd").then(
        function (result) {
          // disable bitmark start button
          if (result.search("start running prooferd") >= 0) {
            $scope.prooferdStatus = "running"
          } else {
            disableStartProoferdBtn(false);
          }
        },
        function (errorMsg) {
          disableStartProoferdBtn(false);
          $scope.setErrorMsg(true, errorMsg);
        });
    }

    $scope.stopProoferd = function () {
      allProoferdDisable();
      $scope.setErrorMsg(false, "");

      $scope.bitmarkInfo = undefined;
      httpService.send("stopProoferd").then(
        function (result) {
          if (result.search("stop running prooferd") >= 0) {
            disableStartProoferdBtn(false);
          } else {
            disableStartProoferdBtn(true);
          }
        },
        function (errorMsg) {
          disableStartProoferdBtn(true);
          $scope.setErrorMsg(true, errorMsg);
        });
    }

    $scope.configProoferd = function () {

    }

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
      $scope.disableBitmarkdStart = startDisableBool;
      $scope.disableBitmarkdStop = !startDisableBool;
    };

    var allBitmarkdDisable = function () {
      bitmarkdisRunning = true;
      $scope.disableBitmarkdStart = true;
      $scope.disableBitmarkdStop = true;
    };


    var disableStartProoferdBtn = function (startDisableBool) {
      prooferdisRunning = startDisableBool;
      $scope.disableProoferdStart = startDisableBool;
      $scope.disableProoferdStop = !startDisableBool;
    };


    var allProoferdDisable = function () {
      prooferdisRunning = true;
      $scope.disableProoferdStart = true;
      $scope.disableProoferdStop = true;
    };


    function getProoferdStatus() {
      httpService.send('statusProoferd').then(
        function (result) {
          if (result.search("stop") >= 0) {
            disableStartProoferdBtn(false);
          } else {
            disableStartProoferdBtn(true);
          }
        },
        function (errorMsg) {
          $scope.setErrorMsg(true, errorMsg);
        });
    }

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
      $interval.cancel(getProoferdStatusPromise);
    });
  });
