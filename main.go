package main

import "lfu-cache/lfu"

func main() {
	/** some examples **/

	//this will set the max cache size to 200MB at the Heap
	//lfu.LFU = lfu.Init(20000000)

	//this will set the max cache size to 200MB at the Heap
	//lfu.LFU = lfu.Init(200000000)

	//this will set the max cache size to 2000MB at the Heap
	//lfu.LFU = lfu.Init(2000000000)

	lfu.LFU = lfu.Init(500000)
}
