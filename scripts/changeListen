#!/bin/bash

path=$1;

while inotifywait -e modify $path ; do 
	echo ELOG ; cat $path ;
done
