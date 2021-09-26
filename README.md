# imghitcntr [![Build Status](https://app.travis-ci.com/drbig/imghitcntr.svg?branch=master)](https://app.travis-ci.com/github/drbig/imghitcntr)

A very simple "hit counter" like in the olden days of the simpler web. Pass it
a key, perhaps desired background and foreground colors, and presto you have an
image-based counter. The digits I believe are based on the Topaz font, because
the olden days always look brighter from the future.

## Showcase

```bash
$ ./imghitcntr-linux-amd64-0.4.0 --help
Usage: ./imghitcntr-linux-amd64-0.4.0 (options...)
imghitcntr v0.4.0 by Piotr S. Staszewski, see LICENSE.txt
binary build by drbig@swordfish on Sun 26 Sep 14:10:57 CEST 2021

Options:
  -b string
        hostname/ip to bind to (default "127.0.0.1")
  -bg string
        background color, HTML hex string (default "#fff")
  -csv string
        path to save and load the CSV data dump
  -endpoint string
        endpoint to mount at (default "/hit")
  -fg string
        foreground color, HTML hex string (default "#000")
  -loglevel string
        log level (default "error")
  -p int
        port to bind to (default 9999)
```

**Example**: ![I'm a counter](https://tensor.work/hit?key=gh-drbig-imghitcntr&bg=%23fff&fg=%23000)

## Benchmarks

Don't really remember the context anymore, but hey, numbers are cool.

The Alpha-mask based, as it is:
```
goos: linux
goarch: amd64
pkg: github.com/drbig/imghitcntr
BenchmarkGenImage-4       205166              6043 ns/op            2963 B/op          7 allocs/op
PASS
ok      github.com/drbig/imghitcntr     1.305s
```

And the crappy pixel-banging versions, I thought will be faster:
```
goos: linux
goarch: amd64
pkg: github.com/drbig/imghitcntr
BenchmarkGenImage-4        41206             33944 ns/op            5013 B/op        608 allocs/op
PASS
ok      github.com/drbig/imghitcntr     1.697s
```

Is actually 5.5x slower. Great! It was also longer.

## Contributing

Follow the usual GitHub development model:

1. Clone the repository
2. Make your changes on a separate branch
3. Make sure you run `gofmt` and `go test` before committing
4. Make a pull request

See licensing for legalese.

## Licensing

Standard two-clause BSD license, see LICENSE.txt for details.

Any contributions will be licensed under the same conditions.

Copyright (c) 2020 - 2021 Piotr S. Staszewski
