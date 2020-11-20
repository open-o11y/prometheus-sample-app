# prometheus-sample-app

This Prometheus sample app generates Prometheus metrics and exposes them at the `/metrics` endpoint

At the same time, the generated metrics are also exposed at the `/expected_metrics` endpoint in the Prometheus remote write format

A health check endpoint also exists at `/`

Steps for running locally:
```bash
docker build . -t prometheus-sample-app
docker run -it -e INSTANCE_ID=test1234 \
-e LISTEN_ADDRESS=127.0.0.1:8080 \
-p 8080:8080 prometheus-sample-app
curl localhost:8080/metrics
```

Note that the port in LISTEN_ADDRESS must match the the second port specified in the port-forward

Another functioning example:

```bash
docker build . -t prometheus-sample-app
docker run -it -e INSTANCE_ID=test1234 \
-e LISTEN_ADDRESS=127.0.0.1:8080 \
-p 9001:8080 prometheus-sample-app
curl localhost:9001/metrics
```