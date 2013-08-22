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

window.Lunchkin = new Marionette.Application

Lunchkin.addInitializer ->
  Lunchkin.addRegions
    mainRegion: '#mainRegion'

Lunchkin.addInitializer ->
  collection = new App.Collections.Standings
  collection.fetch()
  standings = new App.Views.Standings collection: collection
  Lunchkin.mainRegion.show(standings)

$ ->
  Lunchkin.start()
