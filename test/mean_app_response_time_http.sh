#!/bin/bash

echo "Running Mean Application Response Time Test..."
sleep 1

echo "16 bytes payload size"
k6 run -e Q=1 -e RES=16 mean_app_response_http.js
sleep 2

echo "256 bytes payload size"
k6 run -e Q=1 -e RES=256 mean_app_response_http.js
sleep 2

echo "4 kb payload size"
k6 run -e Q=1 -e RES=4k mean_app_response_http.js
sleep 2

echo "64 kb payload size"
k6 run -e Q=1 -e RES=64k mean_app_response_http.js
sleep 2

echo "1 mb payload size"
k6 run -e Q=1 -e RES=1m mean_app_response_http.js

echo "Test Finished"