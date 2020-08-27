#!/bin/bash
set -e

tag=$(git describe --tags)
repo=docker.io/emsixteeen
image=k8s-audit-webhook
image_tag="${repo}/${image}:${tag}"

echo "Building ${image_tag}"
docker build -f Dockerfile -t ${image_tag} .

echo "Pusing ${image_tag}"
docker push ${image_tag}
