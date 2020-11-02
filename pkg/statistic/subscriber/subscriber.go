package subscriber

//Subscriber interface represents a subscriber which returns channel
//to read messages from
type Subscriber interface {
	Subscribe() <-chan []byte
}
