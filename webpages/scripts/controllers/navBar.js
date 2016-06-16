angular.module('bitmarkWebguiApp')
    .controller('NavbarCtrl', function ($rootScope, $scope, $location, httpService, $cookies) {
        $scope.$on('AppAuthenticated', function(event, value){
            $scope.showNavItem = value;
        });

        // to make the navbar item can show correctly
        $rootScope.$on('Authenticated', function(event, value){
            $rootScope.$broadcast('AppAuthenticated', value);
        });

        $scope.init = function(){
            httpService.send("checkAuthenticate").then(
                function(){
                    $scope.$emit('Authenticated', true);
                    $location.path('/');

                },function(){
                    $scope.$emit('Authenticated', false);
                    $location.path('/login');
            });
        };


        $scope.logout = function(){
            httpService.send("logout").then(
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
