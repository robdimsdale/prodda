# prodda [![Build Status](https://travis-ci.org/mfine30/prodda.svg?branch=master)](https://travis-ci.org/mfine30/prodda)

Prods things on schedule.

## Supported Golang versions

The code is tested against the latest versions of golang 1.2, 1.3 and 1.4.

## Getting the code

The [**develop**](https://github.com/mfine30/prodda/tree/develop) branch is where active development takes place; it is not guaranteed that any given commit will be stable.

The [**master**](https://github.com/mfine30/prodda/tree/master) branch points to a stable commit. All tests should pass.

## Running tests

Running the tests will require [ginkgo](http://onsi.github.io/ginkgo/).

In the cloned directory run the following command:

```
ginkgo -p -r -race -failOnPending -randomizeAllSpecs
```

## Contributing

The upcoming work for prodda can be found on [Pivotal Tracker](https://www.pivotaltracker.com/n/projects/1272036)
