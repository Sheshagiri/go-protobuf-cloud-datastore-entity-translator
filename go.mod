module github.com/Sheshagiri/go-protobuf-cloud-datastore-entity-translator

go 1.12

require (
	cloud.google.com/go v0.38.0
	github.com/golang/protobuf v1.3.2
	github.com/google/uuid v1.1.1
	github.com/stretchr/testify v1.4.0
	google.golang.org/genproto v0.0.0-20190708153700-3bdd9d9f5532
	gotest.tools v0.0.0-20181223230014-1083505acf35
)

replace cloud.google.com/go => github.com/Sheshagiri/google-cloud-go v0.41.1-0.20190711043959-301311007500
