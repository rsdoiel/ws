
Two dynamic elements remain in most static websites - search, comment and contact forms.  While both of these can be outsourced
to services like Google Custom Search, Disquis and Wufoo in a development setting it maybe more convient for a proof of concept
to implement these directly in your prototype. Add support for comments and contact forms remains problematic (e.g. Spam management)
but for search should be possible to implement. It could be either an in-memory DB or spill to disc like BoltDB or SQLite file.
The directory to pass to _ws -D_ would be scanned, database populated with an inverted keyword list and the result made available
via the _ws_ JavaScript API for constructing a search results page or JSON API.

