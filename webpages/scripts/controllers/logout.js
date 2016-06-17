angular.module('bitmarkWebguiApp')
    .controller('LogoutCtrl', function ($rootScope, $scope, $location, httpService, $cookies, configuration, BitmarkPayConfig) {

        // var bitmarkPayConfigFile = "/home/yuntai/testWebgui/config/bitmark-pay/bitmark-pay-LOCAL.xml";
        var bitmarkPayConfigFile = BitmarkPayConfig[configuration.getChain()];

        $scope.logout = function(){
            httpService.send("logout", {
                password: $scope.password,
	        bitmark_pay_config_file: bitmarkPayConfigFile
            }).then(
                function(){
                    $scope.$emit('Authenticated', false);
                    $cookies.remove('bitmark-chain');
                    $scope.goUrl('/login');
                }, function(errorMsg){
                    $scope.$emit('Authenticated', true);
                    $scope.errorMsg = errorMsg;
                });
        };
        $scope.goUrl = function(path){
            $location.path(path);
        };
    });
