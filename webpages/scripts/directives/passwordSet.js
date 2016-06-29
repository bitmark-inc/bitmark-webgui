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
                password : '=password',
                result: '=result'
            },
            template: '<table>'+
                '<tr><td><input class="enterpassword" type="password" placeholder="Enter your passcode" ng-model="password" required></td></tr>'+
                '<tr>'+
                '<td ng-class="passwordClass" ng-show="password.length > 0">'+
                '<input class="form-control verifypassword" type="password" placeholder="Verify passcode" ng-model="verifyPassword"></td>'+
                '</tr>'+
                '</table>',
            link: function(scope, element, attrs){
                // check password equality
                scope.verifyPassword = "";
                scope.passwordClass = "";
                scope.result = true;
                scope.$watchGroup(['password','verifyPassword'], function(){
                    scope.result = passwordVerified(scope.password, scope.verifyPassword);
                    // if (!passwordVerified(scope.password, scope.verifyPassword)){
                    //     scope.result = false;
                    // }else{
                    //     scope.result = true;
                    // }
                });

                function passwordVerified(password, verifyPassword){
                    if(password != "" && password != verifyPassword){
                        scope.passwordClass = "has-error has-danger";
                        return false;
                    }
                    scope.passwordClass = "";
                    return true;
                };
            }
        }
    });
