package main

import (
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"
	"os"
	"time"
)

const PaddleHeight = 4
const PaddleSymbol = 0x2588
const BallSymbol = 0x25CF

type GameObject struct {
	row, col, width, height int
	symbol                  rune
}

var screen tcell.Screen
var player1Paddle *GameObject
var player2Paddle *GameObject
var ball *GameObject
var debugLog string
var counter int

var gameObjects []*GameObject

// This program just prints "Hello, World!".  Press ESC to exit.
func main() {
	InitScreen()
	InitGameState()
	inputChan := InitUserInput()
	for {
		//counter++
		DrawState()
		time.Sleep(50 * time.Millisecond)
		key := ReadInput(inputChan)
		HandleUserInput(key)
	}
}

func InitScreen() {
	encoding.Register()
	var err error
	screen, err = tcell.NewScreen()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	if e := screen.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	defStyle := tcell.StyleDefault.
		Background(tcell.ColorBlack).
		Foreground(tcell.ColorWhite)
	screen.SetStyle(defStyle)
}

func InitUserInput() chan string {
	inputChan := make(chan string)
	go func() {
		for {
			switch ev := screen.PollEvent().(type) {
			case *tcell.EventKey:
				inputChan <- ev.Name()
			}
		}
	}()

	return inputChan
}

func HandleUserInput(key string) {
	_, screenHeight := screen.Size()

	if key == "Rune[q]" {
		screen.Fini()
		os.Exit(0)
	} else if key == "Up" && player2Paddle.row > 0 {
		player2Paddle.row--
	} else if key == "Down" && player2Paddle.row+player2Paddle.height < screenHeight {
		player2Paddle.row++
	} else if key == "Rune[w]" && player1Paddle.row > 0 {
		player1Paddle.row--
	} else if key == "Rune[s]" && player1Paddle.row+player1Paddle.height < screenHeight {
		player1Paddle.row++
	}
}

func InitGameState() {
	width, height := screen.Size()
	paddleStart := height/2 - PaddleHeight/2

	player1Paddle = &GameObject{
		row: paddleStart, col: 0, width: 1, height: PaddleHeight, symbol: PaddleSymbol,
	}

	player2Paddle = &GameObject{
		row: paddleStart, col: width - 1, width: 1, height: PaddleHeight, symbol: PaddleSymbol,
	}

	ball = &GameObject{
		row: height / 2, col: width / 2, width: 1, height: 1, symbol: BallSymbol,
	}

	gameObjects = []*GameObject{
		player1Paddle,
		player2Paddle,
		ball,
	}
}

func ReadInput(inputChan chan string) string {
	var key string

	select {
	case key = <-inputChan:
	default:
		key = ""
	}

	return key
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

func DrawState() {
	screen.Clear()
	PrintString(0, 0, debugLog)

	for _, obj := range gameObjects {
		Print(obj.col, obj.row, obj.width, obj.height, obj.symbol)
	}

	//Print(3, 5, 1, 1, rune(counter+'0'))

	//PrintString(screen, 0, 0, "Hello, World!")
	screen.Show()
}
