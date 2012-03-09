package main

import "time"
import "net"
import "fmt"
import "runtime"
import "syscall"
import "strings"

func main() {
	conn, err := net.Dial("tcp", "josh-dev.moovweb.org:2003")
	if err != nil {
		panic("Error connecting to stats server")
	}

	//statTimer := time.NewTicker(1000 * 1000 * 1000 * 5)
	statTimer := time.NewTicker(10)
	var counter uint
	for {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		MemAlloc := m.Alloc/1024
		MemTotalAlloc := m.TotalAlloc/1024
		var MaxMem uint64
		rusage := &syscall.Rusage{}
		ret := syscall.Getrusage(0, rusage)
		if ret == nil && rusage.Maxrss > 0 {
			MaxMem = uint64(rusage.Maxrss)
		}

		timestamp := time.Now().Unix()

		version := strings.Split(runtime.Version(), " ")[0]
		prefix := fmt.Sprintf("carbon.testing.%s.", version)
		buf := fmt.Sprintf(prefix + "rss %d %d\n", MaxMem, timestamp)
		buf += fmt.Sprintf(prefix + "mem_allocated %d %d\n", MemAlloc, timestamp)
		buf += fmt.Sprintf(prefix + "mem_total %d %d\n", MemTotalAlloc, timestamp)
		counter++
		if (counter % 10000) == 0 {
			println(buf)
			conn.Write([]byte(buf))
		}
		<-statTimer.C
	}	
}
