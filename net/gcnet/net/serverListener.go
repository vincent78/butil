package gcnet

type ServerListener struct {
	Id        string
	RemoteIp  string
	LogPrefix string
	SendCH    chan []byte
	ReceiveCH chan []byte
	ErrorCH   chan error
	CloseCH   chan struct{}
}

func NewServerListener(remoteIp string) *ServerListener {
	return &ServerListener{
		Id:        "",
		RemoteIp:  remoteIp,
		LogPrefix: remoteIp,
		SendCH:    make(chan []byte, 10),
		ReceiveCH: make(chan []byte, 10),
		ErrorCH:   make(chan error, 10),
		CloseCH:   make(chan struct{}),
	}
}

func (l *ServerListener) Close() {
	l.CloseCH <- struct{}{}
}
