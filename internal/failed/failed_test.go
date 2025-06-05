package failed

import (
	"fmt"
	"testing"
)

func TestFailed(t *testing.T) {
	failed := Failed{}
	fmt.Printf("%+v", failed.Error())
}
