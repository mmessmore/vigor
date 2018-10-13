package graphite

import (
	"fmt"
	"log"
	"net"
	"time"
)

func SendMetric(dest string, path string, metric int) {
	conn, err := net.Dial("tcp", dest)
	if err != nil {
		log.Panic(err)
	}

	fmt.Fprintln(conn, path, metric, time.Now().Unix())
	fmt.Println(path, metric, time.Now().Unix())
	conn.Close()
}
