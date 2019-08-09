#!/usr/bin/env bash
protoc -I . service_echo.proto --go_out=plugins=grpc:../.
protoc -I . service_greet.proto --go_out=plugins=grpc:../.
protoc -I . message_chunk.proto --go_out=plugins=grpc:../chunk/.