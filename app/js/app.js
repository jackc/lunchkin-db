'use strict';


// Declare app level module which depends on filters, and services
angular.module('myApp', ['myApp.filters', 'myApp.services', 'myApp.directives', 'myApp.controllers']).
  config(['$routeProvider', function($routeProvider) {
    $routeProvider.when('/standings', {templateUrl: 'partials/standings.html', controller: 'Standings'});
    $routeProvider.when('/players', {templateUrl: 'partials/players.html', controller: 'Players'});
    $routeProvider.when('/games', {templateUrl: 'partials/games.html', controller: 'Games'});
    $routeProvider.otherwise({redirectTo: '/standings'});
  }]);
