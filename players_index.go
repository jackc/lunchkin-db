package main

import (
	"html"
	"io"
)

func RenderPlayersIndex(writer io.Writer, players []*player) (err error) {
	RenderHeader(writer)
	io.WriteString(writer, `
<div id="playersPage">

<h1>Players</h1>

<ul>
  `)
	for _, p := range players {
		io.WriteString(writer, `
    <li>
      <div class="name">`)
		io.WriteString(writer, html.EscapeString(p.name))
		io.WriteString(writer, `</div>
      <form action="`)
		io.WriteString(writer, html.EscapeString(deletePlayerPath(p.player_id)))
		io.WriteString(writer, `" method="POST">
        <button>Delete</button>
      </form>
    </li>
  `)
	}
	io.WriteString(writer, `
  <li>
    <form action="/players" method="POST">
      <div class="name"><input type="text" name="player_name" id="player_name" /></div>
      <button>Add</button>
    </form>
  </li>
</ul>

</div>
`)
	RenderFooter(writer)
	io.WriteString(writer, `
`)
	return
}
