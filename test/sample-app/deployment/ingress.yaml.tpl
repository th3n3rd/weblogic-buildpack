apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: $DOMAIN-ingress
  labels:
    weblogic.domainUID: $DOMAIN
  annotations:
    # use the shared ingress-nginx
    kubernetes.io/ingress.class: "nginx"
    nginx.ingress.kubernetes.io/rewrite-target: /app/\$1
spec:
  rules:
  - host: sample-app.weblogic.k8s
    http:
      paths:
      - path: /(.+)
        pathType: ImplementationSpecific
        backend:
          service:
            name: $DOMAIN-cluster-app-server
            port:
              number: 8001