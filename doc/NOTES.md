
What is a spreadsheet (or other document types) was a microservice? The 
spreadsheet would be loaded into memory and it could return columns and rows, be searched by cols of data, return basic metadata about cols, rows or the entire spreadsheet. E.g.

    wsspreadsheetapi -p 8081 -h example.org myexcelfile.xlxs

would start a webserver answering to example.org on port 8081. The root
document would return a JSON blob of metadata 

```
    {
        "Collection": "myexcelfile.xlxs",
        "Demensions: { "Columns": 10, "Rows": 100 },
        "Headings": ["Last Name", "First Name", "etc."],
        "Indexes": [],
        "Updated": "SOME_STANDARD_TIME_REPRESENTATION_OF_MODIFIED_DATA"
    }
```

You could then query the spreadsheet by a simple set of URLs

+ "/" would give you the metadata for the whole
+ "/columns" would return an array of column headings
+ /column/COL-ID would return a column from the spread sheet as a JSON array
+ /row/ROW-ID would return a specific row from the spreadhsheet as a JSON array
+ /search would search all the spreadsheet and return a list of column/row ids as a JSON array
+ /search/columns/COL-ID-LIST would # search column(s) and return the row ids that matched as a JSON array
+ /search/rows/ROW-ID-LIST would search the a row and return a list of column ids that matched as a JSON array
+ /cell/COL-ID/ROW-ID would return the contents of a cell

Specific  indexes could be defined via the command line by column id and
those could speed up search results otherwise the request would do a column scan.

Constraints would be that the microservive would be read only and the spreadsheet would not be larger than memory available.

This idea could in theory be extended to other common structured data types (e.g. tab delimtied files, comma separated files, simple JSON blobs).

If the file was stored someplace like a CDN or S3 then _wsspreadsheet_ would use a URL and proxy the content from the file. Not sure how updates would be triggered without some sort of push notification from the static content source. If this file was on disc you could implement a watch channel that triggered a reload on change.



# random notes and links

+ [Otto JS VM](https://github.com/robertkrimen/otto)
    - [otto with setTimeout(), setInterval()](https://github.com/robertkrimen/natto)
+ [sni-based reverse proxying with Golang](http://www.gilesthomas.com/2013/07/sni-based-reverse-proxying-with-golang/)
+ [OAuth2 Proxy](https://github.com/bitly/oauth2_proxy) - 18F has used this for adding single sign-on/authentication for otherwise staticly generated sites.
+ [Go+Apache+FastCGI](https://github.com/bsingr/golang-apache-fastcgi)
+ [Apache FastCGI](http://www.fastcgi.com/mod_fastcgi/docs/mod_fastcgi.html)
+ [FastCGI](https://en.wikipedia.org/wiki/FastCGI)
+ [Mono and FastCGI](http://www.mono-project.com/docs/web/fastcgi/)
+ [PHP FastCGI-Client example](https://github.com/adoy/PHP-FastCGI-Client/)


