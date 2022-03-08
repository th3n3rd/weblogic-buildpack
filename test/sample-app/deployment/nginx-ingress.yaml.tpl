apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: $DOMAIN-ingress
  labels:
    weblogic.domainUID: $DOMAIN
  annotations:
    kubernetes.io/ingress.class: "nginx"
spec:
  rules:
  - host: $INGRESS_HOST
    http:
      paths:
      - path: /
        pathType: ImplementationSpecific
        backend:
          service:
            name: $DOMAIN-cluster-app-server
            port:
              number: 8001