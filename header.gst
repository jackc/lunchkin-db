package: main
func: RenderHeader
escape: html
---
<!doctype html>
<html>
<head>
  <title>Lunchkin Scoreboard</title>
  <link rel="stylesheet" href="assets/css/app.css"/>
</head>
<body>
  <ul class="menu">
    <li><a href="/standings">Standings</a></li>
    <li><a href="/players">Players</a></li>
    <li><a href="/games">Games</a></li>
    <li><a href="/games/new">Record a Game</a></li>
  </ul>
