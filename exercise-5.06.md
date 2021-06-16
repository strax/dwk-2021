# Exercise 5.06

![Annotated landscape](landscape-annotated.png)

## App Definition & Development
1. PostgreSQL is used for pingpong app and project storage.
2. NATS is used as a message bus in project.
3. Helm is used to install applications to the cluster.
4. Argo was used in Part 4.
5. Flagger was used in Part 5.
6. Flux is used to manage the cluster configuration.

## Orchestration & Management
1. Kubernetes was naturally used during the course.
2. K3s uses CoreDNS.
3. Exercise pingpong and main application communicate over gRPC.
4. Contour was used instead of Traefik in Part 5.
5. Linkerd was the service mesh of choice in Part 5.

## Provisioning
1. Google Container Registry was used to store images.
2. Github Package Registry was used too, but it's not in the picture.

## Platform
1. Docker was used to build the container images used during the course.
2. K3s is used as the development cluster.
3. Google Kubernetes Engine was used in Part 3.

## Observability and Analysis
1. Prometheus was used as the telemetry collector.
2. Grafana was used to provide a nice UI for the data collected by Prometheus.
3. Google Stackdriver was used as an out-of-the-box observability solution in Part 3 for GKE.
4. Grafana Loki was used in the local cluster for log collection.
