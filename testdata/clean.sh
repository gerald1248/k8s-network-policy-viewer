#!/bin/bash

for NAMESPACE in isolated global ingress-isolated egress-isolated; do
  kubectl delete namespace ${NAMESPACE} 
done

rm *.yaml
rm *.json
