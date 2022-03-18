apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: $DOMAIN-ingress
  labels:
    weblogic.domainUID: $DOMAIN
  annotations:
    kubernetes.io/ingress.class: "contour"
    kubernetes.io/tls-acme: "true"
    ingress.kubernetes.io/proxy-body-size: "0"
    ingress.kubernetes.io/ssl-redirect: "true"
    cert-manager.io/cluster-issuer: "letsencrypt"
  generation: 1
spec:
  rules:
  - host: $INGRESS_HOST
    http:
        paths:
        - backend:
            service:
                name: $DOMAIN-cluster-app-server
                port:
                    number: 8001
          path: /
          pathType: ImplementationSpecific
  tls:
  - hosts:
    - $INGRESS_HOST
    secretName: $DOMAIN-cert
