#!/bin/bash

awk '/ERROR/ {print $0}' log.log > error.log
awk '/match/ {print $10}' error.log > error_match.log
echo -e "\n" >> error_match.log
awk '/match/ {print $10}' error.log | grep -oP '(?<=//).*?/' | uniq >> error_match.log
echo -e "\n" >> error_match.log
awk '/match/ {print $10}' error.log | sort | grep -oP '(?<=//).*?/' | uniq >> error_match.log
awk '/:http/ {print $0}' error.log > error_http.log
