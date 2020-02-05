package meteringclient

import "fmt"

type Error struct {
	status  int
	message string
}

func (e Error) Error() string {
	return fmt.Sprintf("status: %d,message: %s", e.status, e.message)
}
