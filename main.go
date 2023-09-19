package main

// Welcome to
// __________         __    __  .__                               __
// \______   \_____ _/  |__/  |_|  |   ____   ______ ____ _____  |  | __ ____
//  |    |  _/\__  \\   __\   __\  | _/ __ \ /  ___//    \\__  \ |  |/ // __ \
//  |    |   \ / __ \|  |  |  | |  |_\  ___/ \___ \|   |  \/ __ \|    <\  ___/
//  |________/(______/__|  |__| |____/\_____>______>___|__(______/__|__\\_____>
//
// This file can be a nice home for your Battlesnake logic and helper functions.
//
// To get you started we've included code to prevent your Battlesnake from moving backwards.
// For more info see docs.battlesnake.com

import (
	"log"
	// "math"
	"math/rand"
	"sort"
	// TODO: add this:
	// "gorgonia.org/gorgonia"
	// "strings"
)

// info is called when you create your Battlesnake on play.battlesnake.com
// and controls your Battlesnake's appearance
// TIP: If you open your Battlesnake URL in a browser you should see this data
func info() BattlesnakeInfoResponse {
	log.Println("INFO")

	return BattlesnakeInfoResponse{
		APIVersion: "1",
		Author:     "username-pass", // TODO: Your Battlesnake username
		Color:      "#888888",       // TODO: Choose color
		Head:       "default",       // TODO: Choose head
		Tail:       "default",       // TODO: Choose tail
	}
}

// start is called when your Battlesnake begins a game
func start(state GameState) {
	log.Println("GAME START")
}

// end is called when your Battlesnake finishes a game
func end(state GameState) {
	log.Printf("GAME OVER\n\n")
}

// move is called on every turn and returns your next move
// Valid moves are "up", "down", "left", or "right"
// See https://docs.battlesnake.com/api/example-move for available data
func move(state GameState) BattlesnakeMoveResponse {

	isMoveSafe := map[string]bool{
		"up":    true,
		"down":  true,
		"left":  true,
		"right": true,
	}

	// We've included code to prevent your Battlesnake from moving backwards
	myHead := state.You.Body[0] // Coordinates of your head
	myNeck := state.You.Body[1] // Coordinates of your "neck"

	if myNeck.X < myHead.X { // Neck is left of head, don't move left
		isMoveSafe["left"] = false

	} else if myNeck.X > myHead.X { // Neck is right of head, don't move right
		isMoveSafe["right"] = false

	}
	if myNeck.Y < myHead.Y { // Neck is below head, don't move down
		isMoveSafe["down"] = false

	} else if myNeck.Y > myHead.Y { // Neck is above head, don't move up
		isMoveSafe["up"] = false
	}

	log.Print("Width: ", state.Board.Width)
	log.Print("Height: ", state.Board.Height)
	log.Print("HeadX: ", myHead.X)
	log.Print("NeckX: ", myNeck.X)
	log.Print("HeadY: ", myHead.Y)
	log.Print("NeckY: ", myNeck.Y)

	// TODO: Step 1 - Prevent your Battlesnake from moving out of bounds
	boardWidth := state.Board.Width
	boardHeight := state.Board.Height

	if myHead.X <= 0 {
		isMoveSafe["left"] = false

	} else if myHead.X >= boardWidth-1 {
		isMoveSafe["right"] = false

	}
	if myHead.Y <= 0 {
		isMoveSafe["down"] = false

	} else if myHead.Y >= boardHeight-1 {
		isMoveSafe["up"] = false
	}

	// TODO: Step 2 - Prevent your Battlesnake from colliding with itself
	mybody := state.You.Body
	log.Print("Body", mybody)

	//check against yourself
	isMoveSafe = checkBody(isMoveSafe, mybody, myHead)
	//check against all other snakes

	// TODO: Step 3 - Prevent your Battlesnake from colliding with other Battlesnakes
	opponents := state.Board.Snakes
	for s := 0; s < len(opponents); s++ {
		isMoveSafe = checkBody(isMoveSafe, opponents[s].Body, myHead)
	}

	hazards := state.Board.Hazards
	isMoveSafe = checkBody(isMoveSafe, hazards, myHead)

	// Are there any safe moves left?
	safeMoves := []string{}
	for move, isSafe := range isMoveSafe {
		if isSafe {
			safeMoves = append(safeMoves, move)
		}
	}

	//TODO: If there are no safe moves left, and there is another snake's head within reach, move there

	if len(safeMoves) == 0 {
		log.Printf("MOVE %d: No safe moves detected! Moving down\n", state.Turn)
		return BattlesnakeMoveResponse{Move: "down"}
	}

	// TODO: Step 4 - Move towards food instead of random, to regain health and survive longer
	food := state.Board.Food

	log.Print(food)
	food = nearSort(food, myHead)

	dir := []string{"left", "down"}
	nearestFood := food[0]
	dirdist := int2Coord(distboth(nearestFood, myHead))
	log.Print("Distance to nearest food: ", dirdist)

	if dirdist.X > 0 {
		dir[0] = "right"
	} else if dirdist.X == 0 {
		dir[0] = "straight"
	}
	if dirdist.Y > 0 {
		dir[1] = "up"
	} else if dirdist.Y == 0 {
		dir[1] = "straight"
	}
	// Choose a random move from the safe ones
	nextMove := "up"

	var Board = make([][]float32, boardWidth, boardHeight)
	Board = generateBoard(Board, state)
	log.Print(Board)
	if absInt(dirdist.X) < absInt(dirdist.Y) && /* absInt(dirdist.X) > 0 && */ isMoveSafe[dir[0]] {
		nextMove = dir[0]
	}
	if absInt(dirdist.Y) < absInt(dirdist.X) && /* absInt(dirdist.Y) > 0 && */ isMoveSafe[dir[1]] {
		nextMove = dir[1]
	}

	if !isMoveSafe[nextMove] {
		// Choose a random move from the safe ones
		nextMove = safeMoves[rand.Intn(len(safeMoves))]
	}

	log.Printf("MOVE %d: %s\n", state.Turn, nextMove)
	return BattlesnakeMoveResponse{Move: nextMove}
}

func absInt(x int) int {
	return absDiffInt(x, 0)
}

func absDiffInt(x, y int) int {
	if x < y {
		return y - x
	}
	return x - y
}

func absDiffUint(x, y uint) uint {
	if x < y {
		return y - x
	}
	return x - y
}

func int2Coord(numarr []int) Coord {
	tempA := numarr[0]
	tempB := numarr[1]
	return Coord{tempA, tempB}
}

func dist(c1, c2 Coord) int {
	var dx int = c1.X - c2.X
	var dy int = c1.Y - c2.Y
	return dx + dy
}

func distboth(c1, c2 Coord) []int {
	var dx int = c1.X - c2.X
	var dy int = c1.Y - c2.Y
	return []int{dx, dy}
}

func nearSort(coords []Coord, head Coord) []Coord {
	sort.Slice(coords, func(i, j int) bool {
		return dist(coords[i], head) < dist(coords[j], head)
	})
	return coords
}

func lerp(A, B, C float32) float32 {
	return A*(1-C) + B*(C)
}

func generateBoard(Board [][]float32, state GameState) [][]float32 {

	var head Coord = state.You.Head

	//set default values on Board
	for x := 0; x < state.Board.Width; x++ {
		for y := 0; y < state.Board.Width; y++ {
			Board[x][y] = 0
		}
	}

	//iterate through all of your character's positions:
	var bodyCostMax float32 = -50
	var bodyCostMin float32 = -3
	for i, coord := range state.You.Body {
		var cost float32 = lerp(bodyCostMin, bodyCostMax, float32(i/len(state.You.Body)))
		Board[coord.X][coord.Y] = cost
	}

	//iterate through all the food on the board
	var foodCostMax float32 = 10
	var foodCostMin float32 = 1
	var food []Coord = nearSort(state.Board.Food, head)
	for i, coord := range food {
		//cost based on distance
		var cost float32 = lerp(foodCostMin, foodCostMax, float32(i/len(food)))
		Board[coord.X][coord.Y] = cost
	}

	return Board
}

func checkBody(issafe map[string]bool, body []Coord, head Coord) map[string]bool {
	for i := 0; i < len(body); i++ {
		bodycell := body[i]

		if (head.X-1) == bodycell.X && head.Y == bodycell.Y {
			issafe["left"] = false
		}
		if (head.X+1) == bodycell.X && head.Y == bodycell.Y {
			issafe["right"] = false
		}
		if (head.Y-1) == bodycell.Y && head.X == bodycell.X {
			issafe["down"] = false
		}
		if (head.Y+1) == bodycell.Y && head.X == bodycell.X {
			issafe["up"] = false
		}

	}
	return issafe
}

func main() {
	RunServer()
}
