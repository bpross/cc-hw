# Post and Caption as a Service

## Problem
Build a RESTful API controller to create, update, or retrieve social media captions for a given URL or ID.

## Solution
My solution uses golang and docker to provide a POC for how I would design and implement a caption service described in the provided documentation. 

### Datastore
My solution uses an in memory map to store the post information. I have provided interfaces, so if we were to productionize this code base, we could just hook this up to our database of choice. I also took the *"Assume a given entry is heavily requested for one hour after the content is sent to a client for approval."* to mean caching. I added a caching layer to my datastore and dao packages. 

Packages of interest:

- datastore
  - this package contains the implementation for the memory map and pseudo-cache implementations
- dao
  - this package contains the "data access object" code, which abstracts out the datastore calls. I have provided a memory map, combined and cache implementations. The combined package is what I am using in my solution, but I wanted to illustrate how easy it is to plug and play different solutions.

#### Discussion
Since the requirements specifically said to not use a datastore, I did not use one. However, it would be pretty simple to plug this code into a postgres/dynamo/etc database. The only layer of code that would need to change would be the dao layer for posting. As far as caching goes, before implementing a caching solution, I would like to see usage statistics and see if we really need to implement a cache. Assuming we find that it makes sense, I would implement the caching layer using redis and most likely a write-through cache with lazy loading. 

### Captions
My solution provides a caption interface and an implementation for that interface using Aylien. Using the assumption *"Assume the response from Aylien API is deterministic. I.e. for a given URL, the summaries will always be the same."* I have implemented a simple in memory "cache" for my Aylien caption generator.

Packages of interest:

- caption
  - this package contains the interface, aylien and cache implementation for caption generator.

#### Discussion
This portion was pretty straight forward and is not very tecnically interesting.

### Handlers
My solution provides three handlers for the three different methods: `POST, PUT and GET`. I used the framework `gin-gonic`, since it seems to perform the best in benchmark tests. My handler interface uses the `gin.Context`, so is tied to that framework. I have implemented a base handler that implements all three methods, as well as a caption generating handler that only implements the `POST` method. I then use composition, so the generating handler implementation satisfies the interface. This is a limitation in `go` as there is not inheritence.

Packages of interest:

- handler
  - this package contains the interface, and implementations for the handler interface.

#### Discussion
I am not married to this framework, or to REST for this if it were to be productionized. If this is an internal-only system, I think gRPC might be a better solution. gRPC works well in multi-language systems, like ours is bound to be. We can easily define and generate code for whatever language we choose.

### Server
My solution creates a simple server that runs on localhost. It lives in `cmd/server/main.go`.

## How to run
### Assumptions
This code was developed using docker, so it is recommended that you have docker and docker installed on your system. Instructions [here](https://docs.docker.com/docker-for-mac/install/). If you choose to not install docker, you will need to have `go` installed on your system. Instructions [here](https://golang.org/doc/install). It is highly recommended that you install docker, as all further instructions use docker commands. Docker also allows all build and test/lint steps to remain the same across developer environments.

### Routes
All routes require the header `x-customer-id` to be set with a string id. The `POST` and `PUT` routes require the header `Content-Type: application/json` to be set.

- `POST /post`
	-  `curl -XPOST -H "Content-Type: application/json" -H "x-customer-id: 1"  localhost:8080/post -d '{"url": "https://blog.cloudcampaign.io/2019/12/04/how-to-register-a-agency-domain/", "captions": ["test1", "test2"]}'`
	-  Body: `{"url": str, "captions": str list}`
- `GET /post/:id`
	- `curl -XGET -H "Content-Type: application/json" -H "x-customer-id: 1" localhost:8080/post/5e154899cb80cb0001000003`
- `PUT /post/:id`
	- `curl -XPUT -H "Content-Type: application/json" -H "x-customer-id: 1"  localhost:8080/post/5e154899cb80cb0001000003 -d '{"captions": ["test1", "test2", "test3"]}'`
	- Body: `{"captions": str list}`  	   

### Get up and running

First, create a file in the root directory `.env`
Add these keys:

- AYLIEN_API_KEY=
- AYLIEN_APP_ID=
- AYLIEN_CAPTION_COUNT=

Run these in order:

- Build the builder image:
  - docker-compose build builder
- Build the binary:
	- docker-compose run --rm builder bin/build
- Build the api image:
	- docker-compose build api
- Run the api image:
	- docker-compse up -d api
- Connect to logs:
	- docker-compose logs -f api

You can now issue commands against `http://localhost:8080` 

### Running unit tests
How do you prove something works without tests?

To run (this assumes you already have the builder container built):

- docker-compose run --rm builder bin/test

### Running lint
Pretty code is important.

To run (this assumes you already have the builder container built):

- docker-compose run --rm builder bin/lint

### Running integration tests
How do I prove to you that my code does what you asked?

To run (this assumes you already have the builder container and api container built):

- Make sure api is up and running
	-  docker-compse up -d api
- Run integration tests
	- docker-compose run --rm builder bin/test_integration