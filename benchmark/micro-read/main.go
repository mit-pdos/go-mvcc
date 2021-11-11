package main

import (
	"time"
	"fmt"
	"runtime/pprof"
	"flag"
	"os"
	"log"
	"math/rand"
	"go-mvcc/txn"
)

var done bool

func populateDB(txn *txn.Txn, r uint64) {
	for k := uint64(0); k < r; k++ {
		txn.Begin()
		txn.Put(k, 2 * k + 1)
		txn.Commit()
	}
}

func reader(txnMgr *txn.TxnMgr, src rand.Source, chCommitted, chTotal chan uint64, rkeys uint64) {
	var c uint64 = 0
	var t uint64 = 0
	r := int64(rkeys)
	txn := txnMgr.New()
	rd := rand.New(src)
	for !done {
		txn.Begin()
		canCommit := true
		for i := 0; i < 1; i++ {
			k := uint64(rd.Int63n(r))
			_, ok := txn.Get(k)
			if !ok {
				canCommit = false
				break
			}
		}
		if canCommit {
			c++
			txn.Commit()
		} else {
			txn.Abort()
		}
		t++
	}
	chCommitted <-c
	chTotal <-t
}

func startCPUProfiler(cpuprof string) {
	if cpuprof != "" {
		f, err := os.Create(cpuprof)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}
}

func main() {
	txnMgr := txn.MkTxnMgr()
	//txnMgr.StartGC()

	var nthrd int
	var rkeys uint64
	var cpuprof string
	flag.IntVar(&nthrd, "nthrd", 1, "number of threads")
	flag.Uint64Var(&rkeys, "rkeys", 1000, "access keys within [0:rkeys)")
	flag.StringVar(&cpuprof, "cpuprof", "cpu.prof", "write cpu profile to cpuprof")
	flag.Parse()

	chCommitted := make(chan uint64)
	chTotal := make(chan uint64)

	startCPUProfiler(cpuprof)

	txn := txnMgr.New()
	populateDB(txn, rkeys)
	fmt.Printf("Database populated.\n")

	done = false
	for i := 0; i < nthrd; i++ {
		src := rand.NewSource(int64(i))
		go reader(txnMgr, src, chCommitted, chTotal, rkeys)
	}
	time.Sleep(3 * time.Second)
	done = true

	var c uint64 = 0
	var t uint64 = 0
	for i := 0; i < nthrd; i++ {
		c += <-chCommitted
		t += <-chTotal
	}
	rate := float64(c) / float64(t)
	fmt.Printf("committed / total = %d / %d (%f).\n", c, t, rate)
}

