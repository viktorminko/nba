package subscriber

type Subscriber interface {
	Subscribe() <-chan []byte
}
