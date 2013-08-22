window.App =
  Collections: {}
  Models: {}
  Views: {}

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
  url: 'api/v1/players'
  comparator: 'name'

class App.Views.PlayerStanding extends Marionette.ItemView
  tagName: 'tr'
  template: '#playerStandingsTemplate'

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

class App.Views.Players extends Marionette.CompositeView
  id: 'playersPage'
  template: '#playersTemplate'
  itemView: App.Views.Player

  appendHtml: (collectionView, itemView)->
    collectionView.$("li:last").before  (itemView.el)

  events:
    'submit form' : 'addPlayer'

  addPlayer: ->
    @collection.create name: @$('#player_name').val()


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
    console.log 'games'
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
