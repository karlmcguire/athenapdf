#!/bin/bash

echo 'Running DEV on docker for port 7330';
docker run -it --rm -p 7330:8080 --dns 1.1.1.1 --env-file weaver.env --shm-size 1024m lachee/athenapdf-service:dev