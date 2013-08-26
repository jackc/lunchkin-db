package: main
imports: strconv
func: RenderGamesIndex
escape: html
parameters: games []Game
---
<% RenderHeader(writer) %>
<h1>Games</h1>

<ul>
  <% for _, g := range games { %>
    <li>
      <%= g.Date.Format("01/02/2006") %>
      <span> - <%=i g.Length %> rounds</span>

      <table>
        <thead>
          <tr>
            <th>Player</th>
            <th>Level</th>
            <th>Effective Level</th>
            <th>Winner</th>
          </tr>
        </thead>
        <% for _, gp := range g.Players { %>
          <tr>
            <td><%= gp.Name %></td>
            <td><%=i gp.Level %></td>
            <td><%=i gp.EffectiveLevel %></td>
            <td><%= strconv.FormatBool(gp.Winner) %></td>
          </tr>
        <% } %>
      </table>
      <form action="<%= deleteGamePath(g.GameId) %>" method="POST">
        <button>Delete</button>
      </form>
    </li>
  <% } %>
</ul>
<% RenderFooter(writer) %>
