package mock

import (
	"fmt"
	"log"
)

type Controller struct{}

func NewController() *Controller {
	return &Controller{}
}

func (c Controller) Errorf(format string, args ...any) {
	fmt.Printf(format, args)
}

func (c Controller) Fatalf(format string, args ...any) {
	log.Fatalf(format, args)
}
