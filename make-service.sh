#!/bin/bash
rm -rf weaver/build/
mkdir weaver/build/
docker build --rm -t lachee/athenapdf-service -f weaver/Dockerfile weaver/