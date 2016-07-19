angular.module('bitmarkWebguiApp')
    .controller('NavbarCtrl', function ($rootScope, $scope, $location, httpService, configuration) {
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
                    configuration.setChain(result.chain);
                    configuration.setBitmarkCliConfigFile(result.bitmark_cli_config_file);
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
