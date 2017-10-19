package emitter_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/telenor-digital-asia/kinesis-connectors/emitter"
)

func Test_Key(t *testing.T) {
	d := time.Now().UTC().Format("2006/01/02")

	k := emitter.Key("", "a", "b")
	if k != fmt.Sprintf("%v/a-b", d) {
		t.Fail()
	}
	k = emitter.Key("prefix", "a", "b")
	if k != fmt.Sprintf("prefix/%v/a-b", d) {
		t.Fail()
	}
}
