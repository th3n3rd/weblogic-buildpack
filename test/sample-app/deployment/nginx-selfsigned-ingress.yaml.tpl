apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: $DOMAIN-ingress
  labels:
    weblogic.domainUID: $DOMAIN
  annotations:
    kubernetes.io/ingress.class: "nginx"
    ingress.kubernetes.io/ssl-redirect: "true"
    cert-manager.io/cluster-issuer: "self-signed-cluster-issuer"
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
  tls:
  - hosts:
    - $INGRESS_HOST
    secretName: $DOMAIN-cert