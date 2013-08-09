'use strict';

/* Controllers */

angular.module('myApp.controllers', []).
  controller('Standings', ['$scope', '$http', function($scope, $http) {
    $http.get('api/v1/standings').success(function(data) {
      $scope.players = data;
    });
    $scope.sort = '-rating';
  }])
  .controller('Players', ['$scope', '$http', function($scope, $http) {
    $http.get('api/v1/players').success(function(data) {
      $scope.players = data;
    });

    $scope.create = function(player) {
      $http.post('api/v1/players', player).success(function(data) {
        player.name = "";
        $http.get('api/v1/players').success(function(data) {
          $scope.players = data;
        });
      })
    };

    $scope.deletePlayer = function(player) {
      $http.delete('api/v1/players/' + player.player_id).success(function(data) {
        $http.get('api/v1/players').success(function(data) {
          $scope.players = data;
        });
      })
    };
  }])
  .controller('Games', ['$scope', '$http', function($scope, $http) {
    $http.get('api/v1/games').success(function(data) {
      $scope.games = data;
    });

    $scope.game = {
      date: moment().format('YYYY-MM-DD'),
      players: []
    };

    $http.get('api/v1/players').success(function(data) {
      $scope.game.players = data;
    });

    $scope.createGame = function(game) {
      var postData = {
        date: game.date,
        players: []
      };

      for(var i = 0; i < game.players.length; i++) {
        var p = game.players[i];
        if(p.played) {
          postData.players.push({
            player_id: p.player_id,
            level: p.level,
            effective_level: p.effective_level,
            winner: p.winner
          });
        }
      }

      $http.post('api/v1/games', postData).success(function(data) {
        $http.get('api/v1/games').success(function(data) {
          $scope.games = data;
        });
      })
    };

    $scope.deleteGame = function(game) {
      $http.delete('api/v1/games/' + game.game_id).success(function(data) {
        $http.get('api/v1/games').success(function(data) {
          $scope.games = data;
        });
      })
    };
  }])
  .controller('NewGame', ['$scope', '$http', '$location', function($scope, $http, $location) {
    $scope.game = {
      date: moment().format('YYYY-MM-DD'),
      players: []
    };

    $http.get('api/v1/players').success(function(data) {
      $scope.game.players = data;
    });

    $scope.createGame = function(game) {
      var postData = {
        date: game.date,
        length: game.length,
        players: []
      };

      for(var i = 0; i < game.players.length; i++) {
        var p = game.players[i];
        if(p.played) {
          postData.players.push({
            player_id: p.player_id,
            level: p.level,
            effective_level: p.effective_level,
            winner: p.winner
          });
        }
      }

      $http.post('api/v1/games', postData).success(function(data) {
        $location.path("/standings")
      })
    };
  }]);