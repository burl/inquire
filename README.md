# inquire
> A collection of common interactive command line user interfaces

inquire attempts to replicate the look and feel of the Node package [inquirer](https://www.npmjs.com/package/inquirer), but for Go.

## WIP

This is a work in progress.  See [demo/grail.go](https://github.com/burl/inquire/blob/master/demo/grail.go)
for a demonstration of using the API.

![Demo](https://github.com/burl/inquire/blob/master/data/inquire-demo.gif)

### API

The API is not yet stable.  I'm torn between building a permissive style
API with interface types and maps -- or --  having a strongly typed
interface with compile-time checking, etc.  Perhaps there can be both,
but for now, its the latter.  Suggestions are welcome.

[API documentation](https://godoc.org/github.com/burl/inquire) can
be found at [godoc.org](https://godoc.org/github.com/burl/inquire).

## License
This library is under the [MIT License](http://opensource.org/licenses/MIT)
