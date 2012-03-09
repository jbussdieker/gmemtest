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
		runtime.UpdateMemStats()
		MemAlloc := runtime.MemStats.Alloc/1024
		MemTotalAlloc := runtime.MemStats.TotalAlloc/1024
		var MaxMem uint64
		rusage := &syscall.Rusage{}
		ret := syscall.Getrusage(0, rusage)
		if ret == 0 && rusage.Maxrss > 0 {
			MaxMem = uint64(rusage.Maxrss)
		}

		timestamp := time.Seconds()

		version := strings.Split(runtime.Version(), " ")[0]
		prefix := fmt.Sprintf("carbon.testing.%s.", version)
		buf := fmt.Sprintf(prefix + "rss %d %d\n", MaxMem, timestamp)
		buf += fmt.Sprintf(prefix + "mem_allocated %d %d\n", MemAlloc, timestamp)
		buf += fmt.Sprintf(prefix + "mem_total %d %d\n", MemTotalAlloc, timestamp)
		counter++
		if (counter % 100000) == 0 {
			println(buf)
			conn.Write([]byte(buf))
		}
		<-statTimer.C
	}	
}
