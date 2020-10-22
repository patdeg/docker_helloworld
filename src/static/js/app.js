/*jslint for:false */

$(document).ready(function() {
  "use strict";
  $('.autoclose').click(function(event) {
    $('.navbar-collapse').collapse('hide');
  });
});

function debug(...args) {
  if (DEBUG) {
    console.log(...args);
  }   
}

function error(...args) {
  console.error(...args);  
}

const myApp = angular.module('myApp', ['ngRoute', 'ui.bootstrap']).
config(['$routeProvider', '$locationProvider',
  function($routeProvider, $locationProvider) {
    "use strict";

    $locationProvider.html5Mode(true);
    $routeProvider
    .when('/', {
      templateUrl: 'static/html/home.html',
      controller:'HomeController'
    })
    .when('/about', {
      templateUrl: 'static/html/about.html',
      controller: 'AboutController'
    })    
    .otherwise({
      redirectTo: '/'
    });
}]);

myApp.controller('BodyController', ['$scope', '$location',
  function($scope, $location) {
    "use strict";

    $scope.world = 'World';

    $scope.go = function(path) {
      debug("Going to ",path);
      $location.path(path);
    };

    $scope.isActive = function(page) {      
      $scope.location = $location.path();            
      if ($scope.location) {
        return $scope.location == page;        
      }      
      return false;
    };
  }
]);

myApp.controller('HomeController', ['$scope', '$http',
  function($scope, $http) {
    "use strict";

    $scope.world = 'World';

    $scope.getList = function() {
      debug(">>>> getList");
      var url = "/api/list";
      debug("Calling ",url);      
      $scope.is_loading = true;
      $http.get(url)
      .success(function(data) {
        $scope.is_loading = false;
        $scope.data = data;            
        debug("List:",$scope.data);            
      })
      .error(function(errorMessage, errorCode, errorThrown) {
        $scope.is_loading = false;        
        error('Error - errorMessage, errorCode, errorThrown:',
            errorMessage, errorCode, errorThrown);
        alert(errorMessage);
      });
    };
    $scope.getList();

}]);


myApp.controller('AboutController', ['$scope', 
  function($scope) {
    "use strict";

    $scope.world = 'World';

}]);