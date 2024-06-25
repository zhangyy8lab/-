#!/bin/bash
set -exu
CertData=`kubectl get user kubectl-action -o jsonpath='{.spec.client-certificate-data}'`
KeyData=`kubectl get user kubectl-action -o jsonpath='{.spec.client-key-data}'`

kubectl get tasks continuous-deploy -o yaml > continuous-deploy.yaml
sed -i -e "s/CertData/$CertData/g" -e "s/KeyData/$KeyData/g"  continuous-deploy.yaml

kubectl apply -f continuous-deploy.yaml

rm -f continuous-deploy.yaml
