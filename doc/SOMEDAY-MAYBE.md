
# Someday, maybe

+ Could you support a simple site search without requiring a DB server?
    + BoltDB
    + goleveldb
+ Add support for on the fly compression (gzip) of text/* content types
+ Is there a core set of APIs that would be easy to support for non-static content?
	- examples (currently available)
		+ Fargo integration with Dropbox
		+ Dropbox [Core API](https://www.dropbox.com/developers/core/docs) 
		+ Dropbox [Drop-ins](https://code.google.com/p/google-api-javascript-client/)
		+ Google [Apps Script](https://developers.google.com/apps-script/) APIs
		+ Google [APIs Client Library for JavaScript](https://code.google.com/p/google-api-javascript-client/)

		+ Box [JavaScript APIs](https://code.google.com/p/google-api-javascript-client/)
	- examples (historical)
		+ [Jaxer](http://www.ibm.com/developerworks/library/wa-aj-jaxer/index.html?ca=drs-tp3508) (apache module for embedding server/client JS processing)
		+ [Lively Kernel](http://www.lively-kernel.org/repository/lively-kernel/trunk/doc/website-index.html)
+ Look at [Echo](http://echo.labstack.com/) router and see if this is a useful way to support FastCGI
+ Look at [libsecurity](https://developer.ibm.com/open/libsecurity/) and see how it might help *ws* stay safe.
+ Look at [takama/router](https://github.com/takama/router) as a possible solution for dynamic routes in _wsjs_
+ Look at [TiDB](https://github.com/pingcap/tidb) MySQL protocol compatible distributed database or [ql](https://github.com/cznic/ql) embedable SQL for Golang




