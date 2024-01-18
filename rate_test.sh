#!/usr/bin/bash

for i in {1..20}
do
    curl -i -X GET "http://localhost:8080/couriers"
done
