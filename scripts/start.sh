#!/bin/bash

go build -o gomss_publisher ./cmd/server 
nohup ./gomss_publisher > gomss_publisher.log 2>&1 &