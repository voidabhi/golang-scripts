package main

import (
	"container/list"
	"fmt"
)

type Observer interface {
	Update()
}

type Subject interface {
	Attach (observer Observer)
	Detach (observer Observer)
	Notify()
}

type DefaultSubject struct {
	observers *list.List
}

func NewDefaultSubject() *DefaultSubject {
	return &DefaultSubject{observers:new(list.List)}
}

func (this *DefaultSubject) Attach(observer Observer) {
	this.observers.PushBack(observer)
}
func (this *DefaultSubject) Detach(observer Observer) {
	for obs := this.observers.Front(); obs != nil; obs = obs.Next() {
		if obs.Value.(Observer) == observer {
			this.observers.Remove(obs)
		}
	}
}

func (this *DefaultSubject) Notify() {
	for obs := this.observers.Front(); obs != nil; obs = obs.Next() {
		observer := obs.Value.(Observer)
		observer.Update()
	}
}

////  concrete obj

type GameState string

func NewGameState(state string) *GameState {
	gs := GameState(state)
	return &gs
}

type Game struct {
	*DefaultSubject
	state *GameState
}

func NewGame() *Game {
	return &Game{DefaultSubject:NewDefaultSubject()}
}

func (this *Game) GetState() *GameState {
	return this.state
}

func (this *Game) SetState(state *GameState) {
	this.state = state
	this.Notify()
}

type Player struct {
	name         string
	lastState *GameState
	game         *Game
}

func NewPlayer(name string, game *Game) *Player {
	this := new(Player)
	this.name = name
	this.game = game
	return this
}

func (this *Player) Update() {
	this.lastState = this.game.GetState()
	fmt.Println(this.name, "noticed that game state has changed to: ", *this.lastState)
}


func main() {
	var game = NewGame()
	var p1 = NewPlayer("Alex", game)
	var p2 = NewPlayer("Tim", game)

	game.Attach(p1)
	game.Attach(p2)

	st := NewGameState("Game started")
	game.SetState(st)

	st2 := NewGameState("Tim's move")
	p2.game.SetState(st2)

}
