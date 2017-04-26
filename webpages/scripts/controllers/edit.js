// Copyright (c) 2014-2016 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

'use strict';

/**
 * @ngdoc function
 * @name bitmarkWebguiApp.controller:EditCtrl
 * @description
 * # EditCtrl
 * Controller of the bitmarkWebguiApp
 */

var defaultBitmarkConfig = {
  "client_rpc": {
    "maximum_connections": 50,
    "announce": ["127.0.0.1:2130"]
  },
  "peering": {
    "announce": {
      "broadcast": ["127.0.0.1:2135"],
      "listen": ["127.0.0.1:2136"]
    }
  },
  "proofing": {
    "publish": ["127.0.0.1:2140"],
    "submit": ["127.0.0.1:2141"]
  }
}

var defaultProoferdConfig = {
  "peering": {
    "connect": [{
      "blocks": "127.0.0.1:2140",
      "submit": "127.0.0.1:2141"
    }]
  }
}

angular.module('bitmarkWebguiApp')
  .controller('EditCtrl', function ($scope, $location, httpService, BitmarkProxyURL, ProxyTemp, utils) {

    $scope.error = {
      show: false,
      msg: ""
    };

    $scope.setErrorMsg = function (show, msg) {
      utils.setErrorMsg($scope.error, show, msg);
    };

    // Check bitamrkd is not running, if it is running, stop it first
    httpService.send('statusBitmarkd').then(
      function (result) {
        if (result.search("stop") >= 0) {
          getAndSetBitmarkConfig();
        } else {
          httpService.send("stopBitmarkd").then(
            function (result) {
              $scope.error.show = true;
              $scope.error.msg = "Bitmarkd has been stopped.";
              getAndSetBitmarkConfig();
            },
            function (errorMsg) {
              $scope.error.show = true;
              $scope.error.msg = errorMsg;
            });
        }
      },
      function (errorMsg) {
        $scope.error.show = true;
        $scope.error.msg = errorMsg;
      }
    );

    // setup proxy temp
    var proxyType = {
      "local": "local",
      "other": "other",
      "testing": "testing",
      "bitmark": "bitmark"
    };

    $scope.bitcoinUseProxy = proxyType.local;

    $scope.otherProxyTemp = angular.copy(ProxyTemp);

    $scope.bitmarkTestNetProxyTemp = angular.copy(ProxyTemp);
    $scope.bitmarkTestNetProxyTemp.URL = BitmarkProxyURL.testing;

    $scope.bitmarkProxyTemp = angular.copy(ProxyTemp);
    $scope.bitmarkProxyTemp.URL = BitmarkProxyURL.bitmark;

    $scope.localBitcoin = {
      Username: "",
      Password: "",
      URL: "",
      Fee: "",
      Address: ""
    };

    $scope.verifyPassowrd = "";
    $scope.bitcoinPasswordEqual = true;
    $scope.publicKeyPattern = /^(\w|\d|\.|\-|:|\+|=|\^|!|\/|\*|\?|&|<|>|\(|\)|\[|\]|\{|\}|@|%|\$|#)+$/;

    // check bitcoin password
    $scope.$watchGroup(['localBitcoin.Password', 'verifyPassword'], function () {
      if ($scope.bitmarkConfig != null && !passwordVerified($scope.localBitcoin.Password, $scope.verifyPassword)) {
        $scope.bitcoinPasswordEqual = false;
      } else {
        $scope.bitcoinPasswordEqual = true;
      }
    });

    $scope.deleteItem = function (list, index) {
      list.splice(index, 1);
    };

    $scope.addItem = function (list, limit) {
      if (typeof limit !== "number" || list.length < limit) {
        list.splice(list.length, 0, "");
      }
    };

    $scope.saveConfig = function () {
      $scope.error.show = false;
      saveConfig(function () {
        $scope.goUrl('/main');
      });
    };


    // $scope.saveConfigAndStart = function(){
    //     $scope.error.show = false;
    //     // send config post api and start bitmark then return to main page
    //     saveConfig(function(){
    //         httpService.send("startBitmarkd").then(
    //             function(result){
    //                 $scope.goUrl('/main');
    //             }, function(errorMsg){
    //                 $scope.error.show = true;
    //                 $scope.error.msg = errorMsg;
    //             });
    //     });
    // };

    $scope.goUrl = function (path) {
      $location.path(path);
    };

    var saveConfig = function (callBackFunc) {
      var bitmarkConfig = $scope.bitmarkConfig;
      var prooferdConfig = $scope.prooferdConfig;

      var configs = {
        bitmarkConfig: bitmarkConfig,
        prooferdConfig: prooferdConfig
      };
      httpService.send('updateBitmarkConfig', configs).then(
        function (result) {
          if (callBackFunc != undefined) {
            callBackFunc();
          }
        },
        function (errorMsg) {
          $scope.error.show = true;
          $scope.error.msg = errorMsg;
        });
    };

    var getAndSetBitmarkConfig = function () {
      var error = [];
      httpService.send('getBitmarkConfig').then(
        function (results) {
          if (results.bitmarkd.err) {
            $scope.bitmarkConfig = defaultBitmarkConfig;
            $scope.setErrorMsg(true, results.bitmarkd.err);
          } else {
            $scope.bitmarkConfig = checkConfig(results.bitmarkd.data);
          }

          if (results.prooferd.err || Object.keys(results.prooferd).length == 0) {
            $scope.prooferdConfig = defaultProoferdConfig
            $scope.setErrorMsg(true, results.prooferd.err);
          } else {
            $scope.prooferdConfig = results.prooferd.data;
          }
        },
        function (errorMsg) {
          $scope.error.show = true;
          $scope.error.msg = errorMsg;
          $scope.bitmarkConfig = defaultBitmarkConfig;
          $scope.prooferdConfig = defaultProoferdConfig;
        });
    };

    var initConfig = function (bitmarkConfig) {
      // give empty array for null field
      var checkItems = ["ClientRPC", "Peering"];
      var checkFields = ["Listen", "Announce", "Connect", "Broadcast"];

      for (var i = 0; i < checkItems.length; i++) {
        var checkItem = checkItems[i];
        if (!bitmarkConfig[checkItem]) {
          continue
        }
        for (var j = 0; j < checkFields.length; j++) {
          var checkField = checkFields[j];
          if (bitmarkConfig[checkItem][checkField] !== undefined && bitmarkConfig[checkItem][checkField] == null) {
            bitmarkConfig[checkItem][checkField] = [];
          }
        }
      }
      return bitmarkConfig;
    };

    var passwordVerified = function (password, verifyPassword) {
      if (password != "" && password != verifyPassword) {
        return false;
      }
      return true;
    };

    // return {bitmarkConfig:{}, error:""}
    var checkConfig = function (config) {

      var clientAnnounce = config.client_rpc.announce
      var peerAnnounce = config.peering.announce

      if (!peerAnnounce.broadcast || peerAnnounce.broadcast.length == 0) {
        peerAnnounce.broadcast = ["127.0.0.1:2135"]
      }

      if (!peerAnnounce.listen || peerAnnounce.listen.length == 0) {
        peerAnnounce.listen = ["127.0.0.1:2136"]
      }

      if (!clientAnnounce || clientAnnounce.length == 0) {
        config.client_rpc.announce = ["127.0.0.1:2130"]
      }

      return config;
    };
  });
