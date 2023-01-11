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
const InitialBallVelocityRow = 1
const InitialBallVelocityCol = 2

/**

How to measure ball velocity in XY Coordinates - Raw demonstration

           (-1,-1)     (-1,0)    (-1,1)
                     \    |    /
	(0, -1) <--------- (ball) ---------> (0,1)
					/	 |    \
		   (1,-1)      (1,0)     (1,1)
*/

type GameObject struct {
	row, col, width, height int
	velRow, velCol          int // Velocity for the ball movement, Defining velocity in this struct is not ideal
	symbol                  rune
}

var screen tcell.Screen
var player1Paddle *GameObject
var player2Paddle *GameObject
var ball *GameObject
var debugLog string
var counter int
var isGamePaused bool

var gameObjects []*GameObject

// This program just prints "Hello, World!".  Press ESC to exit.
func main() {
	InitScreen()
	InitGameState()
	inputChan := InitUserInput()
	for !IsGameOver() {
		//counter++
		HandleUserInput(ReadInput(inputChan))
		UpdateState()
		DrawState()
		time.Sleep(60 * time.Millisecond)

	}

	winner := GetWinner()
	screenWidth, screenHeight := screen.Size()
	PrintStringCenter(screenWidth/2, screenHeight/2-1, "Game Over!")
	PrintStringCenter(screenWidth/2, screenHeight/2, fmt.Sprintf("%s wins...", winner))
	screen.Show()
	time.Sleep(3 * time.Second)
	screen.Fini()
}

func IsGameOver() bool {
	return GetWinner() != ""
}

func GetWinner() string {
	screenWidth, _ := screen.Size()

	if ball.col < 0 {
		return "Player 1"
	} else if ball.col >= screenWidth {
		return "Player 2"
	} else {
		return ""
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
	} else if key == "Rune[p]" {
		isGamePaused = !isGamePaused
	}
}

func InitGameState() {
	width, height := screen.Size()
	paddleStart := height/2 - PaddleHeight/2

	player1Paddle = &GameObject{
		row: paddleStart, col: 0, width: 1, height: PaddleHeight,
		velRow: 0, velCol: 0,
		symbol: PaddleSymbol,
	}

	player2Paddle = &GameObject{
		row: paddleStart, col: width - 1, width: 1, height: PaddleHeight,
		velRow: 0, velCol: 0,
		symbol: PaddleSymbol,
	}

	ball = &GameObject{
		row: height / 2, col: width / 2, width: 1, height: 1,
		velRow: InitialBallVelocityRow, velCol: InitialBallVelocityCol,
		symbol: BallSymbol,
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

func DrawState() {

	if isGamePaused {
		return
	}

	screen.Clear()
	PrintString(0, 0, debugLog)

	for _, obj := range gameObjects {
		Print(obj.col, obj.row, obj.width, obj.height, obj.symbol)
	}

	screen.Show()
}

func UpdateState() {

	if isGamePaused {
		return
	}

	for i := range gameObjects {
		gameObjects[i].row += gameObjects[i].velRow
		gameObjects[i].col += gameObjects[i].velCol
	}

	if CollidesWithWall(ball) {
		ball.velRow = -ball.velRow
	}

	if CollidesWithPaddle(ball, player1Paddle) || CollidesWithPaddle(ball, player2Paddle) {
		ball.velCol = -ball.velCol
	}
}

func CollidesWithWall(obj *GameObject) bool {
	_, screenHeight := screen.Size()
	return obj.row+obj.velRow < 0 || obj.row+obj.velRow >= screenHeight
}

func CollidesWithPaddle(ball, paddle *GameObject) bool {

	var columnCollision bool

	if ball.col < paddle.col {
		columnCollision = ball.col+ball.velCol >= paddle.col
	} else {
		columnCollision = ball.col+ball.velCol <= paddle.col
	}

	return columnCollision &&
		ball.row >= paddle.row &&
		ball.row < paddle.row+paddle.height
}
