package main

import "fmt"

type Strategy func(string, string) string

type Strategic interface {
  SetStrategy(Strategy)
  Result() string
}

type secretStrategy struct {
  first string
  second string
  result string
  strategy Strategy
}

func (ss *secretStrategy) SetStrategy(s Strategy) {
  ss.strategy = s
}

func (ss *secretStrategy) Result() string {
  ss.result = ss.strategy(ss.first, ss.second)
  return ss.result
}

func New(first string, second string) Strategic {
  return &secretStrategy{first: first, second: second}
}

func main() {
  strat := New("vasko", "zdravevski")
  appendd := func(a string, b string) string {
    return a + b
  }
  strat.SetStrategy(appendd)
  fmt.Println(strat.Result())
  prepend := func(a string, b string) string {
    return b + a
  }
  strat.SetStrategy(prepend)
  fmt.Println(strat.Result())
  doubleAndPrepend := func(a string, b string) string {
    return b + b + a + a
  }
  strat.SetStrategy(doubleAndPrepend)
  fmt.Println(strat.Result())
}
