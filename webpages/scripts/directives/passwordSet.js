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
    .directive('passwordSet', function(){
        return{
            restrict : 'EA',
            replace : true,
            transclude : true,
            scope : {
                password : '=password'
            },
            template: '<table>'+
                '<tr>'+
                '<td><input type="password" ng-model="password" required></td>'+
                '<td ng-if="password.length > 0">'+
                '<b>verify</b></td>'+
                '<td ng-show="password.length > 0">'+
                '<input type="password" ng-model="verifyPassword"></td>'+
                '<td ng-show="!passwordEqual">not equal</td>'+
                '</tr>'+
                '</table>'+
                '{{password}}'+
                'verifiy: {{verifyPassword}}',
            link: function(scope, element, attrs){
                // check password equality
                scope.verifyPassword = "";
                scope.$watchGroup(['password','verifyPassword'], function(){
                    if (!passwordVerified(scope.password, scope.verifyPassword)){
                        scope.passwordEqual = false;
                    }else{
                        scope.passwordEqual = true;
                    }
                });

                function passwordVerified(password, verifyPassword){
                    if(password != "" && password != verifyPassword){
                        return false;
                    }
                    return true;
                };
            }
        }
    });
