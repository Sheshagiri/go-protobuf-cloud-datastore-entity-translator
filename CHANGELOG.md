# Changelog

All notable changes to this project will be documented in this file.

## [1.0.1] - 11th September 2019
This release adds support for excluding the fields from being indexed by google cloud datastore using a custom extension
in proto files. Following is an example.
```
datastoreEntity, err := translator.ProtoMessageToDatastoreEntity(dbModel, true, "models.exclude_from_index")
```
