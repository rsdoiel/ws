/**
 * r-args.js - Look at the page's GET args, find targeted marked content 
 * custom elements and have their href element rendering new content.
 * @author: R. S. Doiel, <rsdoiel@yahoo.com>
 * copyright (c) 2014
 * All Rights Reserved.
 * Released under the BSD 2-Clause License.
 */
/*jslint browser: true, indent: 4 */
/*global xtag, console, ActiveXObject, XDomainRequest, marked */
(function () {
    "use strict";
    function parseWindowSearch(arg_string) {
        var args = {},
            kv_pairs;
            
        if (typeof arg_string === 'undefined' || arg_string === '') {
            return {};
        }
        arg_string.substring(1).split("&").forEach(function (pair) {
            var kv_pair = [], key = "", value = "";            
            kv_pair = pair.split('=', 2);
            key = decodeURIComponent(kv_pair[0]);
            value = decodeURIComponent(kv_pair[1]);
            args[key] = value;
        });
        return args;
    }


    xtag.register('r-args', {
        lifecycle: {
            created: function () {
                var mc_elem = document.querySelectorAll("r-marked"),
                    args = parseWindowSearch(window.location.search),
                    ids = Object.keys(args),
                    elem = {},
                    id = '',
                    i = 0;

                for (i = 0; i < mc_elem.length; i += 1) {
                    elem = mc_elem[i];
                    id = elem.getAttribute('id');
                    if (ids.indexOf(id) > -1) {
                        elem.setAttribute('href', args[id]);
                    }
                }
            }
        }
    });
}());
