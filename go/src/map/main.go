package main

import (
	"time"
)

type uids_t map[uint64]uint8

func main() {
	var uids uids_t
	
	uids = make(uids_t)
	for idx := uint64(0);idx < 100000000;idx++ {
		uids[idx] = 1;	
	}

	time.Sleep(30 * time.Second)

	return;
}
