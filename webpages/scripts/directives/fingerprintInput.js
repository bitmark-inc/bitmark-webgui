// Copyright (c) 2014-2016 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

'use strict';

/**
 * @ngdoc function
 * @name bitmarkWebguiApp.directive:fingerprintInput
 * @description
 * Controller of the bitmarkWebguiApp
 */
angular.module('bitmarkWebguiApp')
.directive('fingerprintInput', function() {
  return {
    restrict: 'EA',
    replace: true,
    transclude: true,
    scope: {
      file: '=file',
      fingerprint: '=fingerprint'
    },
    template: '<input type="file" name="file" class="browse" required/>',
    link: function( scope, elm, attrs ) {
      var genFingerprint = function(file){

        var reader = new FileReader();

        // Read in the image file as a data URL.
        reader.readAsArrayBuffer(file);

      };
      // Check for the various File API support.
      if (window.File && window.FileReader && window.FileList && window.Blob) {
        // Great success! All the File APIs are supported.
      } else {
        alert('The File APIs are not fully supported in this browser.');
      }

      elm.bind('change', function( evt ) {
        scope.$apply(function() {
          scope.file = evt.target.files[0];
          var reader = new FileReader();
          // Read in the image file as a data URL.
          reader.readAsArrayBuffer(scope.file);
          reader.onloadend = function(evt) {
            if (evt.target.readyState == FileReader.DONE) {
              var shaObj = new jsSHA("SHA-512", "ARRAYBUFFER");
              shaObj.update(reader.result);

              scope.$apply(function(){
                scope.fingerprint = shaObj.getHash("HEX");

              });
            }
          }
        });
      });
      scope.$on("$destroy", function(){
        ele.remove();
      });
    }
  }
});
