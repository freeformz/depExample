package math

import (
	"fmt"
	"strconv"

	"github.com/boltdb/bolt"
)

type operation string

const (
	add      operation = "add"
	subtract operation = "subtract"
	value    operation = "value"
)

// Calc / store the named variables
type Calc struct {
	DB *bolt.DB
}

var (
	valueKey = []byte("value")
)

// Add i to the key
func (c *Calc) Add(i int, key string) (v int, err error) {
	err = c.DB.Update(func(tx *bolt.Tx) error {
		v, err = transactionOp(tx, i, key, add)
		return err
	})
	return
}

// Subtract i from the key
func (c *Calc) Subtract(i int, key string) (v int, err error) {
	err = c.DB.Update(func(tx *bolt.Tx) error {
		v, err = transactionOp(tx, i, key, subtract)
		return err
	})
	return
}

// Value of the key
func (c *Calc) Value(key string) (v int, err error) {
	err = c.DB.View(func(tx *bolt.Tx) error {
		v, err = transactionOp(tx, 0, key, value)
		return err
	})
	return
}

func transactionOp(tx *bolt.Tx, i int, key string, op operation) (int, error) {
	k := []byte(key)
	var b *bolt.Bucket
	var err error
	switch op {
	case value:
		b = tx.Bucket(k)
	default:
		b, err = tx.CreateBucketIfNotExists(k)
	}
	if err != nil {
		return 0, err
	}
	var rv []byte
	if b != nil {
		rv = b.Get(valueKey)
	}
	var v int
	switch len(rv) {
	case 0: // zero value
	default:
		v, err = strconv.Atoi(string(rv))
	}
	if err != nil {
		return 0, err
	}
	switch op {
	case add:
		v += i
		return v, b.Put(valueKey, []byte(strconv.Itoa(v)))
	case subtract:
		v -= i
		return v, b.Put(valueKey, []byte(strconv.Itoa(v)))
	case value:
		return v, nil
	default:
		return v, fmt.Errorf("unknown operation: %s", op)
	}
}
