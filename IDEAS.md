
# Ideas to explore

## Next

+ Add option -n, -no-cache-header so you shouldn't need to disable cache in web browser while developing sites

## Someday, maybe

+ Integrate [mkpage](https://rsdoiel.github.io/mkpage) site development mode 
+ Add support for ACME SSL/TLS certs (SSL Everywhere)
    + this may wind up in a future version Golang http package
+ Add support for on the fly compression (gzip) of text/* content types
+ Look at [libsecurity](https://developer.ibm.com/open/libsecurity/) and see how it might help *ws* stay safe.
+ Consider adding site search support via Bleve search integration
+ Consider including a mime type map from file exitension setting
+ Consider adding support to treat a XLSX document as a read only DB.
    + would need a good mapping to RESTful URLs row/cell range per URL end point
+ Consider assigning a default template for presenting Markdown files as HTML
    + If file is requested as *.html but only *.md exists, use default template to render as *.html
+ Consider adding support for S3 like options for serving header content


## Old ideas

+ create wsedit for remote editing content of specific files over https connections.
    + look at CodeMirror and AceEditor as candidates for web browser editing
    + research best approach to embedding the editor in the go compiled binary
    + review scripted-editor for general apprach to problem
    + decide how to handle TLS key generation seemlessly
        + use existing certs for host machine
        + create one time self signed certs with signatures in browser display along with onetime URL
+ create a nav generator based on file system organization
    + autogenerate sitemap.xml and siteindex.html for improved search engine ingest
+ develop a generator and JS for browser side site specific search
    + create an inverted word list as JSON file
    + create a sitemap JSON file
+ explore cli tools as CMS to produce static websites
    + a markdown processor for generating static site content
+ explore interfacing with Solr or Bleve for a site search api

