/**
 * r-marked.js - This is a web component based on x-tags/Brick for rendering Markdown content directly into a page.
 * @author: R. S. Doiel, <rsdoiel@yahoo.com>
 * copyright (c) 2014
 * All Rights Reserved.
 * Released under the BSD 2-Clause License.
 */
/*jslint browser: true, indent: 4 */
/*global xtag, console, ActiveXObject, XDomainRequest, marked */
(function () {
    "use strict";

    // httpGet - Grab content via xhr.
    // Based on example from http://docs.webplatform.org/wiki/apis/xhr/XMLHttpRequest
    // @param url - the URL to get the content from
    // @param callback - the function to execute when the content is available.
    // callback's parameters are error, response object.
    // @param progress - an optional parameter to process progress from get.
    // progress has two parameters a string for readState and response object
    // @return request object or false.
    function httpGET(url, callback, progress) {
        var request;

        if (window.XMLHttpRequest) { // Mozilla, Safari, ...
            request = new XMLHttpRequest();
        } else if (typeof XDomainRequest !== 'undefined') { // IE 9
            request = new XDomainRequest();
        } else if (typeof window.ActiveXObject) { // IE 8 and older
            try {
                request = new ActiveXObject("Msxml2.XMLHTTP");
            } catch (err1) {
                try {
                    request = new ActiveXObject("Microsoft.XMLHTTP");
                } catch (err2) {
                    throw err2;
                }
            }
        }

        if (!request) {
            if (typeof callback !== 'undefined') {
                callback("Can't create request object");
            }
            return false;
        }
        if (typeof callback !== 'undefined') {
            request.onreadystatechange = function () {
                if (typeof progress !== 'undefined') {
                    switch (request.readyState) {
                    case 0:
                        progress("uninitialized", request);
                        break;
                    case 1:
                        progress("loading", request);
                        break;
                    case 2:
                        progress("loaded", request);
                        break;
                    case 3:
                        progress("interactive", request);
                        break;
                    case 4:
                        progress("complete", request);
                        if (request.status === 200) {
                            callback(null, request.responseText);
                        } else {
                            callback({status: request.status, error: "http request failed"}, request.responseText);
                        }
                        break;
                    }
                }
            };
        }
        request.open("GET", url);
        request.send();
        return request;
    }
    
    function resolveURL(doc_url, href) {
        var protocol_re = new RegExp('://'),
            last_slash = doc_url.lastIndexOf('/');
        if (href.indexOf('://') === -1) {
            // Concatentate the doc_url and the href
            if (last_slash === -1) {
                return doc_url + '/' + href;
            }
            return doc_url.substring(0, last_slash) + '/' + href;
        }
        return href;
    }

    function loadContent(elem, url) {
        httpGET(url, function (err, data) {
            if (err) {
                console.log("ERROR", err);
                return;
            }
            if (typeof marked === 'undefined') {
                elem.innerHTML = '<pre>' + data + '</pre>';
            } else {
                marked(data, function (err, content) {
                    elem.innerHTML = content;
                });
            }
        }, function (status) {
            // We'll handle the error when complete hits.
        });
    }

    xtag.register('r-marked', {
        lifecycle: {
            created: function () {
                var protocolRe = new RegExp('://');
                if (typeof this.href !== 'undefined' && this.href) {
                    loadContent(this, resolveURL(document.URL, this.href));
                }
                // Check the URL to see if it is relative.
            },
            attributeChanged: function () {
                if (typeof this.href !== 'undefined' && this.href) {
                    loadContent(this, resolveURL(document.URL, this.href));
                }
            }
        },
        accessors: {
            href: {
                attribute: { url: ""}
            }
        }
    });
}());
