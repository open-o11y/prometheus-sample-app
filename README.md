# prometheus-sample-app

This Prometheus sample app generates all 4 Prometheus metric types (counter, gauge, histogram, summary) and exposes them at the `/metrics` endpoint

A health check endpoint also exists at `/`

The following is a list of optional command line flags for configuration:
* `listen_address`: (default = `0.0.0.0:8080`)this defines the address and port that the sample app is exposed to. This is primarily to conform with the test framework requirements.
* `metric_count`: (default=1) the amount of each type of metric to generate. The same amount of metrics is always generated per metric type.

Steps for running locally:
```bash
$ go build .
$ ./prometheus-sample-app -listen_address=0.0.0.0:4567 -metric_count=100
```

Steps for running in docker:

```bash
$ docker build . -t prometheus-sample-app
$ docker run -it -p 8080:8080 prometheus-sample-app /bin/main -listen_address=0.0.0.0:8080
$ curl localhost:8080/metrics
```

Note that the port in LISTEN_ADDRESS must match the the second port specified in the port-forward

More functioning examples:

```bash
$ docker build . -t prometheus-sample-app
$ docker run -it -p 9001:8080 prometheus-sample-app /bin/main -listen_address=0.0.0.0:8080
$ curl localhost:9001/metrics
```

```bash
$ docker build . -t prometheus-sample-app
$ docker run -it -p 9001:8080 prometheus-sample-app /bin/main -listen_address=0.0.0.0:8080 -metric_count=100
$ curl localhost:9001/metrics
```

Running the commands above will require a config file for setting defaults. The config file is provided in this application. To modify it just change the values.
To override config file defaults you can specify your arguments via command line

Usage of generate:

  -is_random

    	Metrics specification

  -metric_count int

    	Amount of metrics to create

  -metric_frequency int

    	Refresh interval in seconds 

  -metric_type string
  
    	Type of metric (counter, gauge, histogram, summary) 

Example: 
```bash
$ docker build . -t prometheus-sample-app
$ docker run -it -p 8080:8080 prometheus-sample-app /bin/main -listen_address=0.0.0.0:8080 generate -metric_type=summary -metric_count=30 -metric_frequency=10
$ curl localhost:8080/metrics
```
```bash
$ docker build . -t prometheus-sample-app
$ docker run -it -p 8080:8080 prometheus-sample-app /bin/main -listen_address=0.0.0.0:8080 generate -metric_type=all -is_random=true
$ curl localhost:8080/metrics
```