#!/bin/sh

OUT=${PWD##*/}
GOOS=linux go build -o $OUT