/**
 * persistence.go - a simple key/value JSON datastore for OttoEngine
 * that persists to disc between invocations of _ws_. 
 * It runs in its own goroutine to keep Push, Get, Set and Pop atomic.
 */
//package persistence
package main

import (
    "bufio"
    "os"
    "path"
    "fmt"
    "time"
    "errors"
    "encoding/json"
)

// Define the basic actions Persistences supports.
const (
    PUSH = iota
    SET  = iota
    GET  = iota // Not need for playback
    POP  = iota
)

//
// A Collection structure. This forms the basic persistence layer
// available OttoEngine Routes. The collection name corresponds to a directory
// that contains one file  - /journal.jsons.
//
// /journal.jsons is a stream of JSON blobs that can be used to recreate
// Collection.db map[string]string
// 
type Collection struct {
    Name string // Name of collection exposed in OttoEngine
    // Collection contents
    db map[string]string
    keys []string
    // Handle to /journal.jsons append file
    journal *os.File
}

//
// Items are the date unit that get journaled.
// key is a string, value is JSON data formatted as string
type Item struct {
    Key string
    Value string
    Action int
    Timestamp string
}


/*
 * Collection level methods
 */
// Open a Collection, check for errors and if Journal playback is needed.
func Open(collection_name string) (*Collection, error) {
    collections_path := os.Getenv("WS_PERSISTENCE")
    basepath := path.Join(collections_path, collection_name)
    fullpath := path.Join(collections_path, collection_name, "journal.jsons")

    // Create an empty Collection structure to populate
    c := new(Collection)
    // Calculate the Colleciton.Name based on the full collection path.
    c.Name = collection_name

    // Read collectin_path, does it exist?
    finfo, err := os.Stat(basepath)
    if err != nil {
        // create directory for collection (mkdir -p)
        err = os.MkdirAll(basepath, 0770)
        if err != nil {
            return nil, err
        }
    } else if finfo.IsDir() == false {
        return nil, errors.New(fmt.Sprintf("%s exists but is not a directory.", basepath))
    }

    // Create the actual journal file if need be.
    finfo, err = os.Stat(fullpath)
    if err != nil {
        // Create an empty journal file
        _, err = os.Create(fullpath)
        if err != nil {
            return nil, err
        }
    }

    fmt.Println("DEBUG running Playback() in Open()")
    err = c.Playback()
    if err != nil {
        fmt.Printf("DEBUG playback failed, %s -> %s\n", fullpath, err)
        return nil, err
    }
    fmt.Println("DEBUG done playing back journal.jsons")

    // OK now we should be ready to get some work done.
    // open /journal.jsons for appending JSON blobs
    journal, err := os.OpenFile(fullpath, os.O_APPEND | os.O_WRONLY, 0770)
    if err != nil {
        fmt.Printf("DEBUG can't open Journal for appending %s -> %s\n", fullpath, err)
        return nil, err
    }
    // Attach the journal append file
    c.journal = journal
    return c, nil
}

// Cleanup and close the collection.
func (c *Collection) Close() error {
    // Sync and close /journal.jsons
    return c.journal.Close()
}

// Return a count of key/value pairs
func (c *Collection) Count() int {
    return len(c.db)
}

// Return a list of object_ids
func (c *Collection) Keys() ([]string, error) {
    // Return a list of all object_ids in collection
    return nil, errors.New("Keys() not implemented.")
}
// Record queue's a action on a item for appending to /journal.jsons.
func (c *Collection) Record(key, value string, action int) error {
    item := fmt.Sprintf(`{"Key": %q, "Value": %q, "Action": %d, "Timestamp": %q}%s`, key, value, action, time.Now(), "\n")
    fmt.Printf("DEBUG Item would be: %s\n", item)
    _, err := c.journal.WriteString(item)
    return err
}

// Playback reads in a Journal and recreates a collection.
func (c *Collection) Playback() error {
    var (
        item *Item
        json_src string
    )
    fmt.Println("DEBUG attempting Playback()")
    collections_path := os.Getenv("WS_PERSISTENCE")
    fullpath := path.Join(collections_path, c.Name, "journal.jsons")


    fmt.Printf("DEBUG opening journal file for playback: %s\n", fullpath)
    file, err := os.Open(fullpath)
    if err != nil {
        fmt.Printf("Problem opening fullpath %s, %v\n", fullpath, c)
        return err
    }
    defer file.Close()

    // Make sure we have an empty map to read into
    if c.db == nil {
        c.db = make(map[string]string)
    }

    fmt.Printf("DEBUG Seting up scanner for %s\n", fullpath)
    scanner := bufio.NewScanner(file)
    for scanner.Scan(){
        json_src = scanner.Text()
        fmt.Printf("DEBUG scanned: %s\n", json_src)
        err := json.Unmarshal([]byte(json_src), &item)
        switch item.Action {
            case PUSH:
                fmt.Printf("DEBUG Pushing %s -> %s\n", item.Key, item.Value)
                c.db[item.Key] = item.Value
            case SET:
                fmt.Printf("DEBUG Setting %s -> %s\n", item.Key, item.Value)
                c.db[item.Key] = item.Value
            case POP:
                fmt.Printf("DEBUG Deleting %s -> %s\n", item.Key, item.Value)
                delete(c.db, item.Key)
                
        }
        if err != nil {
            fmt.Printf("DEBUG can't JSON.parse() item: %s\n", json_src)
        }
        fmt.Printf("DEBUG item is now: %v\n", item)
    }
    fmt.Println("DEBUG finished scanning")
    if err := scanner.Err(); err != nil {
        fmt.Printf("DEBUG scanner error %s\n", err)
        return err;
    }
    // Journal file should be read into map now
    return nil
}

// Write the whole collection to BASE_PATH/db.json.
func (c *Collection) Save() error {
    return errors.New("Save() not implemented.")
}

// Reads an entire collection from BASE_PATH.db.json.
func (c *Collection) Load() error {
    return errors.New("Load() not implemented.")
}

/*
 * Item level methods
 */

// Push creates/updates and in a collection
func (c *Collection) Push(key, value string) error {
    // check to see if we need to initial db map
    if c.db == nil {
        c.db = make(map[string]string)
    }
    // Add to map
    c.db[key] = value
    _, ok := c.db[key]
    if ok == false {
        return errors.New(fmt.Sprintf("Could not add %s -> %s", key, value))
    }
    // Record in Journal
    err := c.Record(key, value, PUSH)
    if err != nil {
        return err
    }
    return nil
}

// Get returns the value for a given key in the collection
func (c *Collection) Get(key string) (string, error) {
    // Get the value of the corresponding object_id, it not found return nil with error.New("ObjectId not found.")
    if c.db == nil {
        return "", errors.New("Map not populated")
    }
    value, ok := c.db[key]
    if ok == false {
        return "", errors.New("Key/Value pair does not exist.")
    }
    // Record GET
    return value, nil
}

// Set only updated an item if it exists, otherwise returns an error
func (c *Collection) Set(key, value string) error {
    if c.db == nil {
        return errors.New("Map not populated")
    }
    _, ok := c.db[key]
    if ok == true {
        c.db[key] = value
        c.Record(key, value, SET)
        return nil
    }
    return errors.New("Key/value pair did not exist.")
}

// Pop removes an item from the collection returning key, value and error
func (c *Collection) Pop(key string) (string, string, error) {
    // Remove a object_id/value pair from the collection
    // Record POP
    return "", "", errors.New("Delete() not implemented.")
}

// temporary main while sorting our Persistence.
func main() {
    var cnt = 0

    test_collection, err := Open("test")
    if err != nil {
        fmt.Printf("DEBUG can't open test_collection: %v -> %s\n", test_collection.db, err)
        os.Exit(1)
    }
    cnt = test_collection.Count()
    fmt.Printf("DEBUG started count of key/value pairs: %d\n", cnt)

    // Try push an item
    err = test_collection.Push("/fred", `{"name": "fred", "cnt": 1, "date": "2015-01-01 00:00:00 PST"}`)
    if err != nil {
        fmt.Printf("DEBUG Push() failed: %s\n", err)
    }
    if test_collection.Count() != cnt {
        fmt.Printf("DEBUG expected one more record: %d != %d\n", cnt, test_collection.Count())
    }


    // Try to GET an item
    value, err := test_collection.Get("/fred")
    if err != nil {
        fmt.Printf("DEBUG Get() failed: %s\n", err)
    }
    fmt.Printf("DEBUG value is: %s\n", value)

    // Try updated an item
    err = test_collection.Set("/fred", `{"name": "fred", "cnt": 2, "date": "2015-01-01 00:00:01 PST"}`)
    if err != nil {
        fmt.Printf("DEBUG Set() failed: %s\n", err)
    }

    // Try to Pop and item
    key, value, err := test_collection.Pop("/fred");
    if err != nil {
        fmt.Printf("DEBUG Pop() failed. %s\n", err)
    }
    fmt.Printf("DEBUG Pop() returned %s, %s, %s\n", key, value, err)


    fmt.Printf("test collection %v\n", test_collection)
    err = test_collection.Close()
    if err != nil {
        fmt.Printf("DEBUG error close %v\n", test_collection)
        os.Exit(1)
    }
}
