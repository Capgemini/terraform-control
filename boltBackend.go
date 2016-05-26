package main

import (
	"bytes"
	"encoding/json"
	"os"
	"io"
	"path/filepath"
	"github.com/boltdb/bolt"
	"encoding/binary"
	"github.com/capgemini/terraform-control/persistence"
)

var (
	boltEnvironmentsBucket  = []byte("environments")
	boltBlobBucket  = []byte("blob")
	boltBuckets     = [][]byte{
		boltEnvironmentsBucket,
		boltBlobBucket,
	}
)

var (
	boltDataVersion byte = 1
)

// BoltBackend - directory where data will be written. This directory will be
// created if it doesn't exist.
type BoltBackend struct {
  Dir string
}

// GetBlob Function to persist in bolt
func (b *BoltBackend) GetBlob(k string) (*persistence.BlobData, error) {
	db, err := b.db()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var data []byte
	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(boltBlobBucket)
		data = bucket.Get([]byte(k))
		return nil
	})
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}

	// We have to copy the data since it isn't valid once we close the DB
	data = append([]byte{}, data...)

	return &persistence.BlobData{
		Key:  k,
		Data: bytes.NewReader(data),
	}, nil
}

// PutBlob function to persist data and update bucket
func (b *BoltBackend) PutBlob(k string, d *persistence.BlobData) error {
	db, err := b.db()
	if err != nil {
		return err
	}
	defer db.Close()

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, d.Data); err != nil {
		return err
	}

	return db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(boltBlobBucket)
		return bucket.Put([]byte(k), buf.Bytes())
	})
}

// GetAllEnvironments returns all environments persisted in bolt
func (b *BoltBackend) GetAllEnvironments() ([]*Environment, error) {
	db, err := b.db()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var result []*Environment
	var env *Environment
	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(boltEnvironmentsBucket).Bucket([]byte(
			boltEnvironmentsBucket))

		// If the bucket doesn't exist, we haven't written this yet
		if bucket == nil {
			return nil
		}

	    c := bucket.Cursor()

	    count := 0
	    for k, v := c.First(); k != nil; k, v = c.Next() {
	        //fmt.Printf("key=%s, value=%s\n", k, v)
	        env = &Environment{}
	        err := b.structRead(env, v)
	        if err != nil {
				return err
	        }
	        result = append(result, env)
	        count = count+1
	    }
		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetEnvironment returns environment specific details stored
func (b *BoltBackend) GetEnvironment(id int) (*Environment, error) {
	db, err := b.db()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var result *Environment
	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(boltEnvironmentsBucket).Bucket([]byte(
			boltEnvironmentsBucket))

		// If the bucket doesn't exist, we haven't written this yet
		if bucket == nil {
			return nil
		}

		// Get the key for this infra
		data := bucket.Get([]byte(itob(id)))
		if data == nil {
			return nil
		}

		result = &Environment{}
		return b.structRead(result, data)
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

// PutEnvironment adds or changes details for existing environment
func (b *BoltBackend) PutEnvironment(environment *Environment) error {

	db, err := b.db()
	if err != nil {
		return err
	}
	defer db.Close()

	return db.Update(func(tx *bolt.Tx) error {

		bucket := tx.Bucket(boltEnvironmentsBucket)
		bucket, err = bucket.CreateBucketIfNotExists([]byte(
			boltEnvironmentsBucket))
		if err != nil {
			return err
		}

		if environment.Id == 0 {
			id, _ := bucket.NextSequence()
	        environment.Id = int(id)
		}

		data, err := b.structData(environment)
		if err != nil {
			return err
		}

		return bucket.Put(itob(environment.Id), data)
	})
}

// itob does some maths function
func itob(v int) []byte {
    b := make([]byte, 8)
    binary.BigEndian.PutUint64(b, uint64(v))
    return b
}

// db returns the database handle, and sets up the DB if it has never
// been created.
func (b *BoltBackend) db() (*bolt.DB, error) {
	// Make the directory to store our DB
	if err := os.MkdirAll(b.Dir, 0755); err != nil {
		return nil, err
	}

	// Create/Open the DB
	db, err := bolt.Open(filepath.Join(b.Dir, "environments.db"), 0644, nil)
	if err != nil {
		return nil, err
	}

	// Create the buckets
	err = db.Update(func(tx *bolt.Tx) error {
		for _, b := range boltBuckets {
			if _, err := tx.CreateBucketIfNotExists(b); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (b *BoltBackend) structData(d interface{}) ([]byte, error) {
	// Let's just output it in human-readable format to make it easy
	// for debugging. Disk space won't matter that much for this data.
	return json.MarshalIndent(d, "", "\t")
}

func (b *BoltBackend) structRead(d interface{}, raw []byte) error {
	dec := json.NewDecoder(bytes.NewReader(raw))
	return dec.Decode(d)
}
