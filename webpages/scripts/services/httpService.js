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
        startBitcoind: {
            method: 'POST',
            url: hostApiPath+'/bitcoind',
            data:{
                option: "start"
            }
        },
        stopBitcoind: {
            method: 'POST',
            url: hostApiPath+'/bitcoind',
            data:{
                option: "stop"
            }
        },
        statusBitcoind: {
            method: 'POST',
            url: hostApiPath+'/bitcoind',
            data:{
                option: "status"
            }
        },
        getBitcoindInfo: {
            method: 'POST',
            url: hostApiPath+'/bitcoind',
            data:{
                option: "info"
            }
        },
        setupBitmarkd: {
            method: 'POST',
            url: hostApiPath+'/bitmarkd',
            data: {
                option: "setup"
            }
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
        },
        setupProoferd: {
            method: 'POST',
            url: hostApiPath+'/prooferd',
            data: {
                option: "setup"
            }
        },
        startProoferd: {
            method: 'POST',
            url: hostApiPath+'/prooferd',
            data:{
                option: "start"
            }
        },
        stopProoferd: {
            method: 'POST',
            url: hostApiPath+'/prooferd',
            data:{
                option: "stop"
            }
        },
        statusProoferd: {
            method: 'POST',
            url: hostApiPath+'/prooferd',
            data:{
                option: "status"
            }
        },
        startBitmarkConsole: {
            method: 'POST',
            url: hostApiPath+'/bitmarkConsole'
        }
    };

    return {
        send: function(api, data){
            var deferred = $q.defer();

            var apiConfig = angular.copy(API[api]);
            if( data != undefined) {
                if(apiConfig.data == undefined) {
                    apiConfig.data = {};
                }
                angular.extend(apiConfig.data, data);
            }


            $http(apiConfig).then(function successCallback(response) {
                if (response.data.ok) {
                    deferred.resolve(response.data.result);
                } else {
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
