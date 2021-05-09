#!/bin/bash

docker build . -f consumer.dockerfile -t eu.gcr.io/strax-dwk/main-app-consumer
docker build . -f producer.dockerfile -t eu.gcr.io/strax-dwk/main-app-producer

docker push eu.gcr.io/strax-dwk/main-app-consumer
docker push eu.gcr.io/strax-dwk/main-app-producer