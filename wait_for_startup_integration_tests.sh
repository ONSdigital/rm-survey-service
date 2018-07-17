#!/usr/bin/env bash
# Hit the survey service endpoint and exit on success or time out after 20 unsuccessful attempts with 2 second interval
for _ in {1..20}; do
    curl -f http://localhost:9090/info
    if [ $? -eq 0 ]
    then
        exit 0
    fi
    sleep 2
done
echo "No successful response from survey service, timing out"
exit 1
