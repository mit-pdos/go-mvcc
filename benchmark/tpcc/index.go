package tpcc

import (
	"bytes"
	"strings"
	"encoding/binary"
	"github.com/mit-pdos/go-mvcc/txn"
)

func readidx(txn *txn.Txn, gkey uint64) ([]uint64, bool) {
	opaque, found := txn.Get(gkey)
	if !found {
		return nil, false
	}
	ents := decodeidx(opaque)
	return ents, true
}

func writeidx(txn *txn.Txn, gkey uint64, ents []uint64) {
	s := encodeidx(ents)
	txn.Put(gkey, s)
}

/**
 * Encode a slice of global keys pointing to table records to an opaque string.
 */
func encodeidx(gkeys []uint64) string {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, uint64(len(gkeys)))
	binary.Write(buf, binary.LittleEndian, gkeys)
	return buf.String()
}

/**
 * Decode an opaque string to a slice of global keys pointing to table records.
 */
func decodeidx(opaque string) []uint64 {
	var n uint64
	reader := strings.NewReader(opaque)
	binary.Read(reader, binary.LittleEndian, &n)
	gkeys := make([]uint64, n)
	binary.Read(reader, binary.LittleEndian, gkeys)
	return gkeys
}