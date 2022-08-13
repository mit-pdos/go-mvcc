package examples

import (
	"github.com/mit-pdos/go-mvcc/txn"
	"github.com/tchajed/goose/machine"
)

func IncrementSeq(txn *txn.Txn, p *uint64) bool {
	v, _ := txn.Get(0)
	*p = v
	if v == 18446744073709551615 {
		return false
	}
	txn.Put(0, v + 1)

	return true
}

func Increment(t *txn.Txn) (uint64, bool) {
	var n uint64
	body := func(txn *txn.Txn) bool {
		return IncrementSeq(txn, &n)
	}
	ok := t.DoTxn(body)
	return n, ok
}

func DecrementSeq(txn *txn.Txn, p *uint64) bool {
	v, _ := txn.Get(0)
	*p = v
	if v == 0 {
		return false
	}
	txn.Put(0, v - 1)
	return true
}

func Decrement(t *txn.Txn) (uint64, bool) {
	var n uint64
	body := func(txn *txn.Txn) bool {
		return DecrementSeq(txn, &n)
	}
	ok := t.DoTxn(body)
	return n, ok
}

func InitializeCounterData(mgr *txn.TxnMgr) {
	// TODO: Initialize key 0 to some value
}

func InitCounter() *txn.TxnMgr {
	mgr := txn.MkTxnMgr()
	InitializeCounterData(mgr)
	return mgr
}

func CallIncrement(mgr *txn.TxnMgr) {
	txn := mgr.New()
	Increment(txn)
}

func CallIncrementTwice(mgr *txn.TxnMgr) {
	txn := mgr.New()
	n1, ok1 := Increment(txn)
	if !ok1 {
		return
	}
	n2, _ := Increment(txn)
	machine.Assert(n1 < n2)
}

func CallDecrement(mgr *txn.TxnMgr) {
	txn := mgr.New()
	Decrement(txn)
}

