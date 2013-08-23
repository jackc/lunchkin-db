package main

import (
	"html"
	"io"
	"strconv"
)

func RenderStandings(writer io.Writer, rows []map[string]interface{}, sortCol, sortDir string) (err error) {
	RenderHeader(writer)
	io.WriteString(writer, `
<div id="standingsPage">
  <table>
    <caption>Standings</caption>
    <thead>
      <th>`)
	io.WriteString(writer, sortLink("Player", "name", "asc", sortCol, sortDir))
	io.WriteString(writer, `</th>
      <th class="number">`)
	io.WriteString(writer, sortLink("Games", "num_games", "desc", sortCol, sortDir))
	io.WriteString(writer, `</th>
      <th class="number">`)
	io.WriteString(writer, sortLink("Wins", "num_wins", "desc", sortCol, sortDir))
	io.WriteString(writer, `</th>
      <th class="number">`)
	io.WriteString(writer, sortLink("Points", "num_points", "desc", sortCol, sortDir))
	io.WriteString(writer, `</th>
      <th class="number">`)
	io.WriteString(writer, sortLink("Rating", "rating", "desc", sortCol, sortDir))
	io.WriteString(writer, `</th>
    </thead>
    `)
	for _, r := range rows {
		io.WriteString(writer, `
      <tr>
        <td>`)
		io.WriteString(writer, html.EscapeString(r["name"].(string)))
		io.WriteString(writer, `</td>
        <td class="number">`)
		io.WriteString(writer, strconv.FormatInt(int64(r["num_games"].(int64)), 10))
		io.WriteString(writer, `</td>
        <td class="number">`)
		io.WriteString(writer, strconv.FormatInt(int64(r["num_wins"].(int64)), 10))
		io.WriteString(writer, `</td>
        <td class="number">`)
		io.WriteString(writer, strconv.FormatInt(int64(r["num_points"].(int64)), 10))
		io.WriteString(writer, `</td>
        <td class="number">`)
		io.WriteString(writer, html.EscapeString(r["rating"].(string)))
		io.WriteString(writer, `</td>
      </tr>
    `)
	}
	io.WriteString(writer, `
  </table>
</div>
`)
	RenderFooter(writer)
	io.WriteString(writer, `
`)
	return
}
