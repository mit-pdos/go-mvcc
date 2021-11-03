package txn

import (
	"sync"
	"go-mvcc/gc"
	"go-mvcc/index"
)

type TxnMgr struct {
	latch		*sync.Mutex
	tidCur		uint64
	tidsActive	map[uint64]bool /* or struct{} if Goose supports. */
	idx			*index.Index
	gc			*gc.GC
}

func MkTxnMgr() *TxnMgr {
	txnMgr := new(TxnMgr)
	txnMgr.latch = new(sync.Mutex)
	txnMgr.tidsActive = make(map[uint64]bool)
	txnMgr.idx = index.MkIndex()
	txnMgr.gc = gc.MkGC(txnMgr.idx)
	return txnMgr
}

func (txnMgr *TxnMgr) New() *Txn {
	txnMgr.latch.Lock()
	txnMgr.tidCur++
	tidNew := txnMgr.tidCur

	/* Make a new txn. */
	txn := new(Txn)
	txn.tid = tidNew
	txn.wset = make(map[uint64]WrEnt)
	txn.idx = txnMgr.idx
	txn.txnMgr = txnMgr

	/* Add `tidNew` to the set of active txns. */
	txnMgr.tidsActive[tidNew] = true

	txnMgr.latch.Unlock()
	return txn
}

/**
 * This function is called by `Txn` at commit/abort time.
 */
func (txnMgr *TxnMgr) deactivate(tid uint64) {
	txnMgr.latch.Lock()
	delete(txnMgr.tidsActive, tid)
	txnMgr.latch.Unlock()
}

/**
 * This function returns the minimal TID of the active txns. If there is no
 * active txns, it returns `tidCur`, which is the largest TID `txnMgr` has ever
 * assigned.
 */
func (txnMgr *TxnMgr) getMinActiveTID() uint64 {
	txnMgr.latch.Lock()

	min := txnMgr.tidCur
	for tid := range txnMgr.tidsActive {
		if tid < min {
			min = tid
		}
	}

	txnMgr.latch.Unlock()
	return min
}

func (txnMgr *TxnMgr) startGC() {
	tidMin := txnMgr.getMinActiveTID()
	txnMgr.gc.Start(tidMin)
}

