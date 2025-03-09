package main

import (
	"encoding/gob"
	"fmt"
	"github.com/google/btree"
	"os"
)

func capture(pathData string, pathIndex string) error {
	var recordCounter = int32(0)

	// Create or open the index file (where the B-tree will be stored)
	indexFile, err := os.Create(pathIndex)
	if err != nil {
		fmt.Println("capture: Error creating index file:", err)
		return err
	}
	defer indexFile.Close()

	// Create or open the data file (where the records will be stored)
	dataFile, err := os.Create(pathData)
	if err != nil {
		fmt.Println("capture: Error creating data file:", err)
		return err
	}
	defer dataFile.Close()

	// Create a B-tree.
	bTree := btree.New(2)

	var rbfr RecordBeginFrame
	rbfr.RID = rtypeBeginFrame
	rbfr.ClassName = "java/lang/String"
	rbfr.MethodName = "getBytes"
	rbfr.MethodType = "()[B"
	err = writeRecordToFile(dataFile, bTree, recordCounter, RecordWrapper{rtypeBeginFrame, rbfr})
	if err != nil {
		fmt.Println("capture: writeRecordToFile(RecordBeginFrame) failed, err:", err)
		return err
	}

	var ri64chg RecordI64Change
	ri64chg.RID = rtypeI64Change
	ri64chg.ValueOld = int64(0)

	for recordCounter = 1; recordCounter < (maxRecords + 1); recordCounter++ {

		// Create a record.
		ri64chg.ValueNew = int64(recordCounter)

		// Write int64 change record to the data file.
		err := writeRecordToFile(dataFile, bTree, recordCounter, RecordWrapper{rtypeI64Change, ri64chg})
		if err != nil {
			fmt.Println("capture: writeRecordToFile(RecordI64Change) failed, err:", err)
			return err
		}

		// Update old value of int64.
		ri64chg.ValueOld = ri64chg.ValueNew
	}

	// Write end frame record.
	var refr RecordEndFrame
	refr.RID = rtypeEndFrame
	refr.ClassName = "java/lang/String"
	refr.MethodName = "getBytes"
	refr.MethodType = "()[B"
	err = writeRecordToFile(dataFile, bTree, recordCounter, RecordWrapper{rtypeEndFrame, refr})
	if err != nil {
		fmt.Println("capture: writeRecordToFile(RecordEndFrame) failed, err:", err)
		return err
	}

	// Serialize the B-tree and write to the index file
	err = saveBtree(bTree, indexFile)
	if err != nil {
		fmt.Println("capture: saveBtree failed, err:", err)
		return err
	}

	return nil
}

// writeRecordToFile:
// * Write a record to the data file.
// * Insert its file offset associated with the record number into the B-tree.
func writeRecordToFile(dataFile *os.File, bTree *btree.BTree, recordNumber int32, wrapper RecordWrapper) error {

	// Get current file size as the next write position.
	offset, err := fileSize(dataFile)
	if err != nil {
		return err
	}

	encoder := gob.NewEncoder(dataFile)
	err = encoder.Encode(wrapper)
	if err != nil {
		return err
	}

	// Insert index record into the B-tree.
	bTree.ReplaceOrInsert(&IndexRecord{Key: recordNumber, Value: offset})

	// Return success to caller.
	return nil
}

// saveBtree serializes the B-tree and saves it to the file
func saveBtree(tree *btree.BTree, indexFile *os.File) error {
	var records []IndexRecord
	tree.Ascend(func(item btree.Item) bool {
		records = append(records, *(item.(*IndexRecord)))
		return true
	})

	encoder := gob.NewEncoder(indexFile)
	return encoder.Encode(records)
}
