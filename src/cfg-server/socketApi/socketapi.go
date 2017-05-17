package socketApi

import (
	"fmt"
	log "github.com/auxten/logrus"
	"net"
)

func handleConn(c net.Conn) {
	defer c.Close()
	for {
		// read from the connection
		var buf = make([]byte, 10)
		log.Debug("start to read from conn")
		n, err := c.Read(buf)
		if err != nil {
			log.Error("conn read error:", err)
			return
		}
		log.Debug("read %d bytes, content is %s\n", n, string(buf[:n]))
	}
}

func SocketServerStart() {
	strListenPort := fmt.Sprintf(":%d", 2000)
	l, err := net.Listen("tcp", strListenPort)
	if err != nil {
		log.Fatal("listen error:", err)
		return
	}

	for {
		c, err := l.Accept()
		if err != nil {
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				continue
			}

			log.Fatal("accept error:", err)
			break
		}

		// start a new goroutine to handle
		// the new connection.
		log.Debug("accept a new connection")
		go handleConn(c)
	}
}
