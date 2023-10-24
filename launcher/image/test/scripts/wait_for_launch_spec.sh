 #!/bin/bash
set -euo pipefail
source util/read_serial.sh

# This test requires the workload to run and print
# corresponding messages to the serial console.
SERIAL_OUTPUT=$(read_serial $1 $2 "Launch Spec") 
print_serial=false

if $print_serial; then
    echo $SERIAL_OUTPUT
fi
