package main

import (
	"html"
	"io"
	"strconv"
)

func RenderGamesIndex(writer io.Writer, games []Game) (err error) {
	RenderHeader(writer)
	io.WriteString(writer, `
<h1>Games</h1>

<ul>
  `)
	for _, g := range games {
		io.WriteString(writer, `
    <li>
      `)
		io.WriteString(writer, html.EscapeString(g.Date.Format("01/02/2006")))
		io.WriteString(writer, `
      <span> - `)
		io.WriteString(writer, strconv.FormatInt(int64(g.Length), 10))
		io.WriteString(writer, ` rounds</span>

      <table>
        <thead>
          <tr>
            <th>Player</th>
            <th>Level</th>
            <th>Effective Level</th>
            <th>Winner</th>
          </tr>
        </thead>
        `)
		for _, gp := range g.Players {
			io.WriteString(writer, `
          <tr>
            <td>`)
			io.WriteString(writer, html.EscapeString(gp.Name))
			io.WriteString(writer, `</td>
            <td>`)
			io.WriteString(writer, strconv.FormatInt(int64(gp.Level), 10))
			io.WriteString(writer, `</td>
            <td>`)
			io.WriteString(writer, strconv.FormatInt(int64(gp.EffectiveLevel), 10))
			io.WriteString(writer, `</td>
            <td>`)
			io.WriteString(writer, html.EscapeString(strconv.FormatBool(gp.Winner)))
			io.WriteString(writer, `</td>
          </tr>
        `)
		}
		io.WriteString(writer, `
      </table>
      <form action="`)
		io.WriteString(writer, html.EscapeString(deleteGamePath(g.GameId)))
		io.WriteString(writer, `" method="POST">
        <button>Delete</button>
      </form>
    </li>
  `)
	}
	io.WriteString(writer, `
</ul>
`)
	RenderFooter(writer)
	io.WriteString(writer, `
`)
	return
}
