#!/bin/bash

set -e

EDGEX_ROOT=$(dirname "${BASH_SOURCE}")/..
INTERVAL=0
SERVICE_LIST="consul-service mongo-service logging-service notifications-service metadata-service data-service 
                command-service scheduler-service export-client-service export-distro-service rulesengine-service 
                ui-service device-mqtt-service"

DEPLOY_LIST="consul-deployment mongo-deployment logging-deployment notifications-deployment metadata-deployment 
              data-deployment command-deployment scheduler-deployment export-client-deployment export-distro-deployment 
              rulesengine-deployment ui-deployment device-mqtt-deployment"

function create_services {
  for svc in $SERVICE_LIST
  do
    kubectl apply -f "${EDGEX_ROOT}/services/$svc.yaml" || true
  done
}

function delete_services {
  for svc in $SERVICE_LIST
  do
    kubectl delete -f "${EDGEX_ROOT}/services/$svc.yaml" || true
  done
}

function create_deployments {
  for deploy in $DEPLOY_LIST
  do
    kubectl apply -f "${EDGEX_ROOT}/deployments/$deploy.yaml" || true
    sleep $INTERVAL
  done
}

function delete_deployments {
  for deploy in $DEPLOY_LIST
  do
    kubectl delete -f "${EDGEX_ROOT}/deployments/$deploy.yaml" || true
    sleep $INTERVAL
  done
}

function clear_volumes {
  rm -rf /data/db
  rm -rf /edgex/logs
  rm -rf /consul/*
}

function start {
  echo "Creating EdgeX services now!"
  create_services
  echo "EdgeX services created successfully !"

  echo "Creating EdgeX deployments now!"
  create_deployments
  echo "EdgeX deployments created successfully!"
}

function stop {
  echo "Deleting EdgeX deployments now!"
  delete_deployments 
  echo "EdgeX deployments created successfully!"

  echo "Deleting EdgeX services now!"
  delete_services
  echo "EdgeX services deleted successfully !"

  clear_volumes
}

if [ $1 == 'stop' ];then
  stop
elif [ $1 == 'start' ];then
  start
fi


