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
		port        = fs.String("port", "", "")
		server      = fs.Bool("server", false, "")
		dtls        = fs.Bool("dtls", false, "")
		keyFile     = fs.String("key-file", "", "")
		certFile    = fs.String("cert-file", "", "")
		resumeCount = fs.Int("resume-count", 0, "")
	)

	fmt.Println("Args:", os.Args[1:])
	if err := fs.Parse(os.Args[1:]); err != nil {
		log.Println(err)
		if os.Getenv("FAIL") == "1" {
			os.Exit(1)
		} else {
			os.Exit(89)
		}
	}
	if *dtls || !*server {
		os.Exit(89)
	}

	if *port == "" {
		log.Fatal("missing -port")
	}

	config := &tls.Config{
		MinVersion:         tls.VersionSSL30,
		MaxVersion:         tls.VersionTLS13,
		InsecureSkipVerify: true,
	}

	if *keyFile != "" {
		cert, err := tls.LoadX509KeyPair(*certFile, *keyFile)
		if err != nil {
			log.Fatal(err)
		}
		config.Certificates = []tls.Certificate{cert}
	}

	for i := 0; i < *resumeCount+1; i++ {
		conn, err := net.Dial("tcp", net.JoinHostPort("localhost", *port))
		if err != nil {
			log.Fatal(err)
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
}
