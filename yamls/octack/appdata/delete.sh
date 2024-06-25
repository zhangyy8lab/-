#!/bin/bash
kubectl delete sts -n appdata ck-sts
kubectl delete pvc -n appdata ck-data-ck-sts-3 ck-data-ck-sts-2 ck-data-ck-sts-1 ck-data-ck-sts-0
