#!/bin/bash
export DOCKER_IMAGE_TAG=$(git show -s --format=%ct $CI_COMMIT_SHA)
cat dp.json |sed 's/$DOCKER_IMAGE_TAG/'"$DOCKER_IMAGE_TAG"'/g' |kubectl apply -f -
kubectl apply -f svc.json
