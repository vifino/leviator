package ipc

func MakeChans(num int) []chan string {
	retval := make([]chan string, num)
	for i := range retval {
		retval[i] = make(chan string)
	}
	return retval
}

func Broadcast(chans []chan string, msg string) {
	for i := range chans {
		chans[i] <- msg
	}
}

func Send(channel chan string, msg string) {
	channel <- msg
}

func Receive(channel chan string) string {
	return <-channel
}

func ReceiveNonBlocking(channel chan string) string {
	select {
	case msg := <-channel:
		return msg
	default:
		return ""
	}
}
