package main

import "lfu-cache/lfu"

func main() {
	lfu.LFU = lfu.Init(20)
}
