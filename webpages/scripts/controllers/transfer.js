// Copyright (c) 2014-2016 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

'use strict';

/**
 * @ngdoc function
 * @name bitmarkWebguiApp.controller:TransferCtrl
 * @description
 * # TransferCtrl
 * Controller of the bitmarkWebguiApp
 */
angular.module('bitmarkWebguiApp')
    .controller('TransferCtrl', function ($scope, $location, httpService, configuration) {
        if(configuration.getConfiguration().bitmarkCliConfigFile.length == 0){
            $location.path('/login');
        }
  });
