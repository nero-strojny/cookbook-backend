#!/bin/bash
cd /home/ec2-user/cookbook-backend
go run main.go > ../output.log 2>&1 &