## What is this project? 
A GoLang backend which is a CRUD API for inventory management. Could be a library inventory management, a to-do list, or even a bug-tracker. 
Change the model and controller accordingly.

## Steps
1. Clone the repository
2. Run `docker build --tag "name_for_image"` to build an image using the Dockerfile
3. Create a K8s namespace using the `ns.yaml` file
4. Create pods for DB by applying the config files in`/k8s/db` using `kubectl apply -f .`
5. Make sure pods are running using `kubectl get pods -ns <namespace_name>`
6. Create application pods for the go-server by applying the config files in `/k8s/apps` using `kubectl apply -f .`
7. Access the endpoints mentioned in the controller at port 30004 (unless you changed the ports)

# DOCKER
A `Dockerfile` is needed to build a docker container, you can also use it to create a multistage deployment. 
```
FROM golang:1.24.3-alpine AS builder
WORKDIR /app
```
`FROM` specifies what image you want to pull or where you want to start your build from, you can name the stages

```
# Copy files from host
# Copies from first path in host and pastes at second path (workdir) in container
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o Go_CRUD
```
`COPY` the path from host to path in container and add command to be run using `RUN`,  
if your configuration on base machine is different from the target machine, make sure to specify the specs

```
FROM gcr.io/distroless/base
WORKDIR /app
COPY --from=builder /app/Go_CRUD .

CMD ["./Go_CRUD"]
```
This is the second stage, fetch the distroless OS since its lighter.
`COPY` using from the earlier stage with the path mentioned and paste it in the current directory as indicated by `.`

`CMD` to specifcy the commands to run the executable


# KUBERNETES

### Creating a namespace
To make two containers interact with each other in docker, they need to be in the same namespace. And to do that, we need to create a namespace by creating a k8s cluster. 
- Go to settings in Docker Desktop and enable Kubernetes, wait till it starts.
- Create a directory `k8s` and add `ns.yaml` to create a namespace file
- Then to build using that, 
	- run `kubectl create ns <name_of_k8s_ns>` this name is the name given in the ns.yaml file
	- Or you can go to k8s dir and run `kubectl apply -f ns.yaml` This will create the namespace.

### Create a Postgres DB pod
- Create a `pv.yaml` to create a PersistentVolume. This the storage space for your app's pods which exists independently of any specific application.
- Create a `pvc.yaml`. This PersistentVolumeClaim (PVC) is a request for storage by user or app, allows you to specify storage specification and permissions.
- Create a `cm.yaml`. This ConfigMap is to store key-value pair of config data for k8s, allows to separate config from application code making it easier to make changes.
- Create a `svc.yaml` and `deploy.yaml` to create the service and deploy it. In our case its a db so do mention the db config like port, username, dbname, host and password.

### Create the application and connect to DB
- Set up ConfigMap and add env variables to connect to postgres-service we created earlier
- Set up deployment where we create a container using the image created earlier (in step 2) and mention the port which exposes this container
- Create a service file to expose the deployed container to the outside world. In our case we have exposed 8080 (on container) to 30004 in outside world

## Common Errors I ran into while developing this
- ### Connection refused
    Check if the ports are correct. Check if you're in the same namespace and if you're accessing the exposed ports.
- ### K8s pod fails immediately with exit code 1
    Meaning the internal application fails. 
    Try running `kubectl logs <pod_name> -n <namespace_name>` to see the logs. 
    Error was that db we are connecting is uninitialized so added the env tag to initialize the port.
     containers:
            - name: postgres-db
              image: postgres:17.5
              env:
                - name: POSTGRES_PASSWORD
                  value: !
              ports:
