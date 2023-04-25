package publicip

import (
	"testing"
)

func TestV4(t *testing.T) {
    ip, _ := V4()
    if ip == nil {
        t.Errorf("expected an ip address but got nil")
    }
}
