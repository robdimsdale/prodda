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

The endpoint for managing prods is found at `/prods/`.

#### Get all prods

```
curl -XGET /prods/
```

#### Create new prod

The contents of the request body for creates must contain a `schedule` field (the contents of which must be valid cron syntax) and sufficient information to create or update a task, which varies by the type of task. See [supported tasks](#supported-tasks) for further information.

```
curl -XPOST /prods/ -d '{<prod-body-as-json>}'
```

#### Get specific prod

```
curl -XGET /prods/:id
```

#### Update existing prod

The contents of the request body for updates must contain a `schedule` field (the contents of which must be valid cron syntax). Updating attributes of a task is not currently supported - instead the recommended approach is to delete the prod and create a new one with the desired attributes.

```
curl -XPUT /prods/:id -d '{<updated-prod-body-as-json>}'
```

#### Delete existing prod

```
curl -XDELETE /prods/:id
```

## <a name="supported-tasks"</a> Supported tasks

Prodda supports multiple task types.

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
  "schedule":"15 03 * * *",
  "task": {
    "type": "travis-re-run",
    "token":"my-travis-token",
    "buildID":123456789
  }
}

```

### URL Get

A URL Get task is one which will perform an get request to the specified URL, logging the response and any errors encountered. The URL should be fully-formed, including the protocol.

Running a url-get task can be achieved by creating a prod with the following body:

```
{
  "schedule":"15 03 * * *",
  "task": {
    "type": "url-get",
    "url": "http://localhost/"
  }
}
```

### No-op

A no-op task is one which will log its start and finish points, sleeping for a configurable duration in between. The duration must comply with the [golang time.ParseDuration specification](http://golang.org/pkg/time/#ParseDuration):

```

ParseDuration parses a duration string.
A duration string is a possibly signed sequence of decimal numbers,
each with optional fraction and a unit suffix, such as "300ms", "-1.5h" or "2h45m".
Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".

```

Running a no-op task can be achieved by creating a prod with the following body:

```
{
  "schedule":"15 03 * * *",
  "task": {
    "type": "no-op",
    "sleepDuration": "1m"
  }
}
```

## Supported Golang versions

The code is tested against the latest version of golang 1.4.

## Getting the code

The [**develop**](https://github.com/mfine30/prodda/tree/develop) branch is where active development takes place; it is not guaranteed that any given commit will be stable.

The [**master**](https://github.com/mfine30/prodda/tree/master) branch points to a stable commit. All tests should pass.

### Dependency management

Dependencies are managed via [godep](http://godoc.org/github.com/tools/godep). To ensure the dependencies are correct, run `godep restore`. Adding a new dependency requires running `godep save`; the resultant changes in the `Godeps/` directory must be committed.

### Git hooks

To set up git hooks, run the following command:

```
./scripts/install-git-hooks
```

This will need to be re-run if any new git hooks are added, but not if existing ones change.

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
