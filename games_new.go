package main

import (
	"html"
	"io"
	"strconv"
)

func RenderGamesNew(writer io.Writer, players []Player) (err error) {
	RenderHeader(writer)
	io.WriteString(writer, `
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
    `)
	for _, p := range players {
		io.WriteString(writer, `
      <tr>
        <td><input type="checkbox" id="PlayerId`)
		io.WriteString(writer, strconv.FormatInt(int64(p.PlayerId), 10))
		io.WriteString(writer, `" name="Game.Players.Ids" value="`)
		io.WriteString(writer, strconv.FormatInt(int64(p.PlayerId), 10))
		io.WriteString(writer, `" /></td>
        <td><label for="PlayerId`)
		io.WriteString(writer, strconv.FormatInt(int64(p.PlayerId), 10))
		io.WriteString(writer, `">`)
		io.WriteString(writer, html.EscapeString(p.Name))
		io.WriteString(writer, `</label></td>
        <td><input type="number" min="1" max="20" name="Game.Players.`)
		io.WriteString(writer, strconv.FormatInt(int64(p.PlayerId), 10))
		io.WriteString(writer, `.Level" /></td>
        <td><input type="number" min="-1000" max="1000" name="Game.Players.`)
		io.WriteString(writer, strconv.FormatInt(int64(p.PlayerId), 10))
		io.WriteString(writer, `.EffectiveLevel" /></td>
        <td><input type="checkbox" name="Game.Players.`)
		io.WriteString(writer, strconv.FormatInt(int64(p.PlayerId), 10))
		io.WriteString(writer, `.Winner" /></td>
      </tr>
    `)
	}
	io.WriteString(writer, `
  </table>
  <button>Save</button>
</form>
`)
	RenderFooter(writer)
	io.WriteString(writer, `
`)
	return
}
