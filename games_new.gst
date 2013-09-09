package: main
func: RenderGamesNew
escape: html
parameters: players []Player
---
<% RenderHeader(writer) %>
<h1>Record a Game</h1>

<form action="/games" method="POST">
  <label for="game_date">Date</label>
  <input type="date" name="Game.Date" id="game_date" /><br />

  <label for="game_length">Length</label>
  <input type="number" min="1" max="32767" name="Game.Length" id="game_length" /><br />
  <table>
    <thead>
      <tr>
        <th></th>
        <th>Player</th>
        <th>Level</th>
        <th>Effective Level</th>
        <th>Winner</th>
      </tr>
    </thead>
    <% for _, p := range players { %>
      <tr>
        <td><input type="checkbox" id="PlayerId<%=i p.PlayerId %>" name="Game.Players.Ids" value="<%=i p.PlayerId %>" /></td>
        <td><label for="PlayerId<%=i p.PlayerId %>"><%= p.Name %></label></td>
        <td><input type="number" min="1" max="20" name="Game.Players.<%=i p.PlayerId %>.Level" /></td>
        <td><input type="number" min="-1000" max="1000" name="Game.Players.<%=i p.PlayerId %>.EffectiveLevel" /></td>
        <td><input type="checkbox" name="Game.Players.<%=i p.PlayerId %>.Winner" /></td>
      </tr>
    <% } %>
  </table>
  <button>Save</button>
</form>
<% RenderFooter(writer) %>
