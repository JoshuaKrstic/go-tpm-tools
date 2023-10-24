 #!/bin/bash
set -euo pipefail

# This test requires the workload to run and print
# corresponding messages to the serial console.
FILE_OUTPUT=$(cat $2) 
print_serial=false

if echo $FILE_OUTPUT | grep -q "EnableTestFeatureForImage:$1"
then
    echo "- test experiment verified true"
else
    echo "FAILED: experiment status expected to be true"
    echo 'TEST FAILED. Test experiment status expected to be true' > /workspace/status.txt
    print_serial=true
fi

if $print_serial; then
    echo $FILE_OUTPUT
fi
