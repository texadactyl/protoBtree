package main

import (
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/google/btree"
	"io"
	"math/rand"
	"os"
)

func analysis(pathData string, pathIndex string) error {
	// Create or open the index file (where the B-tree will be stored)
	indexFile, err := os.Open(pathIndex)
	if err != nil {
		fmt.Println("analysis: Error creating index file:", err)
		return err
	}
	defer indexFile.Close()

	// Create or open the data file (where the records will be stored)
	dataFile, err := os.Open(pathData)
	if err != nil {
		fmt.Println("analysis: Error creating data file:", err)
		return err
	}
	defer dataFile.Close()

	// Load the B-tree and data file back later (simulate load)
	loadedTree, err := loadBTree(indexFile)
	if err != nil {
		fmt.Println("analysis: Error loading B-tree from file:", err)
		return err
	}

	// Use the loaded tree to retrieve data records.
	for ix := 0; ix < maxDisplaySamples; ix++ {
		err = reportData(randomPal(), loadedTree, dataFile)
		if err != nil {
			return err
		}
	}
	reportData(int32(0), loadedTree, dataFile)
	reportData(maxRecords+1, loadedTree, dataFile)

	return nil
}

// loadBTree: Load the B-tree items from the specified open file and reconstructs the B-tree.
func loadBTree(indexFile *os.File) (*btree.BTree, error) {

	var indexRecords []IndexRecord
	decoder := gob.NewDecoder(indexFile)
	if err := decoder.Decode(&indexRecords); err != nil {
		fmt.Println("loadBTree: decoder.Decode failed, err:", err)
		return nil, err
	}

	// Create a B-tree from the index records
	bTree := btree.New(2)
	for _, record := range indexRecords {
		bTree.ReplaceOrInsert(&IndexRecord{Key: record.Key, Value: record.Value})
	}

	return bTree, nil
}

// Given a key, report the associated data record.
func reportData(recordNumber int32, loadedTree *btree.BTree, dataFile *os.File) error {
	item := loadedTree.Get(&IndexRecord{Key: recordNumber})
	if item == nil {
		errMsg := fmt.Sprintf("reportData: Cannot find index recordNumber: %d", recordNumber)
		println(errMsg)
		return errors.New(errMsg)
	}
	indexRecord := item.(*IndexRecord)

	// Seek to the stored byte offset in the data file
	_, err := dataFile.Seek(indexRecord.Value, io.SeekStart)
	if err != nil {
		fmt.Printf("reportData: dataFile.Seek(%d) failed, err: %v\n:", indexRecord.Value, err)
		return err
	}

	// Read and decode the data record.
	decoder := gob.NewDecoder(dataFile)
	var wrapper RecordWrapper
	err = decoder.Decode(&wrapper)
	if err != nil {
		fmt.Println("reportData: decoder.Decode(%d) failed, err: %v\n:", indexRecord.Value, err)
		return err
	}

	// Show data.
	switch wrapper.Type {
	case rtypeBeginFrame:
		record := wrapper.Record.(RecordBeginFrame)
		fmt.Printf("reportData: begin frame: Record %d, datafile offset %d: FQN = %s.%s%s\n",
			indexRecord.Key, indexRecord.Value, record.ClassName, record.MethodName, record.MethodType)
	case rtypeI64Change:
		record := wrapper.Record.(RecordI64Change)
		fmt.Printf("reportData: i64 change: Record %d, datafile offset %d: old = %d, new = %d\n",
			indexRecord.Key, indexRecord.Value, record.ValueOld, record.ValueNew)
	case rtypeEndFrame:
		record := wrapper.Record.(RecordEndFrame)
		fmt.Printf("reportData: end frame: Record %d, datafile offset %d: FQN = %s.%s%s\n",
			indexRecord.Key, indexRecord.Value, record.ClassName, record.MethodName, record.MethodType)
	default:
		fmt.Printf("*** ERROR, reportData: unrecognizable wrapper type: Record %d, datafile offset %d: wrapper type = 0x%02x\n",
			indexRecord.Key, indexRecord.Value, wrapper.Type)
	}

	return nil
}

// We don't want a random number generator range of [0, n). We actually want (0, n].
func randomPal() int32 {
	var rn int32
	for {
		rn = rand.Int31n(maxRecords + 1)
		if rn > 0 {
			break
		}
	}
	return rn
}
