package core

import (
	"fmt"
)

//Initable represents something that can be initialized (usually a service)
type Initable interface {
	Init() error
}

//InitMany will initialize all the passed Initables and it will return on first error
func InitMany(things ...Initable) error {
	for _, thing := range things {
		err := thing.Init()
		if err != nil {
			return fmt.Errorf("intiialize %T: %w", thing, err)
		}
	}
	return nil
}
