angular.module('bitmarkWebguiApp')
    .controller('LogoutCtrl', function ($rootScope, $scope, $location, httpService, configuration, BitmarkPayConfig, $interval) {
        var chain = configuration.getConfiguration().chain;
        var bitmarkPayConfigFile = BitmarkPayConfig[chain];

        $scope.error = {
            show: false,
            msg: ""
        };

        $scope.msg = [];

        $scope.decryptAlert = {
            show: false,
            msg: ""
        };

        var logout = function(){
            httpService.send("logout").then(
                function(){
                    $scope.$emit('Authenticated', false);
                    $scope.goUrl('/login');
                }, function(errorMsg){
                    $scope.$emit('Authenticated', true);
                    $scope.error.msg = errorMsg;
                    $scope.error.show = true;
                });
        };

        var decryptPromise;
        var decryptWaitingTime = 60; // 60s
        var pollDecryptCount = 0;
        var decryptJobHash;

        var decryptWalletAndLogout = function(keyPair){
            // decrypt bitmark-wallet
            var net = chain;
            if(net == "local") {
                net = "local_bitcoin_reg";
            }

            httpService.send("decryptBitmarkPay", {
                config: bitmarkPayConfigFile,
                net: net,
                password: keyPair.private_key
            }).then(function(decryptPayJobHash){
                // check decrypt jobs success
                $interval.cancel(decryptPromise);
                decryptJobHash = decryptPayJobHash;
                decryptPromise = $interval(function(){
                    httpService.send("getBitmarkPayStatus", {
                        job_hash: decryptPayJobHash
                    }).then(function(payStatusResult){
                        switch(payStatusResult){
                        case "success":
                            $interval.cancel(decryptPromise);
                            pollDecryptCount = 0;
                            logout();
                        case "running":
                            pollDecryptCount++;
                            if(pollDecryptCount*3 > decryptWaitingTime){
                                $scope.decryptAlert.msg = "Bitmark-pay has been running for "
                                    +pollDecryptCount*3+
                                    " seconds, normally it could cost 7 mins, would you want to stop the process?";
                                $scope.decryptAlert.show = true;
                            }
                            break;
                        case "fail":
                            $interval.cancel(decryptPromise);
                            $scope.error.show = true;
                            $scope.error.msg = "failed to decrypt wallet, please check your bitcoin status";
                            break;
                        case "stopped":
                            $interval.cancel(decryptPromise);
                            $scope.error.show = true;
                            $scope.error.msg = "wallet was decrypted before, please decrypt your bitmark wallet first";
                            break;
                        }
                    });
                }, 3*1000);

                // stop bitcoind
                if(chain == "local") {
                    httpService("stopBitcoind").then(function(stop){
                        logout();
                    }, function(stopErr){
                        $scope.error.msg = "failed to stop: "+ stopErr;
                        $scope.error.show = true;
                    });
                } else {
                    logout();
                }
            }, function(decryptErr){
                $scope.error.msg = "failed to decrypt wallet: "+decryptErr;
                $scope.error.show = true;
            });
        };


        $scope.logout = function(){
	    // get privateKey from bitmark-cli
            $scope.msg.push("get keypair...");
            httpService.send("getBitmarkKeyPair", {
                password: $scope.password
            }).then(function(keypair){
                // check bitcoind status
                $scope.msg.push("checking bitcoind...");
                if(chain == "local"){
                    httpService.send("statusBitcoind").then(function(bitcoinStatus){
                        if(bitcoinStatus == "stopped"){ // start bitcoind for the user
                            $scope.msg.push("bitcoind is stopped, try to start it");
                            httpService.send("startBitcoind").then(function(startSuccess){
                                $scope.msg.push("bitcoind is started");
                                decryptWalletAndLogout(keypair);
                            }, function(startErr){
                                $scope.error.msg = "failed to start bitcoind: "+startErr;
                                $scope.error.show = true;
                            });
                        }else{
                            $scope.msg.push("bitcoind is started...");
                            decryptWalletAndLogout(keypair);
                        }
                });
                } else{
                    decryptWalletAndLogout(keypair);
                }
            }, function(keypairErr){
                $scope.error.msg = keypairErr;
                $scope.error.show = true;
            });
        };

        var killPromise;
        var killBitmarkPayStatusProcess = function(jobHash, alertObj){
            httpService.send('stopBitmarkPayProcess', {"job_hash": jobHash}).then(function(result){
                $interval.cancel(killPromise);
                killPromise = $interval(function(){
                    httpService.send("getBitmarkPayStatus", {
                        job_hash: jobHash
                    }).then(function(payStatus){
                        if(payStatus == "stopped"){
                            $interval.cancel(killPromise);
                            alertObj.show = false;
                        }
                    });
                }, 3*1000);
            }, function(err){
                alertObj.show = true;
                alertObj.msg = err;
            });
        };

        $scope.killPayProcess = function(kill){
            if(kill){
                $interval.cancel(decryptPromise);
                pollDecryptCount = 0;
                if(decryptJobHash == "" || decryptJobHash == null) {
                    httpService.send('getBitmarkPayJob').then(function(jobHash){
                        decryptJobHash = jobHash;
                        killBitmarkPayStatusProcess(decryptJobHash, $scope.decryptAlert);
                    });
                }else{
                    killBitmarkPayStatusProcess(decryptJobHash, $scope.decryptAlert);
                }
            }else{
                $scope.decryptAlert.show = false;
                pollDecryptCount = 0;
            }
        };

        $scope.goUrl = function(path){
            $location.path(path);
        };

        $scope.$on("$destroy", function(){
            $interval.cancel(decryptPromise);
            $interval.cancel(killPromise);
        });
    });
