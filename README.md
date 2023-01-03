# dishy

[![Go Documentation](https://godocs.io/olowe.co/dishy?status.svg)](https://godocs.io/olowe.co/dishy)

Package dishy provides an interface to manage a Starlink Dishy device over the network.

## `dishy` command

The `dishy` command controls a dishy device over the network.

To install using the go tools:

	go install olowe.co/dishy/cmd/dishy@latest

## gRPC Code Generation

A starlink device supports gRPC [server reflection][reflection].
This is used to eventually generate code to interact with its gRPC service.
Regenerating the code requires:

- [grpcurl][grpcurl]
- [protoc][protoc]
- connectivity to a Starlink Dishy's gRPC service listening on the default address (192.168.100.1:9200).

To regenerate the code:

	go generate

This calls protoc.sh. For more information on protoc.sh, see the inline documentation in the script.

[grpcurl]: https://pkg.go.dev/github.com/fullstorydev/grpcurl
[protoc]: https://grpc.io/docs/protoc-installation/
[reflection]: https://grpc.github.io/grpc/cpp/md_doc_server-reflection.html
