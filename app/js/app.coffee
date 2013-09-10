window.App =
  Collections: {}
  Models: {}
  Views: {}

class App.Models.Player extends Backbone.Model
  idAttribute: 'player_id'

class App.Models.Game extends Backbone.Model
  idAttribute: 'game_id'

class App.Collections.Standings extends Backbone.Collection
  url: "api/v1/standings"
  sortAttribute: 'rating'
  sortDescending: true
  comparator: (left, right)->
    left = left.get @sortAttribute
    right = right.get @sortAttribute

    lessValue = if @sortDescending then -1 else 1
    greaterValue = -lessValue

    if left < right
      lessValue
    else if left > right
      greaterValue
    else
      0

class App.Collections.Players extends Backbone.Collection
  model: App.Models.Player
  url: 'api/v1/players'
  comparator: (player)->
    player.get('name').toLowerCase()

class App.Collections.Games extends Backbone.Collection
  model: App.Models.Game
  url: 'api/v1/games'
  comparator: 'date'

class App.Views.PlayerStanding extends Marionette.ItemView
  tagName: 'tr'
  # template: '#playerStandingsTemplate'
  template : (serialized_model)->
    _.template($('#playerStandingsTemplate').html(), serialized_model, {variable: 'ps'})


class App.Views.Game extends Marionette.ItemView
  tagName: 'li'
  template : (serialized_model)->
    _.template($('#gameTemplate').html(), serialized_model, {variable: 'g'})

class App.Views.Games extends Marionette.CompositeView
  tagName: 'div'
  itemView: App.Views.Game
  template : (serialized_model)->
    _.template($('#gamesTemplate').html(), serialized_model, {variable: 'games'})

  appendHtml: (collectionView, itemView)->
    collectionView.$("ul").append(itemView.el)


class App.Views.Standings extends Marionette.CompositeView
  tagName: 'table'
  className: 'standings'
  template: '#standingsTemplate'
  itemView: App.Views.PlayerStanding

  events:
    'click th' : 'onHeaderClick'

  collectionEvents:
    'sort': 'render'

  onHeaderClick: (e)->
    sortAttribute = $(e.target).data("sort")

    if @collection.sortAttribute == sortAttribute
      @collection.sortDescending = !@collection.sortDescending
    else
      @collection.sortAttribute = sortAttribute
      @collection.sortDescending = true

    @collection.sort()

class App.Views.Player extends Marionette.ItemView
  tagName: 'li'
  template: '#playerTemplate'

  events:
    'click button' : 'deleteSelf'

  deleteSelf: ->
    @model.destroy dataType: 'text', wait: true

class App.Views.Players extends Marionette.CompositeView
  id: 'playersPage'
  template: '#playersTemplate'
  itemView: App.Views.Player

  appendHtml: (collectionView, itemView)->
    collectionView.$("li:last").before  (itemView.el)

  events:
    'submit form' : 'addPlayer'

  collectionEvents:
    'sort': 'render'

  addPlayer: ->
    @collection.create({name: @$('#player_name').val()}, {wait: true})


class window.Controller extends Marionette.Controller
  standings: ->
    collection = new App.Collections.Standings
    collection.fetch()
    standings = new App.Views.Standings collection: collection
    Lunchkin.mainRegion.show(standings)
  players: ->
    collection = new App.Collections.Players
    collection.fetch()
    players = new App.Views.Players collection: collection
    Lunchkin.mainRegion.show(players)
  games: ->
    collection = new App.Collections.Games
    collection.fetch()
    games = new App.Views.Games collection: collection
    Lunchkin.mainRegion.show(games)
  gamesNew: ->
    console.log 'record a game'

controller = new Controller

window.Router = new Marionette.AppRouter
  controller: controller
  appRoutes:
    "": "standings"
    "standings": "standings"
    "players": "players"
    "games": "games"
    "games/new": "gamesNew"

window.Lunchkin = new Marionette.Application

Lunchkin.addInitializer ->
  Lunchkin.addRegions
    mainRegion: '#mainRegion'

$ ->
  Lunchkin.start()
  Backbone.history.start()
