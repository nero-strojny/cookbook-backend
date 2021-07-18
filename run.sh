#!/bin/bash
cd /home/ec2-user/cookbook-backend
docker build -t docker-example .
docker run --detach --publish 8080:8080 docker-example