angular.module('bitmarkWebguiApp')
    .controller('NavbarCtrl', function ($rootScope, $scope, $window, $uibModal, $log, $location, httpService) {
        $scope.$on('AppAuthenticated', function(event, value){
            $scope.showNavItem = value;
        });

        // to make the navbar item can show correctly
        $rootScope.$on('Authenticated', function(event, value){
            $rootScope.$broadcast('AppAuthenticated', value);
            if(value) {
                $scope.navbarClass = "navbar navbar-default";
            } else {
                $scope.navbarClass = "";
            }
        });

        $scope.init = function(){
            httpService.send("checkAuthenticate").then(
                function(result){
                    $scope.$emit('Authenticated', true);
                },function(){
                    $scope.$emit('Authenticated', false);
                    $location.path('/login');
            });
        };

        $scope.leftNavItems = [
            {
                url: "/main",
                active: true,
                name: "Node"
            },
            {
                url: "/console",
                active: false,
                name: "Console"
            }
        ];

        $scope.dropdownNavItems = [
            {
                url: "/chain",
                active: false,
                name: "Switch chain"
            },
            {
                url: "/network",
                active: false,
                name: "Change password"
            },
            {
                active: false,
                divider: true
            },
            {
                url: "/logout",
                active: false,
                name: "LOGOUT"
            }
        ];

        $scope.goUrl = function(navItem, type){
            switch(navItem.url){
            case "/chain":
                // check if bitmarkd is running
                httpService.send('statusBitmarkd').then(
                    function(result){
                        if(result.search("stop") >= 0) {
                            activeUrl(navItem, type);
                        }else{ // bitmarkd is running
                            var modalInstance = $uibModal.open({
                                templateUrl: 'views/stopBitmarkdModal.html',
                                controller: 'StopBitmarkdModalCtrl',
                                resolve: {
                                    type: function () {
                                        return "switch";
                                    }
                                }
                            });

                            modalInstance.result.then(function(){
                                // stop bitmarkd
                                httpService.send("stopBitmarkd").then(
                                    function(result){
                                        if(result.search("stop running bitmarkd")>=0) {
                                            activeUrl(navItem, type);
                                        }
                                    }, function(errorMsg){
                                        $log.error("stopBitmarkd error: "+errorMsg);
                                    });
                            });
                        }
                    }, function(errorMsg){
                        $log.error(errorMsg);
                    });
                break;
            case "/logout":
                httpService.send("logout").then(
                    function(){
                        $scope.$emit('Authenticated', false);
                        $location.path('/login');
                    }, function(errorMsg){
                        $scope.$emit('Authenticated', true);
                    });
                break;
            default:
                activeUrl(navItem, type);
            };
        };

        var activeUrl = function(navItem, type){
            // setup active class
            for(var i=0; i<$scope.leftNavItems.length; i++){
                var item = $scope.leftNavItems[i];
                item.active = false;
            }

            for(var i=0; i<$scope.dropdownNavItems.length; i++){
                var item = $scope.dropdownNavItems[i];
                item.active = false;
            }

            if(type == 'dropdown') {
                $scope.dropdownActive = true;
            } else {
                $scope.dropdownActive = false;
            }

            navItem.active = true;

            $location.path(navItem.url);
        };
    });
