package txn

import (
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assert := assert.New(t)
	txnMgr := MkTxnMgr()
	assert.Equal(0, len(txnMgr.tidsActive))

	txn := txnMgr.New()
	assert.Equal(0, len(txn.wset))
}

func TestPutCommitAndGet(t *testing.T) {
	assert := assert.New(t)
	txnMgr := MkTxnMgr()

	txnPut := txnMgr.New()
	txnPut.Put(10, 20)
	txnPut.Put(11, 22)
	txnPut.Commit()

	txnGet := txnMgr.New()
	v, found := txnGet.Get(10)
	assert.Equal(true, found)
	assert.Equal(uint64(20), v)

	v, found = txnGet.Get(11)
	assert.Equal(true, found)
	assert.Equal(uint64(22), v)

	_, found = txnGet.Get(12)
	assert.Equal(false, found)
	txnGet.Commit()
}

func TestPutAbortAndGet(t *testing.T) {
	assert := assert.New(t)
	txnMgr := MkTxnMgr()

	txnPut := txnMgr.New()
	txnPut.Put(10, 20)
	txnPut.Abort()

	txnGet := txnMgr.New()
	_, found := txnGet.Get(10)
	assert.Equal(false, found)
	txnGet.Commit()
}

/**
 * Interleaved `txnPut.Put` and `txnGet.Get` with `txnPut.tid < txnGet.tid`.
 */
func TestInterleavedPutAndGet1(t *testing.T) {
	assert := assert.New(t)
	txnMgr := MkTxnMgr()

	txnPut := txnMgr.New()
	txnGet := txnMgr.New()

	txnPut.Put(10, 20)

	go func() {
		time.Sleep(10 * time.Millisecond)
		txnPut.Commit()
	}()

	v, found := txnGet.Get(10)
	assert.Equal(true, found)
	assert.Equal(uint64(20), v)

	txnGet.Commit()
}

/**
 * Interleaved `txnPut.Put` and `txnGet.Get` with `txnPut.tid > txnGet.tid`.
 */
func TestInterleavedPutAndGet2(t *testing.T) {
	assert := assert.New(t)
	txnMgr := MkTxnMgr()

	txnGet := txnMgr.New()
	txnPut := txnMgr.New()

	txnPut.Put(10, 20)

	go func() {
		time.Sleep(10 * time.Millisecond)
		txnPut.Commit()
	}()

	_, found := txnGet.Get(10)
	assert.Equal(false, found)

	txnGet.Commit()
}

func TestInOrderPuts(t *testing.T) {
	assert := assert.New(t)
	txnMgr := MkTxnMgr()

	txnA := txnMgr.New()
	txnB := txnMgr.New()

	txnA.Put(10, 20)
	txnA.Commit()

	ok := txnB.Put(10, 200)
	assert.Equal(true, ok)
	txnB.Commit()

	txnGet := txnMgr.New()
	v, found := txnGet.Get(10)
	assert.Equal(true, found)
	assert.Equal(uint64(200), v)
	txnGet.Commit()
}

/**
 * `Put` fails due to later writers.
 */
func TestFailedReversedPuts(t *testing.T) {
	assert := assert.New(t)
	txnMgr := MkTxnMgr()

	txnA := txnMgr.New()
	txnB := txnMgr.New()

	txnB.Put(10, 20)
	txnB.Commit()

	ok := txnA.Put(10, 200)
	assert.Equal(false, ok)
	txnA.Abort()

	txnGet := txnMgr.New()
	v, found := txnGet.Get(10)
	assert.Equal(true, found)
	assert.Equal(uint64(20), v)
	txnGet.Commit()
}

func TestInOrderGetAndPut(t *testing.T) {
	assert := assert.New(t)
	txnMgr := MkTxnMgr()

	txnA := txnMgr.New()
	txnB := txnMgr.New()

	txnA.Get(10)
	txnA.Commit()

	ok := txnB.Put(10, 20)
	assert.Equal(true, ok)
	txnB.Commit()

	txnGet := txnMgr.New()
	v, found := txnGet.Get(10)
	assert.Equal(true, found)
	assert.Equal(uint64(20), v)
	txnGet.Commit()
}

/**
 * `Put` fails due to later readers.
 */
func TestFailedReversedGetAndPut(t *testing.T) {
	assert := assert.New(t)
	txnMgr := MkTxnMgr()

	txnA := txnMgr.New()
	txnB := txnMgr.New()

	txnB.Get(10)
	txnB.Commit()

	ok := txnA.Put(10, 200)
	assert.Equal(false, ok)
	txnA.Abort()

	txnGet := txnMgr.New()
	_, found := txnGet.Get(10)
	assert.Equal(false, found)
	txnGet.Commit()
}

/**
 * `Put` fails due to concurrent writes (`txnA` writes first).
 */
func TestFailedConcurrentPuts1(t *testing.T) {
	assert := assert.New(t)
	txnMgr := MkTxnMgr()

	txnA := txnMgr.New()
	txnB := txnMgr.New()

	txnA.Put(10, 20)

	ok := txnB.Put(10, 200)
	assert.Equal(false, ok)

	txnA.Commit()
	txnB.Abort()

	txnGet := txnMgr.New()
	v, found := txnGet.Get(10)
	assert.Equal(true, found)
	assert.Equal(uint64(20), v)
	txnGet.Commit()
}

/**
 * `Put` fails due to concurrent writes (`txnB` writes first).
 */
func TestFailedConcurrentPuts2(t *testing.T) {
	assert := assert.New(t)
	txnMgr := MkTxnMgr()

	txnA := txnMgr.New()
	txnB := txnMgr.New()

	txnB.Put(10, 20)

	ok := txnA.Put(10, 200)
	assert.Equal(false, ok)

	txnA.Abort()
	txnB.Commit()

	txnGet := txnMgr.New()
	v, found := txnGet.Get(10)
	assert.Equal(true, found)
	assert.Equal(uint64(20), v)
	txnGet.Commit()
}

func TestReadMyWrite(t *testing.T) {
	assert := assert.New(t)
	txnMgr := MkTxnMgr()

	txn := txnMgr.New()
	txn.Put(10, 20)
	v, found := txn.Get(10)
	assert.Equal(true, found)
	assert.Equal(uint64(20), v)
}

func TestWriteMyWrite(t *testing.T) {
	assert := assert.New(t)
	txnMgr := MkTxnMgr()

	txn := txnMgr.New()
	txn.Put(10, 20)
	txn.Put(10, 200)
	v, found := txn.Get(10)
	assert.Equal(true, found)
	assert.Equal(uint64(200), v)

	txn.Commit()

	txnRd := txnMgr.New()
	v, found = txnRd.Get(10)
	assert.Equal(true, found)
	assert.Equal(uint64(200), v)
}

func TestActiveTxns(t *testing.T) {
	assert := assert.New(t)
	txnMgr := MkTxnMgr()

	txn := txnMgr.New()
	assert.Equal(1, len(txnMgr.tidsActive))
	txn.Commit()
	assert.Equal(0, len(txnMgr.tidsActive))

	txnA := txnMgr.New()
	txnB := txnMgr.New()
	txnC := txnMgr.New()
	assert.Equal(3, len(txnMgr.tidsActive))

	txnC.Abort()
	txnA.Commit()
	txnB.Commit()
	assert.Equal(0, len(txnMgr.tidsActive))
}

func TestMinActiveTxns(t *testing.T) {
	assert := assert.New(t)
	txnMgr := MkTxnMgr()

	txns := make([]*Txn, 10)
	for i := 0; i < 10; i++ {
		txns[i] = txnMgr.New()
	}
	assert.Equal(uint64(1), txnMgr.getMinActiveTID())

	txns[0].Commit()
	txns[1].Commit()
	txns[2].Abort()
	txns[5].Commit()
	txns[9].Abort()
	assert.Equal(uint64(4), txnMgr.getMinActiveTID())

	txns[7].Commit()
	assert.Equal(uint64(4), txnMgr.getMinActiveTID())

	txns[3].Abort()
	assert.Equal(uint64(5), txnMgr.getMinActiveTID())

	txns[4].Commit()
	txns[6].Commit()
	txns[8].Abort()
	assert.Equal(uint64(10), txnMgr.getMinActiveTID())
}

