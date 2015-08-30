
Two dynamic elements remain in most static websites - search, comment and contact forms.  While both of these can be outsourced
to services like Google Custom Search, Disquis and Wufoo in a development setting it maybe more convient for a proof of concept
to implement these directly in your prototype. Add support for comments and contact forms remains problematic (e.g. Spam management)
but for search should be possible to implement. It could be either an in-memory DB or spill to disc like BoltDB or SQLite file.
The directory to pass to _wssearch -d_ would be scanned, database populated with an inverted keyword list and the result made available
via the _wsjs_ JavaScript API for constructing a search results page or JSON API.

## background resources

+ [BoldDB](https://github.com/boltdb/bolt)
    + [BoltDB article on Progville](https://www.progville.com/go/bolt-embedded-db-golang/)
+ [https://github.com/google/cayley](cayley) - graph database
+ [goleveldb](https://github.com/syndtr/goleveldb) (an implementation of LevelDB in Go)
+ [LedisDB](http://ledisdb.com/)
+ [etcd](https://github.com/coreos/etcd) coreos key/value configuration db

Alternatives that my not work on a Pi

+ [Tiedot](https://github.com/HouzuoGuo/tiedot)



