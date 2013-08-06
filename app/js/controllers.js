'use strict';

/* Controllers */

angular.module('myApp.controllers', []).
  controller('MyCtrl1', ['$scope', function($scope) {
    $scope.name = "Jack";
    $scope.noteHeaders = [
      {title: "Note 1", body: "Bla bla bla..."},
      {title: "Note 2", body: "Yada yada yada..."},
      {title: "Note 3", body: "Fe Fi Fo Fum..."},
      {title: "Note 4", body: "Out of things to say..."}
    ];
    $scope.currentNote = $scope.noteHeaders[0];

    $scope.selectNote = function(noteHeader) {
      $scope.currentNote = noteHeader;
    }

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

    $http.get('api/v1/players').success(function(data) {
      $scope.players = data;
    });

    $scope.game = {players: []};

    $scope.addPlayerToGame = function(player) {
      $scope.game.players.push({player_id: player.player_id, name: player.name});
    };

    $scope.createGame = function(game) {
      $http.post('api/v1/games', game).success(function(data) {
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
  }]);