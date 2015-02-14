package ipc

type IPCData struct {
	id      int
	content string
}

func MakeChans(num int) []chan IPCData {
	retval := make([]chan IPCData, num)
	for i := range retval {
		retval[i] = make(chan IPCData)
	}
	return retval
}

func Broadcast(chans []chan IPCData, id int, msg string) {
	for i := range chans {
		chans[i] <- IPCData{id, msg}
	}
}

func Send(channel chan IPCData, id int, msg string) {
	channel <- IPCData{id, msg}
}

func Receive(channel chan IPCData) (int, string) {
	tmp := <-channel
	return tmp.id, tmp.content
}

func ReceiveNonBlocking(channel chan IPCData) (int, string) {
	select {
	case msg := <-channel:
		return msg.id, msg.content
	default:
		return -1, ""
	}
}
