#!/bin/bash

go build ws.go
./ws -otto=true \
   -otto-path="otto-demo"
