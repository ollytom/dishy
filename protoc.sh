#!/bin/sh

# protoc.sh generates Go code for the Starlink Device gRPC service.
# First it calls grpcurl[1] to fetch the protobuf file descriptor set using server reflection.
# Using the descriptor set, protoc[2] generates the Go code for the proto files we care about.
#
# By default it connects to the default Dishy address (see DefaultDishyAddr).
# A different address can be specified as an argument. For example:
#
#	protoc.sh 192.0.2.1:9200
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

# First compile the dependencies of the Device service
deps="spacex/api/common/status/status.proto
	spacex/api/common/protobuf/internal.proto
	spacex/api/telemetron/public/common/time.proto
	spacex/api/satellites/network/ut_disablement_codes.proto
	spacex/api/common/protobuf/internal.proto"

# Map the dependencies between protofiles to Go package import paths.
# We need a clean import path with a different root (not the
# automatically generated "spacex.com") so that packages in this module
# can import each other.
# https://developers.google.com/protocol-buffers/docs/reference/go-generated
gomodule=`grep module go.mod | awk '{print $2}'`
import_opt="
	--go_opt Mspacex/api/common/status/status.proto=$gomodule/status
	--go_opt Mspacex/api/common/protobuf/internal.proto=$gomodule/protobuf
	--go-grpc_opt Mspacex/api/common/protobuf/internal.proto=$gomodule/protobuf
	--go_opt Mspacex/api/satellites/network/ut_disablement_codes.proto=$gomodule/satellites
	--go-grpc_opt Mspacex/api/satellites/network/ut_disablement_codes.proto=$gomodule/satellites
	--go_opt Mspacex/api/telemetron/public/common/time.proto=$gomodule/telemetron
	--go-grpc_opt Mspacex/api/telemetron/public/common/time.proto=$gomodule/telemetron"

# Strip our module prefix when writing generated files to disk; we import packages from within this module.
# https://developers.google.com/protocol-buffers/docs/reference/go-generated#invocation
module_opt="--go_opt module=$gomodule --go-grpc_opt module=$gomodule"

protoc --go_out . --go-grpc_out . --descriptor_set_in device.protoset \
	$module_opt $import_opt \
	$deps

files="spacex/api/device/command.proto
	spacex/api/device/common.proto
	spacex/api/device/device.proto
	spacex/api/device/dish.proto
	spacex/api/device/dish_config.proto
	spacex/api/device/transceiver.proto
	spacex/api/device/wifi.proto
	spacex/api/device/wifi_config.proto"

# Strip the original spacex.com/api prefix so package imports are resolved from our Go module.
module_opt="--go_opt module=spacex.com/api --go-grpc_opt module=spacex.com/api"

protoc --go_out . --go-grpc_out . --descriptor_set_in device.protoset \
	$import_opt $module_opt \
	$files
