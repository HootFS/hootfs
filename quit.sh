#!/bin/bash
# Read a string with spaces using for loop
processes=$(ps aux | grep '[d]iscovery_main' | awk '{print $2}')
length=${#processes[@]}
if [ -z "$processes" ]
then
	echo "No discovery_main processes to kill!"
else
	kill $processes
fi
