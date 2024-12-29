# GO HTMX Proof of Concept

This is a simple proof of concept using GO, HTMX and dynamic templates.  It only tests the Active Search pattern.

Goals
-----
* Support development using dynamic template reloading but use embedded templates in production
* Get the static part of the web server right
* Support both nested templates and single file templates (partials)
* Use HTMX to dynamically load search results (or any partial template)

Features
--------
* Embeds the [HTMX](https://htmx.org/) files, CSS, and templates in an embedded file system
* Uses the [PICO CSS](https://picocss.com/) framework for simplicity
* Accepts a `-local` flag on the command-line to indicate whether to read templates from the file system for each page load or use the embedded file system
* The `static` package provides the embedded file system
* The `web` package parses and executes templates.  It really should return an object but this works for the POC.

Notes
-----
* Early versions had weird errors related to MIME types.  I think this was a combination of Windows 10, the way I serve static files, and something else I forgot.  The fix is in `app.init()` but is commented out.

TO DO
-----
* Display an HTMX error event back to the user
* Populate a drop down and reload based on choosing an item
* Populate a DIV or page based on a timer
* `web` should really be a structure with fields for `Local`
* Use an accent-insensitive data filter