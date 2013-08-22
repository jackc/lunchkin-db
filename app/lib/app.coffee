window.App = {}
App.Models = {}
class App.Models.Todo extends Backbone.Model
  initialize: ->
    console.log 'This model has been initialized'
    @listenTo this, 'change', ->
      console.log 'Something has changed'

todo = new App.Models.Todo
todo.set 'foo', 'bar'
