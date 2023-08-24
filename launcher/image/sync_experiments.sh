#!/bin/bash
set -euo pipefail

readonly OEM_PATH='/usr/share/oem'
readonly CS_PATH="${OEM_PATH}/confidential_space"
readonly EXPERIMENTS_BINARY="confidential_space_experiments"
readonly SYNC_RESULT_DIRECTORY="/tmp/container_launcher"
readonly SYNC_RESULT_FILE="experiment_data"

getImageDetails() {
    # image will look like `projects/confidential-space-images/global/images/confidential-space-hardened-230600`
    #                   or `projects/confidential-space-images/global/images/confidential-space-debug-230600`
    image=$(curl --silent "http://metadata.google.internal/computeMetadata/v1/instance/image" -H "Metadata-Flavor: Google")

    csVersionString=$(echo $image | tr '/' '\n' | grep confidential-space-[dh])

    # get image name like `confidential-space-debug-230600`
    versionStringSegments=(${csVersionString//-/ })

    imageType="${versionStringSegments[2]}"
    version="${versionStringSegments[3]}"

    local -n arr=$1 
    arr=(${imageType} ${version})
}

main() {
    getImageDetails results
    if [ ! -d "$SYNC_RESULT_DIRECTORY" ]; then
        echo "$SYNC_RESULT_DIRECTORY does not exist, creating it."
        mkdir $SYNC_RESULT_DIRECTORY
    fi
   
    echo "calling $CS_PATH/$EXPERIMENTS_BINARY with flags: output=$SYNC_RESULT_DIRECTORY/$SYNC_RESULT_FILE, image_type=\"${results[0]}\" image_version=\"${results[1]}\""
    $CS_PATH/$EXPERIMENTS_BINARY -output=$SYNC_RESULT_DIRECTORY/$SYNC_RESULT_FILE -image_type="${results[0]}" -image_version="${results[1]}"
}

main
