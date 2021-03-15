#!/bin/sh

exec docker build . -t gogot-tests -f ./test_suite/Dockerfile
exec docker run -it gogot-tests
