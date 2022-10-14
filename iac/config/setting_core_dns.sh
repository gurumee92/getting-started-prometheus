#!/bin/bash

# module.eks.module.fargate_profile["default"].aws_eks_fargate_profile.this[0]: Still creating... [10s elapsed] 단계 때 실행
aws eks update-kubeconfig --region ap-northeast-2 --name gurumee-test --alias gurumee-test
kubectl patch deployment coredns -n kube-system --type json -p='[{"op": "remove", "path": "/spec/template/metadata/annotations/eks.amazonaws.com~1compute-type"}]'
