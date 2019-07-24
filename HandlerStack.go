package gohttp


type HandlerStack struct {
	Handler func()
	Stack []string

}