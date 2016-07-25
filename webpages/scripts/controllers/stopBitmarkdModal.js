angular.module('bitmarkWebguiApp')
    .controller('StopBitmarkdModalCtrl', function ($scope, $uibModalInstance, type) {
        switch(type) {
        case 'config':
            $scope.message = "To config bitmarkd, we need to stop the bitmarkd first, would you like to stop the bitmarkd now ?";
            break;
        case 'switch':
            $scope.message = "To switch chain, we need to stop the bitmarkd first, would you like to stop the bitmarkd now ?";
        };

        $scope.ok = function () {
            $uibModalInstance.close(true);
        };

        $scope.cancel = function () {
            $uibModalInstance.dismiss('cancel');
        };
    });
