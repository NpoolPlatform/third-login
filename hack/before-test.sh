#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

PLATFORM=linux/amd64
OUTPUT=./output

mkdir -p $OUTPUT/$PLATFORM


for file in `ls $path`
do
    echo $path"/"$file
done


for service_name in `ls -l $(pwd)/cmd | grep ^d |awk '{print $9}'`; do
    cp $(pwd)/cmd/*.viper.yaml $OUTPUT/$PLATFORM
    cd $OUTPUT/$PLATFORM; ./$service_name run | grep error &
done

sleep 60
