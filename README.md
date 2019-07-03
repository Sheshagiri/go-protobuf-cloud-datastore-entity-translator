# go-protobuf-cloud-datastore-entity-translator


[![Build Status](https://travis-ci.org/Sheshagiri/go-protobuf-cloud-datastore-entity-translator.svg?branch=master)](https://travis-ci.org/Sheshagiri/go-protobuf-cloud-datastore-entity-translator)
[![codecov](https://codecov.io/gh/Sheshagiri/go-protobuf-cloud-datastore-entity-translator/branch/master/graph/badge.svg)](https://codecov.io/gh/Sheshagiri/go-protobuf-cloud-datastore-entity-translator)

# Background

This is largely inspired from being able to persist the protocol buffers into google cloud datastore. Though the datastore
api's by protocol buffers, the client SDK in "cloud.google.com/go/datastore". While it works out of the box since the `Put` function
in `datastore` takes an `interface{}` as input argument, is doesn't support if your protobuf has google.profobuf.Struct or 
google.protobuf.Value in it. 

Issue: https://github.com/googleapis/google-cloud-go/issues/1474
