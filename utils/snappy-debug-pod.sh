#!/usr/bin/env bash
#copyright (c) 2018 SnappyData, Inc. All rights reserved.
# 
# Licensed under the Apache License, Version 2.0 (the "License"); you
# may not use this file except in compliance with the License. You
# may obtain a copy of the License at
#   
# http://www.apache.org/licenses/LICENSE-2.0
#     
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
# implied. See the License for the specific language governing
# permissions and limitations under the License. See accompanying
# LICENSE file.
#

usage() {
echo "Usage: snappy-debug-pod.sh --pvc <comma separated pvc list> [--namespace <namespace in which the PVCs exist>] [--image <docker cordinates on SnappyData image to be used>]"
echo "Options namespace and image are optional. Default namespace is 'spark'"
echo ""
echo "This script launches a pod in the K8S cluster with user specified persistent volumes mounted on it. User must provide list of persistent volume claims as an input for the volumes to be mounted. 
Volumes will be mounted on the path starting with /data (volume1 on /data0 and so on). This script can be used to inspect logs on volumes even when SnappyData system is not online"
echo ""
echo "Example usage: snappy-debug-pod.sh --pvc snappy-disk-claim-snappydata-leader-0,snappy-disk-claim-snappydata-server-0 --namespace default"
}

namespace=spark
image="snappydatainc/snappydata:1.0.1.1-test.1" 
imageString='"image": '\""$image"\"','
mountString=""
volumeString=""
gotPVC=false

while (( "$#" )); do
  case "$1" in
    --pvc)
      pvclist=$2
      gotPVC=true
      
      # read the comma separated PVC names into an array
      oIFS="$IFS"; IFS=, ; 
      read -r -a array <<< "$pvclist";
      IFS="$oIFS"

      # using PVCs provided by the user, 
      # create valid JSON string for volumes and volumeMounts atrbutes of a pod
      for index in "${!array[@]}"
      do
        if [ $index -eq 0 ]
        then
          mountString="{\"mountPath\": \"/data$index/\", \"name\": \"snappy-disk-claim$index\"}"
          volumeString="{\"name\": \"snappy-disk-claim$index\", \"persistentVolumeClaim\": {\"claimName\": \"${array[index]}\"}}"
        else
          mountString="$mountString, {\"mountPath\": \"/data$index/\", \"name\": \"snappy-disk-claim$index\"}"
          volumeString="$volumeString, {\"name\": \"snappy-disk-claim$index\", \"persistentVolumeClaim\": {\"claimName\": \"${array[index]}\"}}"
        fi  
#        echo "$index ${array[index]}"
#        echo $mountString
#        echo $volumeString
      echo "Volume for ${array[index]} will be mounted on /data$index"
      done
      shift 2
    ;;
    --namespace)
      namespace=$2
      shift 2
    ;;
    --image)
      image=$2
      imageString='"image": '\""$image"\"','
      shift 2
    ;;
    --help)
      usage
      exit 0 
    ;;  
    *)
      break
    ;;
  esac
done

if [ "$gotPVC" = false ] ; then
    echo 'ERROR: PVC list not provided'
    usage
    exit 1
fi


# debug-pod-override-template contains, the template JSON in which 
# JSON string for volumes and volumeMounts is added
# first create a copy to modify JSON 
cp debug-pod-override-template.json /tmp/debug-pod-override-actual.json
sed -i 's|IMAGE_MARKER|'"$imageString"'|; s|VOLUME_MOUNTS_MARKER|'"$mountString"'|; s|VOLUMES_MARKER|'"$volumeString"'|' /tmp/debug-pod-override-actual.json

#run the actual command that will launch a pod
overrides=$(</tmp/debug-pod-override-actual.json)
rm -f /tmp/debug-pod-override-actual.json 
cmd="kubectl run -i --rm --tty snappy-debug-pod --overrides='$overrides' --image="$image" --restart=Never --namespace="$namespace" -- bash"

echo "Launching the POD"

#echo "command to be run is=$cmd"
# Fix this. Can we avoid writing the command to a temp file 
echo $cmd > /tmp/kubecommand
chmod +x /tmp/kubecommand
/tmp/kubecommand 
rm -f /tmp/kubecommand
#$cmd
