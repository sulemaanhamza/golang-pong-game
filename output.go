package main

import "github.com/gdamore/tcell"

func PrintStringCenter(col, row int, str string) {
	col = col - len(str)/2
	PrintString(col, row, str)
}

func PrintString(col, row int, str string) {
	for _, c := range str {
		screen.SetContent(col, row, c, nil, tcell.StyleDefault)
		col += 1
	}
}

func Print(col, row, width, height int, ch rune) {
	for r := 0; r < height; r++ {
		for c := 0; c < width; c++ {
			screen.SetContent(col+c, row+r, ch, nil, tcell.StyleDefault)
		}
	}
}
