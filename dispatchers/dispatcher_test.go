package dispatchers

import (
	"fmt"
	"testing"
)

func TestRegisterDispatcher(t *testing.T) {
	RegisterDispatcher("bla", NewNoopDispatcher)

	dispatcher := MustGetDispatcher(DispatcherConfig{
		Dispatcher: "bla",
	})

	_, ok := dispatcher.(*NoopDispatcher)
	if !ok {
		t.Error("Expected MustGetDispatcher to return NoopDispatcher")
	}
}

func TestMustGetDispatcher(t *testing.T) {
	defer func() {
		err := recover()
		if err == nil {
			t.Error("Shoud have panicked while getting non-existing-dispatcher")
		}
		msg := fmt.Sprintf("%s", err)
		if msg != "Dispatcher 'non-existing-dispatcher' not found" {
			t.Error("Unexpeceted error messsage:", msg)
		}
	}()

	MustGetDispatcher(DispatcherConfig{
		Dispatcher: "non-existing-dispatcher",
	})
}
