#!/bin/bash
cd /home/ec2-user/cookbook-backend
docker-compose build --no-cache
docker-compose up -d