#go-rencode

[![GoDoc](https://godoc.org/github.com/gdm85/go-rencode?status.png)](https://godoc.org/github.com/gdm85/go-rencode)

go-rencode is a Go implementation of [aresch/rencode](https://github.com/aresch/rencode).

The rencode logic is similar to [bencode](https://en.wikipedia.org/wiki/Bencode). For complex, heterogeneous data structures with many small elements, r-encodings take up significantly less space than b-encodings.

#Usage

You can use either specific methods to encode one of the supported types, or the interface-generic `Encode()` method.

The `DecodeNext()` method can be used to decode the next value from the rencode stream.

#Credits

* This Go version: [gdm85](https://github.com/gdm85)
* Original Python version: Petru Paler, Connelly Barnes et al.
* Cython version: [Andrew Resch](https://github.com/aresch)

##License

go-rencode is licensed under GNU GPL v2, see [COPYING] (./COPYING) for license information.
