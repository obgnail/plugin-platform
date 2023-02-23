package errors

import "testing"

func TestErrors(t *testing.T) {
	err := New(InvalidParameter, "Plugin.Desc", "TooLong")

	err = WithValue(err, "key1", "value1")
	t.Log(err)
}
