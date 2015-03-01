# prodda [![Build Status](https://travis-ci.org/mfine30/prodda.svg?branch=master)](https://travis-ci.org/mfine30/prodda)

Prods things on schedule.

## Overview

Prods are scheduled tasks. Scheduling is defined using the extended cron syntax:

- Full crontab specs e.g. `"* * * * * ?"`
- Descriptors, e.g. `"@midnight", "@every 1h30m"`

## API reference

### Root Endpoint

The root endpoint for the API is at `/api/v0`. All endpoints are nested below this path,

### Authentication and authorization

All API requests must be made using basic authentication e.g:

`curl https://username:password@127.0.0.1/api/v0/`

### Prods endpoint

The endpoint for managing prods is found at `/prods/`. The contents of the request body for creates and updates must contain a `schedule` field (the contents of which must be valid cron syntax) and sufficient information to create or update a task, which varies by the type of task.

#### Get all prods

```
curl -XGET /prods/
```

#### Create new prod

```
curl -XPOST /prods/ -d '{<prod-body-as-json>}'
```

#### Get specific prod

```
curl -XGET /prods/:id
```

#### Update existing prod

```
curl -XPOST /prods/:id -d '{<updated-prod-body-as-json>}'
```

#### Delete existing prod

```
curl -XDELETE /prods/:id
```

## Supported tasks

Currently prodda supports executing the following tasks:

### Travis builds

All travis tasks require a travis token. This can be obtained via the following:

```
gem install travis && travis login && travis token
```

More detailed information can be found on the official [travis blog](http://blog.travis-ci.com/2013-01-28-token-token-token/).

#### Re-running an existing travis build

Re-running a specific travis build can be accomplished by creating a new prod with the following body:

```
{
  "token":"my-travis-token",
  "buildID":123456789,
  "schedule":"15 03 * * *"
}
```

## Supported Golang versions

The code is tested against the latest version of golang 1.4.

## Getting the code

The [**develop**](https://github.com/mfine30/prodda/tree/develop) branch is where active development takes place; it is not guaranteed that any given commit will be stable.

The [**master**](https://github.com/mfine30/prodda/tree/master) branch points to a stable commit. All tests should pass.

### Dependency management

Dependencies are managed via [godep](http://godoc.org/github.com/tools/godep). To ensure the dependencies are correct, run `godep restore`. Adding a new dependency requires running `godep save`; the resultant changes in the `Godeps/` directory must be committed.

### Go vet

```
go tool vet -composites=false $(ls -d */ | grep -v Godeps)
go tool vet -composites=false *.go
```

### Running tests

Running the tests will require [ginkgo](http://onsi.github.io/ginkgo/) and [gomega](http://onsi.github.io/gomega/):

```
go get github.com/onsi/ginkgo/ginkgo
go get github.com/onsi/ginkgo/gomega
```

In the cloned directory run the following command:

```
ginkgo -p -r -race -failOnPending -randomizeAllSpecs
```

## Contributing

The upcoming work for prodda can be found on [Pivotal Tracker](https://www.pivotaltracker.com/n/projects/1272036).

Rull-requests are welcome; please make them against the [**develop**](https://github.com/mfine30/prodda/tree/develop) branch. Please also include ginkgo tests. Multiple small commits are preferred over a single monolithic commit.
