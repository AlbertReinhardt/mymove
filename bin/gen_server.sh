#!/bin/bash

gendir=./pkg/gen

rm -rf $gendir
mkdir -p $gendir
swagger generate server -f swagger.yaml -t $gendir --model-package messages --exclude-main -A mymove
