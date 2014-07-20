/**
 * persistence.go - a simple JSON data store for OttoEngine. It is designed to run in a goroutine
 * operating where the convertion of JavaScript data structures to JSON blobs has already taken place.
 */
package persistence

import (
    "errors"
)

const (
    SET    atoi
    DELETE atoi
    SYNC   atoi
)

type Collection struct (
    Name string // Name of collection, this gets exposed in OttoEngine
    Path string // Path to disc file(s), this is NOT exposed in OttoEngine, Only Golang knows where it is stored.
    JSON map[string]string // This holds the in-memory database that needs to get persisted to disc between program restarts.
    Hostname string // Used to create ObjectIDs that are portable. Based on MongoDB's ObjectID algorythm.
    MaxID int // A serial counter used to build a MongoDB style ObjectID.
)


// Object ID
func ObjectID() string, error {
    return "", errors.New("ObjectID() not implemented. This is designed to be MongoDB ObjectID compatible.")
}

// Open a collection for persistence. Creates on if necessary. Ties a OttoEngine name 
// with journaling disc storage.
func Open(name string, pathname string) (Collection, error) {
    // 1. See if a jounral file exists, if yes, rename it to .TIMESTAMP and read it in (using Playback), otherwise create one.
    // 2. Write compacted journal and prep for appending changes.
    return nil, errors.New("Open() not implemented.")
}

// Close the collection, flush everything to disc, compacting the journal.
func (c Collection) Close() error {
    return errors.New("Close() not implemented.")
}

// Return the total count of keys stored in the map
func (c Collection) Count() (int, errors) {
    // Return the total number of key/value pairs stored
    return 0, errors.New("Count() not implemented.")
}

// Return a list of keys
func (c Collection) Keys(from, to int) ([]string, error) {
    // Return a list of all keys in collection
    return nil, errors.New("Keys() not implemented.)
}

// Get a record from the Map
func (c Collection) Get(key string) (Item, error) {
    // Get the value of the corresponding key, it not found return nil with error.New("Key not found.")
    return nil, errors.New("Get() not implemented.")
}

// Save a specific key/value pair
func (c Collection) Set(key string, value string) error {
    // Save the key/value pair, queue the update to the journal
    return errors.New("Set() not implemented.")
}

func (c Collection) Delete(key string) error {
    // Remove a key/value pair from the collection
    // queue the update to the journal
    return errors.New("Delete() not implemented.")
}

// Record appends a single change to the DB Journal
func (c Collection) Record(key string, value string, action int, timestamp Time) error {
    // Write a change out to the Journal file.
    return errors.New("Record() not implemented.")
}

// Playback reads in a DB Journal from disc into collection
func (c Collection) Playback(collection Collection) error {
    // Read in a Journal file into an in-memory collection
    return errors.New("Playback() not implemented.")
}

// Write the whole collection to disc.
func (c Collection) Drain() error {
    // Write the entire collection to disc as a DB Journal file.
    return errors.New("Drain() not implemented.")
}

