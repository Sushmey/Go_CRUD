Run a `Docker pull postgres:17.5` to pull the docker image of the postgres version you want.

Run `docker run --name go_backend_db -e POSTGRES_PASSWORD=password -d postgres:17.5` to whip up the postgres container. Make sure to add the version in postgres else it will download the most recent version.


Use `docker container start <name_of_container>` to start an already existing container.

Just connected to the postgresql server inside the connected from the outside. Connection details should be the exposed ones and not of the docker container.

In my case 

- Host address -> localhost
- port -> 8001 (or  whatever mentioned, if not mentioned should be 5432. If that fails create another container and explictly mention a port)
- username -> postgres (or whatever used to create)
- password -> ! (password for db)


## What is this project? 
A GoLang backend which is a CRUD API for inventory management. Could be a library inventory management, a to-do list, or even a bug-tracker. 
Change the model and controller accordingly.

## Implement interfaces in Golang
Interfaces are implemented implictly without any keyword. To implement an interface, you need to define all the functions of it.

## Steps
- Initalized a go.mod file by running `go mod init <name_of_module>`
- Import gorilla mux for routing
- Import gorm and gorm postgres driver (gorm is Go ORM)

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
