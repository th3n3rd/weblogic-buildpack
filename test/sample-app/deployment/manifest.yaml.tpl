apiVersion: "weblogic.oracle/v8"
kind: Domain
metadata:
  name: $DOMAIN
  labels:
    weblogic.domainUID: $DOMAIN
spec:
  domainHomeSourceType: FromModel
  domainHome: /u01/domains/$DOMAIN
  image: "container-registry.oracle.com/middleware/weblogic:12.2.1.4"
  imagePullPolicy: "IfNotPresent"
  webLogicCredentialsSecret:
    name: $DOMAIN-weblogic-credentials
  includeServerOutInPodLog: true
  serverStartPolicy: "IF_NEEDED"
  auxiliaryImageVolumes:
  - name: auxiliaryImage
    mountPath: "/workspace"
  serverPod:
    env:
    - name: JAVA_OPTIONS
      value: "-Dweblogic.StdoutDebugEnabled=false"
    - name: USER_MEM_ARGS
      value: "-XX:+UseContainerSupport -Djava.security.egd=file:/dev/./urandom "
    auxiliaryImages:
    - image: "$CONTAINER_IMAGE"
      imagePullPolicy: IfNotPresent
      volume: auxiliaryImage
  adminServer:
    serverStartState: "RUNNING"
    adminService:
      channels:
      - channelName: default
        nodePort: 30701
  replicas: 1
  clusters:
  - clusterName: app-server
    serverStartState: "RUNNING"
    replicas: 1
  restartVersion: '1'
  configuration:
    model:
      domainType: "WLS"
      modelHome: "/workspace/models"
      wdtInstallHome: "/workspace/weblogic-deploy"
      runtimeEncryptionSecret: $DOMAIN-runtime-encryption-secret
---
apiVersion: v1
kind: Secret
metadata:
  name: $DOMAIN-weblogic-credentials
  labels:
    weblogic.domainUID: $DOMAIN
type: Opaque
data: # admin:admin123
  password: YWRtaW4xMjM=
  username: YWRtaW4=
---
apiVersion: v1
kind: Secret
metadata:
  name: $DOMAIN-runtime-encryption-secret
  labels:
    weblogic.domainUID: $DOMAIN
type: Opaque
data: # secret
  password: c2VjcmV0