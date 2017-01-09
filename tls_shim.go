package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	fs := &flag.FlagSet{}
	var (
		port     = fs.String("port", "", "")
		server   = fs.Bool("server", false, "")
		dtls     = fs.Bool("dtls", false, "")
		no_tls13 = fs.Bool("no-tls13", false, "")
	)

	fmt.Println("Args:", os.Args[1:])
	if err := fs.Parse(os.Args[1:]); err != nil {
		log.Println(err)
		os.Exit(89)
	}
	if *dtls || !*no_tls13 {
		os.Exit(89)
	}

	if *port == "" {
		log.Fatal("missing -port")
	}
	conn, err := net.Dial("tcp", net.JoinHostPort("127.0.0.1", *port))
	if err != nil {
		log.Fatal(err)
	}

	config := &tls.Config{
		MinVersion:         tls.VersionSSL30,
		InsecureSkipVerify: true,
	}
	var tlsConn *tls.Conn
	if *server {
		tlsConn = tls.Server(conn, config)
	} else {
		tlsConn = tls.Client(conn, config)
	}

	for {
		buf := make([]byte, 500)
		n, err := tlsConn.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		buf = buf[:n]
		for i := range buf {
			buf[i] ^= 0xff
		}
		if _, err := tlsConn.Write(buf); err != nil {
			log.Fatal(err)
		}
	}
}
