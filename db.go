package main

import (
	"encoding/json"
	"fmt"
	"log"

	badger "github.com/dgraph-io/badger/v2"
)

// NewDB create a new db
func NewDB() *DB {
	return &DB{}
}

// DB wraps badger db with domain operations
type DB struct {
	db *badger.DB
}

// DBImage the image metadata stored in the database
type DBImage struct {
	Image
	Flag bool
}

func (d *DBImage) toBytes() ([]byte, error) {
	return json.Marshal(d)
}

// dbImageFromBytes convers bytes to DBImage
func dbImageFromBytes(b []byte) (DBImage, error) {
	d := DBImage{}
	err := json.Unmarshal(b, &d)
	return d, err
}

// Seed seeds the database if it's empty
func (d *DB) Seed(records []Image) error {
	imgs, err := d.ListImages()
	if err != nil {
		return err
	}
	if len(imgs) > 0 {
		log.Println("data already exists in db")
		return nil
	}
	dbimgs := []DBImage{}
	for _, img := range records {
		dbimgs = append(dbimgs, DBImage{Image: img})
	}
	err = d.SaveImages(dbimgs...)
	if err != nil {
		return fmt.Errorf("failed to seed images %w", err)
	}
	return nil
}

// Start starts the database at a specific path
func (d *DB) Start(p string) error {
	if p == "" {
		p = "/tmp/badger"
	}
	var err error
	d.db, err = badger.Open(badger.DefaultOptions(p))
	if err != nil {
		return fmt.Errorf("failed to open db %w", err)
	}
	return nil
}

// Close closes the database
func (d *DB) Close() error {
	d.db.Close()
	return nil
}

//ListImages lists all the images in the database
func (d *DB) ListImages() ([]DBImage, error) {
	imgs := []DBImage{}
	err := d.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			_ = k
			err := item.Value(func(v []byte) error {
				img, err := dbImageFromBytes(v)
				if err != nil {
					return err
				}
				imgs = append(imgs, img)
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	return imgs, err
}

// SaveImages saves images
func (d *DB) SaveImages(images ...DBImage) error {
	err := d.db.Update(func(txn *badger.Txn) error {
		for _, v := range images {
			b, err := v.toBytes()
			if err != nil {
				return err
			}
			err = txn.Set([]byte(v.ID), b)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

// FlagImage flags an image
func (d *DB) FlagImage(imgID string) (DBImage, error) {
	var img DBImage
	err := d.db.Update(func(txn *badger.Txn) error {
		k := []byte(imgID)
		itm, err := txn.Get(k)
		if err != nil {
			return err
		}
		err = itm.Value(func(val []byte) error {
			var err error
			img, err = dbImageFromBytes(val)
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return err
		}
		img.Flag = !img.Flag
		b, err := img.toBytes()
		if err != nil {
			return err
		}
		return txn.Set(k, b)
	})
	return img, err
}
