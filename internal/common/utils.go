package common

import "fmt"

func NotImplemented() error {
	return fmt.Errorf("not implemented")
}
