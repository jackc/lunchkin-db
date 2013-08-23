package main

import (
	"io"
)

func RenderFooter(writer io.Writer) (err error) {
	io.WriteString(writer, `</body>
</html>`)
	return
}
