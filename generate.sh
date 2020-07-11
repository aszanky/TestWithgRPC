#!/bin/bash

protoc justmessagepb/justmessage.proto --go_out=plugins=grpc:.