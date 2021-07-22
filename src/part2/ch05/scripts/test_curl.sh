#!/bin/bash

while :
do
  (curl localhost:8080/) &
  (curl localhost:8080/) &
  (curl localhost:8080/) &
  (curl localhost:8080/) &
  (curl localhost:8080/) &
  (curl localhost:8080/) &
  sleep 1
done