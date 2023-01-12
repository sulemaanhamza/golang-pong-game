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
var isGamePaused bool

var gameObjects []*GameObject

func main() {

	// Initializing default states
	InitScreen()
	InitGameState()

	inputChan := InitUserInput()

	for !IsGameOver() {
		
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
