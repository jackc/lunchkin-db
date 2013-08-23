package: main
imports: fmt
func: RenderStandings
escape: html
parameters: rows []map[string]interface{}, sortCol, sortDir string
---
<% RenderHeader(writer) %>
<div id="standingsPage">
  <table>
    <caption>Standings</caption>
    <thead>
      <th><% io.WriteString(writer, sortLink("Player", "name", "asc", sortCol, sortDir)) %></th>
      <th class="number"><% io.WriteString(writer, sortLink("Games", "num_games", "desc", sortCol, sortDir)) %></th>
      <th class="number"><% io.WriteString(writer, sortLink("Wins", "num_wins", "desc", sortCol, sortDir)) %></th>
      <th class="number"><% io.WriteString(writer, sortLink("Points", "num_points", "desc", sortCol, sortDir)) %></th>
      <th class="number"><% io.WriteString(writer, sortLink("Rating", "rating", "desc", sortCol, sortDir)) %></th>
    </thead>
    <% for _, r := range rows { %>
      <tr>
        <td><%= r["name"].(string) %></td>
        <td class="number"><%=i r["num_games"].(int64) %></td>
        <td class="number"><%=i r["num_wins"].(int64) %></td>
        <td class="number"><%=i r["num_points"].(int64) %></td>
        <td class="number"><%= r["rating"].(string) %></td>
      </tr>
    <% } %>
  </table>
</div>
<% RenderFooter(writer) %>
