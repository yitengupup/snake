// snake project main.go
package main

import (
	"github.com/nsf/termbox-go"
	"math/rand"
	"time"
)

// var of map
const hight = 20
const wide = 20

var (
	left   int
	top    int
	bottom int
)

// type of snake
var (
	area        = [hight][wide]rune{}
	direction   rune // U,L,D,R   means Up,Left,Down,Right
	status      bool //game is not over
	head        location
	tail        location
	length      int
	isfoodExist bool //food exists
)

//var of difficulty
const (
	easy      = 0 + iota //1000 ms --> easy
	normal               //800 ms --> normal
	hard                 //600 ms --> hard
	maddening            //400 ms --> hard
)

var frequncy = [4]int{800, 600, 300, 100} //声明并初始化一个长度为4的int数组

func print_wide(x, y int, char rune) {
	red := true
	c := termbox.ColorDefault
	if red {
		c = termbox.ColorRed
	}
	termbox.SetCell(x, y, char, termbox.ColorDefault, c)
}

//type of location
type location struct {
	row    int
	colume int
}

// create food location
func foodCreate() location {
	var foodLoca location // food will appear randomly in the map(can not on the snake)
	for {
		randomInt := rand.Int() % (hight * wide)
		h := randomInt / hight
		w := randomInt % wide
		if area[h][w] == 0 {
			foodLoca.row = h
			foodLoca.colume = w
			area[h][w] = '$'
			break
		}
	}
	//fmt.Fprintf(os.Stderr, "food is [%d %d]$\\n", foodLoca.row, foodLoca.colume)
	redraw(foodLoca, left, top, '$')
	return foodLoca
}

const (
	defaultColor = termbox.ColorDefault
	bgColor      = termbox.ColorDefault
	snakeColor   = termbox.ColorGreen
)

// main function for snake
func main() {
	err := termbox.Init()
	defer termbox.Close()
	if err != nil {
		panic(err)
	}
	//initial map
	initialMap()
	//go to check input
	go controlDirection()
	//main loop
	var directionBefore rune = 0
	for {
		time.Sleep(time.Millisecond * time.Duration(frequncy[hard]))

		// create food if need
		if !isfoodExist {
			foodCreate() //area[][] := '$'
			isfoodExist = true
		}
		// let the tail knows what is the direction
		area[head.row][head.colume] = direction
		tmpdirection := direction
		// move head
		switch direction {
		case 'U':
			if directionBefore == 'D' {
				area[head.row][head.colume] = 'D'
				tmpdirection = 'D'
				head.row++
			} else {
				head.row--
			}
		case 'L':
			if directionBefore == 'R' {
				area[head.row][head.colume] = 'R'
				tmpdirection = 'R'
				head.colume++
			} else {
				head.colume--
			}
		case 'R':
			if directionBefore == 'L' {
				area[head.row][head.colume] = 'L'
				tmpdirection = 'L'
				head.colume--
			} else {
				head.colume++
			}
		case 'D':
			if directionBefore == 'U' {
				area[head.row][head.colume] = 'U'
				tmpdirection = 'U'
				head.row--
			} else {
				head.row++
			}
		}
		directionBefore = tmpdirection
		// judge game state
		if head.row < 0 || head.row >= hight || head.colume < 0 || head.colume >= wide {
			//fmt.Fprintf(os.Stderr, "reach the edge\n")
			break
		}
		// determine next step
		state := area[head.row][head.colume]
		if state == 0 {
			redraw(head, left, top, '*')
			redraw(tail, left, top, ' ')
			tRow := tail.row
			tColume := tail.colume
			switch area[tail.row][tail.colume] {
			case 'U':
				tail.row--
			case 'L':
				tail.colume--
			case 'R':
				tail.colume++
			case 'D':
				tail.row++
			}
			area[tRow][tColume] = 0
		} else if state == '$' {
			redraw(head, left, top, '*')
			length++
			isfoodExist = false
			//fmt.Fprintf(os.Stderr, "eat food, perfect! \n")
		} else {
			//fmt.Fprintf(os.Stderr, "eat itself!!!\n")
			break
		}

		//fmt.Fprintf(os.Stderr, "head is [%d %d]! \n", head.row, head.colume)
	}
}

// initial draw map
func initialMap() error {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	head, tail = location{4, 4}, location{4, 4}
	direction = 'R'
	area[4][4] = 'R'
	isfoodExist = false
	rand.Seed(int64(time.Now().Unix()))
	length = 1
	var (
		w, h = termbox.Size()
		midY = h / 2
	)
	left = (w - wide) / 2
	//right  = (w + wide) / 2
	top = midY - (hight / 2)
	bottom = midY + (hight / 2) + 1

	InitialArena(top, bottom, left)
	//	renderSnake(left, bottom, g.arena.snake)
	//	renderFood(left, bottom, g.arena.food)
	//	renderScore(left, bottom, g.score)
	//	renderQuitMessage(right, bottom)

	return termbox.Flush()
}

func InitialArena(top, bottom, left int) {
	for i := top; i < bottom; i++ {
		termbox.SetCell(left-1, i, '│', defaultColor, bgColor)
		termbox.SetCell(left+wide, i, '│', defaultColor, bgColor)
	}

	termbox.SetCell(left-1, top, '┌', defaultColor, bgColor)
	termbox.SetCell(left-1, bottom, '└', defaultColor, bgColor)
	termbox.SetCell(left+wide, top, '┐', defaultColor, bgColor)
	termbox.SetCell(left+wide, bottom, '┘', defaultColor, bgColor)

	fill(left, top, wide, 1, termbox.Cell{Ch: '─'})
	fill(left, bottom, wide, 1, termbox.Cell{Ch: '─'})

}

func fill(x, y, w, h int, cell termbox.Cell) {
	for ly := 0; ly < h; ly++ {
		for lx := 0; lx < w; lx++ {
			termbox.SetCell(x+lx, y+ly, cell.Ch, cell.Fg, cell.Bg)
		}
	}
}

func redraw(point location, leftPos, topPos int, char rune) {
	termbox.SetCell(leftPos+point.colume, topPos+point.row, char, defaultColor, bgColor)
	termbox.Flush()
	//fmt.Fprintf(os.Stderr, "redraw [x=%d y=%d] with %c\n", point.colume, point.row, char)
}

// control the direction of snake
func controlDirection() {
LOOP:
	for { // 函数只写入lead，外部只读取lead，无需设锁
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyArrowUp:
				direction = 'U'
			case termbox.KeyArrowDown:
				direction = 'D'
			case termbox.KeyArrowLeft:
				direction = 'L'
			case termbox.KeyArrowRight:
				direction = 'R'
			default:
				break LOOP
			}
		}
	}
}
