module github.com/Sheshagiri/go-protobuf-cloud-datastore-entity-translator

go 1.12

require (
	cloud.google.com/go v0.38.0
	github.com/golang/protobuf v1.3.2
	github.com/google/go-cmp v0.3.0
	github.com/google/uuid v1.1.1
	github.com/googleapis/gax-go v1.0.3
	github.com/hashicorp/golang-lru v0.5.1
	github.com/pborman/uuid v0.0.0-20180906182336-adf5a7427709
	github.com/pkg/errors v0.8.1
	go.opencensus.io v0.22.0
	golang.org/x/net v0.0.0-20190628185345-da137c7871d7
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45
	golang.org/x/sys v0.0.0-20190710143415-6ec70d6a5542
	golang.org/x/text v0.3.2
	google.golang.org/api v0.7.0
	google.golang.org/appengine v1.6.1
	google.golang.org/genproto v0.0.0-20190708153700-3bdd9d9f5532
	google.golang.org/grpc v1.22.0
	gotest.tools v0.0.0-20181223230014-1083505acf35
)

replace cloud.google.com/go => github.com/Sheshagiri/google-cloud-go v0.41.1-0.20190711043959-301311007500
