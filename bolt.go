package demo

import (
    "fmt"
    "log"
    "time"
    "path"
    "bytes"
    "runtime"
    "encoding/json"
    "encoding/gob"

    "github.com/boltdb/bolt"
)

var db *bolt.DB
var open bool

func Open() error {
    var err error
    _, filename, _, _ := runtime.Caller(0)  // get full path of this file
    dbfile := path.Join(path.Dir(filename), "data.db")
    config := &bolt.Options{Timeout: 1 * time.Second}
    db, err = bolt.Open(dbfile, 0600, config)
    if err != nil {
        log.Fatal(err)
    }
    open = true
    return nil
}

func Close() {
    open = false
    db.Close()
}

type Person struct {
    ID string
    Name string
    Age string
    Job string
}

func (p *Person) save() error {
    if !open {
        return fmt.Errorf("db must be opened before saving!")
    }
    err := db.Update(func(tx *bolt.Tx) error {
        people, err := tx.CreateBucketIfNotExists([]byte("people"))
        if err != nil {
            return fmt.Errorf("create bucket: %s", err)
        }
        enc, err := p.encode()
        if err != nil {
            return fmt.Errorf("could not encode Person %s: %s", p.ID, err)
        }
        err = people.Put([]byte(p.ID), enc)
        return err
    })
    return err
}

func (p *Person) gobEncode() ([]byte, error) {
    buf := new(bytes.Buffer)
    enc := gob.NewEncoder(buf)
    err := enc.Encode(p)
    if err != nil {
        return nil, err
    }
    return buf.Bytes(), nil
}

func gobDecode(data []byte) (*Person, error) {
    var p *Person
    buf := bytes.NewBuffer(data)
    dec := gob.NewDecoder(buf)
    err := dec.Decode(&p)
    if err != nil {
        return nil, err
    }
    return p, nil
}

func (p *Person) encode() ([]byte, error) {
    enc, err := json.Marshal(p)
    if err != nil {
        return nil, err
    }
    return enc, nil
}

func decode(data []byte) (*Person, error) {
    var p *Person
    err := json.Unmarshal(data, &p)
    if err != nil {
        return nil, err
    }
    return p, nil
}

func GetPerson(id string) (*Person, error) {
    if !open {
        return nil, fmt.Errorf("db must be opened before saving!")
    }
    var p *Person
    err := db.View(func(tx *bolt.Tx) error {
        var err error
        b := tx.Bucket([]byte("people"))
        k := []byte(id)
        p, err = decode(b.Get(k))
        if err != nil {
            return err
        }
        return nil
    }) 
    if err != nil {
        fmt.Printf("Could not get Person ID %s", id)
        return nil, err
    }
    return p, nil
}

func List(bucket string) {
    db.View(func(tx *bolt.Tx) error {
        c := tx.Bucket([]byte(bucket)).Cursor()
        for k, v := c.First(); k != nil; k, v = c.Next() {
            fmt.Printf("key=%s, value=%s\n", k, v)
        }
        return nil
    })
}

func ListPrefix(bucket, prefix string) {
    db.View(func(tx *bolt.Tx) error {
        c := tx.Bucket([]byte(bucket)).Cursor()
        p := []byte(prefix)
        for k, v := c.Seek(p); bytes.HasPrefix(k, p); k, v = c.Next() {
            fmt.Printf("key=%s, value=%s\n", k, v)
        }
        return nil
    })
}

func ListRange(bucket, start, stop string) {
    db.View(func(tx *bolt.Tx) error {
        c := tx.Bucket([]byte(bucket)).Cursor()
        min := []byte(start)
        max := []byte(stop)
        for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; 
            k, v = c.Next() {
            fmt.Printf("%s: %s\n", k, v)
        }
        return nil
    })
}
