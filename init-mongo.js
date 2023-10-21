apiVersion: v1
kind: Pod
metadata:
  name: mongodb
spec:
  containers:
  - name: mongodb
    image: mongo:latest
    command:
    - mongod
    - --port
    - 27017
    - --dbpath
    - /data/db
    volumeMounts:
    - name: data
      mountPath: /data/db
    - name: config
      mountPath: /docker-entrypoint-initdb.d
  volumes:
  - name: data
    emptyDir: {}
  - name: config
    configMap:
      name: mongodb-init
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: mongodb-init
data:
  init-mongo.js: |-
    db.getSiblingDB("mydb").createUser({
      user: "myuser",
      pwd: "mypassword",
      roles: [ { role: "readWrite", db: "mydb" } ]
    })
    db.mydb.insertOne({ name: "John Doe" })
    db.mydb.insertOne({ name: "Jane Doe" })
