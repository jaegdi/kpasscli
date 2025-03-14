#!/usr/bin/env bash

function test_kpasscli {
    local input=$1
    local expected_output=$2
    local field=$3

    if [[ -z "$field" ]]; then
        result=$(./dist/kpasscli -c config.yaml -i "$input")
    else
        result=$(./dist/kpasscli -c config.yaml -i "$input" -f "$field")
    fi

    if [[ "$result" == "$expected_output" ]]; then
        echo "Test passed for input: $input"
    else
        echo "Test failed for input: $input. Expected: $expected_output, Got: $result"
        exit 1
    fi
}

# Test cases
test_kpasscli "pw2" "password2"
test_kpasscli "pw1.1" "user1.1" "username"
