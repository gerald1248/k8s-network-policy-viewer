#!/bin/bash

for NAMESPACE in isolated global ingress-isolated egress-isolated ingress-isolated-whitelist; do
  kubectl delete namespace ${NAMESPACE} 
done

rm -f *.yaml
rm -f *.json
