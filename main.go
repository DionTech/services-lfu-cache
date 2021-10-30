package main

import (
	"fmt"
	"lfu-cache/lfu"
	"runtime"
	"strconv"
	"time"
)

func main() {
	lfu.LFU = lfu.Init(20)

	PrintMemUsage()

	lfu.LFU.Put("foo", "bar")

	PrintMemUsage()

	lfu.LFU.Put("foo2", []string{"foo", "bar", "irgendwas"})
	//var m1, m2 runtime.MemStats
	for i := 0; i < 100; i++ {

		// Allocate memory using make() and append to overall (so it doesn't get
		// garbage collected). This is to create an ever increasing memory usage
		// which we can track. We're just using []int as an example.
		//runtime.ReadMemStats(&m1)
		lfu.LFU.Put(strconv.Itoa(i), i)
		//runtime.ReadMemStats(&m2)
		//memUsage(&m1, &m2)

		// Print our memory usage at each interval
		PrintMemUsage()
		time.Sleep(time.Second)
	}

	PrintMemUsage()

}

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", m.Alloc)
	fmt.Printf("\tTotalAlloc = %v MiB", m.TotalAlloc)
	fmt.Printf("\tHEAPAlloc = %v MiB", m.HeapAlloc)
	fmt.Printf("\tSys = %v MiB", m.Sys)
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
	fmt.Printf("\tMCacheSys = %v\n", m.MCacheSys)
	/*
		fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
		fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
		fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
		fmt.Printf("\tNumGC = %v\n", m.NumGC)
	*/
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func memUsage(m1, m2 *runtime.MemStats) {
	fmt.Println("Alloc:", m2.Alloc-m1.Alloc,
		"TotalAlloc:", m2.TotalAlloc-m1.TotalAlloc,
		"HeapAlloc:", m2.HeapAlloc-m1.HeapAlloc)
}
