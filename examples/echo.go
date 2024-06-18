package main

import (
	"io"
	"log"
	"time"

	"github.com/xtaci/kcp-go/v5"
)

func main() {
	go server1()
	server2()
}
func server1() {
	if listener, err := kcp.ListenWithOptions("127.0.0.1:8091", nil, 10, 3); err == nil {
		// spin-up the client
		// go client()
		for {
			s, err := listener.AcceptKCP()
			if err != nil {
				log.Fatal(err)
			}
			go handleEcho(s)
		}
	} else {
		log.Fatal(err)
	}
}
func server2() {
	if listener, err := kcp.ListenWithOptions("127.0.0.1:8092", nil, 10, 3); err == nil {
		// spin-up the client
		ss, err := listener.NewConn("127.0.0.1:8091")
		if err != nil {
			panic(err)
		}
		go Write(ss)
		go handleEcho2(ss)
		for {
			s, err := listener.AcceptKCP()
			if err != nil {
				log.Fatal(err)
			}
			go handleEcho(s)
		}
	} else {
		log.Fatal(err)
	}
}
func Write(s *kcp.UDPSession) {
	for {
		data := time.Now().String()
		buf := []byte(data)
		log.Println("sent:", data)
		s.Write(buf)
		time.Sleep(time.Second * 3)
	}
}

// handleEcho send back everything it received
func handleEcho(conn *kcp.UDPSession) {
	buf := make([]byte, 4096)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("recv==>:" + string(buf[:n]))
		n, err = conn.Write(buf[:n])
		if err != nil {
			log.Println(err)
			return
		}
	}
}
func handleEcho2(conn *kcp.UDPSession) {
	buf := make([]byte, 4096)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("echo 2==>:" + string(buf[:n]))
		time.Sleep(time.Second * 3)
	}
}
func client() {
	// key := pbkdf2.Key([]byte("demo pass"), []byte("demo salt"), 1024, 32, sha1.New)
	// block, _ := kcp.NewAESBlockCrypt(key)

	// wait for server to become ready
	time.Sleep(time.Second)

	// dial to the echo server
	if sess, err := kcp.DialWithOptions("127.0.0.1:8091", nil, 10, 3); err == nil {
		for {
			data := time.Now().String()
			buf := make([]byte, len(data))
			log.Println("sent:", data)
			if _, err := sess.Write([]byte(data)); err == nil {
				// read back the data
				if _, err := io.ReadFull(sess, buf); err == nil {
					log.Println("recv:", string(buf))
				} else {
					log.Fatal(err)
				}
			} else {
				log.Fatal(err)
			}
			time.Sleep(time.Second)
		}
	} else {
		log.Fatal(err)
	}
}
