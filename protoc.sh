#!/bin/sh

# protoc.sh generates Go code for the Starlink Device gRPC service.
# First it calls grpcurl[1] to fetch the protobuf file descriptor set using server reflection.
# Using the descriptor set, protoc[2] generates the Go code for the proto files we care about.
#
# By default it connects to the default Dishy address (see DefaultDishyAddr).
# A different address can be specified as an argument. For example:
#
#	protoc.sh 172.19.248.42:9100
#
# This can be useful if the device is only available through a tunnel or VPN.
#
# [1]: https://pkg.go.dev/github.com/fullstorydev/grpcurl
# [2]: https://grpc.io/docs/protoc-installation/

usage="usage: protoc.sh [address]"

# The Dishy exposes its gRPC service on this adress by default.
default_dishy_addr=`grep DefaultDishyAddr *.go | awk '{print $4}' | tr -d '"'`

if test $# -ge 2
then
	echo $usage
	exit 1
fi

addr=$default_dishy_addr
if test -n "$1"
then
	addr=$1
fi
grpcurl -plaintext -protoset-out device.protoset $addr describe SpaceX.API.Device.Device

# The protoc command, without proto file arguments, to generate the Go code.
# From https://grpc.io/docs/languages/go/quickstart/#regenerate-grpc-code
# For a protobuf reference:
# https://developers.google.com/protocol-buffers/docs/reference/go-generated

# Strip the spacex.com/api prefix so package imports can be resolved within our Go module.
# https://developers.google.com/protocol-buffers/docs/reference/go-generated#invocation
module_opt="--go_opt module=spacex.com/api --go-grpc_opt module=spacex.com/api"

protoc --go_out . --go-grpc_out . --descriptor_set_in device.protoset \
	$module_opt \
	spacex/api/common/status/status.proto

# Device proto files import status.proto, so set the corresponding Go
# package import path for the status package we generated previously.
gomodule=`grep module go.mod | awk '{print $2}'`
import_opt="--go_opt Mspacex/api/common/status/status.proto=$gomodule/status"
files="spacex/api/device/command.proto
	spacex/api/device/common.proto
	spacex/api/device/device.proto
	spacex/api/device/dish.proto
	spacex/api/device/transceiver.proto
	spacex/api/device/wifi.proto
	spacex/api/device/wifi_config.proto"
protoc --go_out . --go-grpc_out . --descriptor_set_in device.protoset \
	$module_opt $import_opt \
	$files
