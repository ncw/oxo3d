# Oxo3d for web #

This is a version of Oxo3d with a web interface.

The Go code is transpiled into Javascript with GopherJS so all the
code runs in your browser - there is no server component.

You can see it in action here: https://www.craig-wood.com/nick/oxo3d/

See also: the [go/wasm version](https://github.com/ncw/oxo3d/tree/master/oxo3dwasm).

## Testing ##

Run

    gopherjs serve

Then go to http://localhost:8080/github.com/ncw/oxo3d/oxo3dweb
