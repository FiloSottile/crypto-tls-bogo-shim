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
		minVersion  = fs.Int("min-version", tls.VersionSSL30, "")
		maxVersion  = fs.Int("max-version", tls.VersionTLS13, "")
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
	if *dtls {
		os.Exit(89)
	}

	if *port == "" {
		log.Fatal("missing -port")
	}

	config := &tls.Config{
		MinVersion:         uint16(*minVersion),
		MaxVersion:         uint16(*maxVersion),
		InsecureSkipVerify: true,
	}

	if keyLogFile := os.Getenv("SSLKEYLOGFILE"); config.KeyLogWriter == nil && keyLogFile != "" {
		var err error
		config.KeyLogWriter, err = os.OpenFile(keyLogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			log.Fatalf("Cannot open keylog file: %v", err)
		}
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
