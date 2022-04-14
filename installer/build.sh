#!/bin/bash

docker build -t fpmbuilder .
rm -Rf vauth-*
docker run -v "${PWD}:/usr/src/app" fpmbuilder