package: main
func: RenderPlayersIndex
escape: html
parameters: players []*player
---
<% RenderHeader(writer) %>
<div id="playersPage">

<h1>Players</h1>

<ul>
  <% for _, p := range players { %>
    <li>
      <div class="name"><%= p.name %></div>
      <form action="<%= deletePlayerPath(p.player_id) %>" method="POST">
        <button>Delete</button>
      </form>
    </li>
  <% } %>
  <li>
    <form action="/players" method="POST">
      <div class="name"><input type="text" name="player_name" id="player_name" /></div>
      <button>Add</button>
    </form>
  </li>
</ul>

</div>
<% RenderFooter(writer) %>
