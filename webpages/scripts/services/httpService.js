// Copyright (c) 2014-2016 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

app.factory('httpService', function($http, $q, $location, $rootScope){
    var hostApiPath = "/api";

    var API = {
        getBitmarkConfig: {
            method: 'GET',
            url: hostApiPath+'/config'
        },
        updateBitmarkConfig: {
            method: 'POST',
            url: hostApiPath+'/config'
        },
        updateBitmarkWebguiPassword: {
            method: 'POST',
            url: hostApiPath+'/password'
        },
        login: {
            method: 'POST',
            url: hostApiPath+'/login'
        },
        checkAuthenticate: {
            method: 'GET',
            url: hostApiPath+'/login'
        },
        logout: {
            method: 'POST',
            url: hostApiPath+'/logout'
        },
        startBitmarkd: {
            method: 'POST',
            url: hostApiPath+'/bitmarkd',
            data:{
                option: "start"
            }
        },
        stopBitmarkd: {
            method: 'POST',
            url: hostApiPath+'/bitmarkd',
            data:{
                option: "stop"
            }
        },
        statusBitmarkd: {
            method: 'POST',
            url: hostApiPath+'/bitmarkd',
            data:{
                option: "status"
            }
        },
        getBitmarkdInfo: {
            method: 'POST',
            url: hostApiPath+'/bitmarkd',
            data:{
                option: "info"
            }
        }
    };

    return {
        send: function(api, data){
            var deferred = $q.defer();

            var apiConfig = angular.copy(API[api]);
            if( data != undefined) {
                apiConfig.data = data;
            }


            $http(apiConfig).then(function successCallback(response) {
                if (response.data.ok) {
                    deferred.resolve(response.data.result);
                }else {
                    var errorMsg = "";
                    switch(api){
                    case 'getBitmarkConfig':
                        errorMsg = "Failed to get bitmark config";
                        break;
                    case 'updateBitmarkConfig':
                        errorMsg = "Failed to update bitmark config";
                        if(response.data.result.invalid_field != null){
                            errorMsg = errorMsg + ". Invalid field: " + response.data.result.invalid_field;
                        }
                        break;
                    default:
                        errorMsg = response.data.result;
                    };
                    deferred.reject(errorMsg);

                }
            }, function errorCallback(response) {
                // backend internal error
                if(response.status === 401) {
                    $rootScope.$broadcast('AppAuthenticated', false);
                    $location.path('/login');

                }
                // else if(response.status === -1){
                //     deferred.reject("");
                // }
                else {
                    deferred.reject(response);
                }

            });

            return deferred.promise;
        }

    };
});
