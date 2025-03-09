package main

import (
	"encoding/gob"
	"github.com/google/btree"
	"os"
)

// Number of records to write to the data file.
const maxRecords = 10000
const maxDisplaySamples = 20

// Record types.
const rtypeBeginFrame = 'B'
const rtypeEndFrame = 'E'
const rtypeI64Change = 'I'

// Data record definitions.
type RecordBeginFrame struct {
	RID        byte // 'B': begin frame
	ClassName  string
	MethodName string
	MethodType string
}

type RecordEndFrame struct {
	RID        byte // 'E': end frame
	ClassName  string
	MethodName string
	MethodType string
}

type RecordI64Change struct {
	RID      byte // 'I': int variable
	ValueOld int64
	ValueNew int64
}

// IndexRecord is the wrapper for storing data in the B-bTree
type IndexRecord struct {
	Key   int32 // Record number in the data file
	Value int64 // Byte offset in the data file
}

// Force the B-tree to have exportable fields for gob Encoder.
type SerializableBTree struct {
	Records []IndexRecord
}

// Implement the `Less` function for btree.Item interface.
func (this IndexRecord) Less(that btree.Item) bool {
	return this.Key < that.(*IndexRecord).Key
}

// Wrapper for record type identification.
type RecordWrapper struct {
	Type   byte
	Record interface{}
}

// fileSize gets the current size of the file using its open file handle.
func fileSize(file *os.File) (int64, error) {
	info, err := file.Stat()
	if err != nil {
		return -1, err
	}
	return info.Size(), nil
}

func gobRegisterRecords() {
	gob.Register(RecordBeginFrame{})
	gob.Register(RecordI64Change{})
	gob.Register(RecordEndFrame{})
	gob.Register(RecordWrapper{})
}
