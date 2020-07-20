# ConfigMap Map Operator

```bash
git clone https://github.com/joesonw/configmap-map-operator
kubectl apply -f ./configmap-map-operator/deploy
```

```bash
apiVersion: operators.dstream.cloud/v1alpha1
kind: ConfigMapMap            
metadata:         
  name: app-configmap-map
spec:                                                           
  name: app-configmap
  namespace: default
  data:                        
    mongo.yaml:
      kind: configmap
      namespace: default
      name: config-mongo
      subPath: mongo.yaml
    harbor:
      kind: secret
      namespace: default
      name: harbor
      subPath: .dockerconfigjson
EOF

```