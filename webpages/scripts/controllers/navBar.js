angular.module('bitmarkMgmtApp')
    .controller('NavbarCtrl', function ($scope, $location, httpService) {
        $scope.$on('AppAuthenticated', function(event, value){
            $scope.showNavItem = value;
        });
        httpService.send("checkAuthenticate").then(
            function(){
                $scope.showNavItem = true;
                $location.path('/');

            },function(){
                $scope.showNavItem = false;

                $location.path('/login');
            });
        $scope.logout = function(){
            httpService.send("logout").then(
                function(){
                    $scope.showNavItem = false;

                    $scope.goUrl('/login');
                }, function(errorMsg){
                    $scope.showNavItem = true;

                    $scope.errorMsg = errorMsg;
                });
        };
        $scope.goUrl = function(path){
            $location.path(path);
        };
    });
