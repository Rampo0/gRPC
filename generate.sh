#!/bin/bash

protoc .\greet\greetpb\greet.proto --go_out=plugins=grpc:.
protoc .\calculator\calculatorpb\calculatorpb.proto --go_out=plugins=grpc:.
protoc .\blog\blogpb\blog.proto --go_out=plugins=grpc:.
