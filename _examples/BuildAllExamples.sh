#!/bin/sh
set -e
wd=$PWD
mkdir -p $wd/bin
echo $wd
for ex in $(ls); do
if [ "$ex" != "$0" ]; then
cd $ex && make && mv $ex $wd/bin/
cd $wd
fi
done;
