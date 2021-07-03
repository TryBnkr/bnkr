# Bnkr

The Kubernetes backup solution that designed for human beings!
Bnkr is a height performance backup application written in go.

## Current Roadmap

- [x] Support MongoDB backup & restore.
- [ ] Accept kubeconfig
- [ ] Support PostgreSQL backup & restore.
- [ ] Support S3 compilable object storages.
- [ ] Use Redis for cache and session.
- [ ] Helm chart.

## What I can do with Bnkr?

Bnkr goal is to backup only your valuable data inside Kubernetes cluster not the whole cluster, any contribution in this regard is more than welcome.

## Installation

Bnkr itself is single binary application however it depends on some other tools and applications, like for example it uses `PostgreSQL` to store its data and uses `mysqldump` to backup MySQL databases so it is better to use Bnkr official image because you will guarantee that all the dependencies already exist.

Here I'll mention everything you need to get Bnkr up and running in your cluster but of course you can ignore any part that you already have or you don't need it.

First create a service account with:

```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: backup-sa
  namespace: default
```

Create cluster admin:

```yaml
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRoleBinding
metadata:
  name: backup-cluster-admin
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cluster-admin
subjects:
- kind: ServiceAccount
  name: backup-sa
  namespace: default
```

Create a pvc for PostgreSQL:

```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: postgres
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 5Gi
```

Create PostgreSQL StatefulSet with headless service(**don't forget to change the password!**):

```yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: postgres
spec:
  selector:
    matchLabels:
      app: postgres
  serviceName: postgres
  replicas: 1
  template:
    metadata:
      labels:
        app: postgres
    spec:
      containers:
        - name: postgres
          image: postgres:13.3-alpine
          env:
            - name: POSTGRES_PASSWORD
              value: postgres
          ports:
            - name: postgres
              containerPort: 5432
          volumeMounts:
            - name: data
              mountPath: /data/db
      volumes:
        - name: data
          persistentVolumeClaim:
            claimName: postgres
---
apiVersion: v1
kind: Service
metadata:
  name: postgres
spec:
  clusterIP: None
  selector:
    app: postgres
  ports:
    - name: db
      protocol: TCP
      port: 5432
      targetPort: 5432
```

Next create the Bnkr deployment with NGINX Ingress:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: bnkr
spec:
  replicas: 1
  selector:
    matchLabels:
      app: bnkr
  template:
    metadata:
      labels:
        app: bnkr
    spec:
      serviceAccountName: backup-sa
      containers:
        - name: bnkr
          image: engrmth/bnkr
          imagePullPolicy: Always
          env:
            - name: PORT
              value: "5000"
            - name: USERNAME
              value: "John Doe"
            - name: USERPASSWORD
              value: password
            - name: USEREMAIL
              value: me@example.com
            - name: SETUP
              value: "true"
            - name: PRODUCTION
              value: "true"
            - name: DB_HOST
              value: postgres-0.postgres
            - name: DB_USER
              value: postgres
            - name: DB_PASSWORD
              value: postgres
            - name: DB_NAME
              value: postgres
            - name: DB_PORT
              value: "5432"
            - name: DB_TIMEZONE
              value: Europe/Istanbul
            - name: DB_SSLMODE
              value: disable
---
apiVersion: v1
kind: Service
metadata:
  name: bnkr-clusterip-srv
spec:
  selector:
    app: bnkr
  ports:
    - name: bnkr
      protocol: TCP
      port: 80
      targetPort: 5000
---
apiVersion: networking.k8s.io/v1beta1
kind: Ingress
metadata:
  name: bnkr-ingress-srv
  annotations:
    kubernetes.io/ingress.class: nginx
    nginx.ingress.kubernetes.io/use-regex: 'true'
    nginx.ingress.kubernetes.io/proxy-body-size: '0'
    nginx.ingress.kubernetes.io/proxy-read-timeout: '600'
    nginx.ingress.kubernetes.io/proxy-send-timeout: '600'
spec:
  rules:
    - host: bnkr.example.com
      http:
        paths:
          - path: /
            backend:
              serviceName: bnkr-clusterip-srv
              servicePort: 80
```

## License

Copyright (c) 2021 [Mohammed Al-Mahdawi](https://al-mahdawi.is/)

Licensed under the **MIT** license.