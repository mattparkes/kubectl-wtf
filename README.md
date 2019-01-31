# kubectl-wtf

## Overview:

Check for common problems preventing a Kubernetes workload from working.

_Note: This is very much a work in progress and it doesn't (yet) actually talk to the Kubernetes API._  
_I'm also using this project to learn golang, so excuse how awful the code/repo is currently laid out (a bunch of functions in `main.go`). I'll gladly accept any help making things more idiomatic._

## Installation:

### kubectl >=1.12.0
1. Compile and add to your `PATH` as a binary named `kubectl-wtf` (or something more politically correct such as `kubectl-diagnose`)
2. `kubectl wtf ingress myingress`

### kubectl >= 1.10.0
1. Add to `~/.kube/plugins/wtf/`
2. `kubectl plugin wtf ingress myingress`

## Usage

```
$ kubectl wtf ingress test

Checking Ingress 'test' in namespace 'default' to ensure it exists:
  [OK]: Ingress test exists in namespace default:
    Hostnames: [mattparkes.net notreal.host]
    Certificates: [mattparkes.net]
    Paths: [/]
    Backend Services: [test]
    Backend Pods: []
  [Error]: Service 'test' points to no Running Pods
  [Error]: Service 'test' points to no Pods in any state

Checking local hosts file for [mattparkes.net notreal.host]:
  [Warn]: Local hosts file entry for 'mattparkes.net' found:
    line 5:    205.134.241.102  mattparkes.net
  [OK]: No local hosts file entry found for 'notreal.host'

Checking DNS resolution for [mattparkes.net notreal.host]:
  [OK]: Hostname 'mattparkes.net' resolves to: [205.134.241.102]
  [Error]: Unable to resolve Hostname 'notreal.host'

Checking TCP Ports [80 443] for [mattparkes.net notreal.host]:
  [OK]: TCP connection to mattparkes.net:80 successfully established
  [OK]: TCP connection to mattparkes.net:443 successfully established
  [Error]: TCP connection to notreal.host:80 could not be established
    dial tcp: lookup notreal.host on 192.168.0.1:53: no such host
  [Error]: TCP connection to notreal.host:443 could not be established
    dial tcp: lookup notreal.host on 192.168.0.1:53: no such host

TODO Checking SSL Certificates for [mattparkes.net notreal.host]:
    [TODO]: mattparkes.net
    [TODO]: notreal.host

Skipped: Checking Pods []:
```