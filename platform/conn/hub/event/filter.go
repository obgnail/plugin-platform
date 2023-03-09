package event

type Filter func(topic string, payload []byte, filter map[string][]string) bool

func DummyFilter(topic string, payload []byte, filter map[string][]string) bool {
	return true
}
