#!/bin/sh

docker build -t gogot/test_suite -f ./test_suite/Dockerfile . && docker run -it gogot/test_suite
