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
                    if(result.bitmark_cli_config_file.length != 0 )  {
                        $scope.$emit('Authenticated', true);
                    }
                },function(){
                    $scope.$emit('Authenticated', false);
                    $location.path('/login');
            });
        };

        $(".nav a").on("click", function(){
         $(".nav").find(".active").removeClass("active");
          $(this).parent().addClass("active");
        });
        
        $scope.goUrl = function(path){
            $location.path(path);
        };
    });
