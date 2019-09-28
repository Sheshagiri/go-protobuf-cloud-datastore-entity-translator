module github.com/Sheshagiri/go-protobuf-cloud-datastore-entity-translator

go 1.12

require (
	cloud.google.com/go v0.38.0
	github.com/golang/protobuf v1.3.2
	github.com/google/uuid v1.1.1 // indirect
	github.com/pborman/uuid v0.0.0-20180906182336-adf5a7427709
	github.com/pkg/errors v0.8.1 // indirect
	github.com/stretchr/testify v1.4.0
	golang.org/x/net v0.0.0-20190628185345-da137c7871d7 // indirect
	golang.org/x/sys v0.0.0-20190710143415-6ec70d6a5542 // indirect
	google.golang.org/genproto v0.0.0-20190708153700-3bdd9d9f5532
	google.golang.org/grpc v1.22.0 // indirect
	gotest.tools v0.0.0-20181223230014-1083505acf35
)

replace cloud.google.com/go => github.com/Sheshagiri/google-cloud-go v0.41.1-0.20190711043959-301311007500
