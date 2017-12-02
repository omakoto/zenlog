#!/bin/bash

set -e

for test in "${0%/*}/../e2etests/"test*.sh ; do
    echo "$test"
    "$test"
done

