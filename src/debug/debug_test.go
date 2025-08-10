package debug

import (
	"testing"
)

func TestLog_NoPanic(t *testing.T) {
	Log("test log: %s", "value")
}

func TestErrMsg_NoError(t *testing.T) {
	ErrMsg(nil, "should not exit")
}
