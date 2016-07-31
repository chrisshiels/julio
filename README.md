# julio


## Overview

Explore blue / green microservice deployment in Rancher-managed Kubernetes.


## Notes - Rancher

- Rancher
  - It's a container orchestrator orchestrator.

  - Open source.

  - Handles deployment and management of container orchestration frameworks.

  - Supports multiple container orchestration frameworks:
    - Cattle
      - Rancher's native container orchestration framework.
      - Superset of Docker Compose and Swarm functionality.
    - Kubernetes
    - Mesos
    - Swarm

  - Manages multiple Rancher Environments
    - Environment is a collection of hosts running the same
      container orchestration framework.

  - "Rancher Monthly Training - June 2016"
    - https://www.youtube.com/watch?v=JqcpblO9OIs
      - 1.   Setting up Rancher.
      - 2.   Setting up Access Control.
      - 3.   Creating Environments.
      - 4.   Attaching private registries.
      - 5.   Adding resources / hosts.
      - 6.   UI, API and CLI access
      - 7.   Launching containers.
      - 8.   Creating stacks and services.
      - 9.   Rancher compose.
      - 10.  Load Balancers and Service Discovery.
      - 11.  Rancher Catalog.
      - 12.  Kubernetes / Swarm / Mesos.

  - http://docs.rancher.com/

  - See:
    - http://docs.rancher.com/rancher/latest/en/rancher-services/networking/
    - http://docs.rancher.com/rancher/latest/en/rancher-services/load-balancer/
    - http://docs.rancher.com/rancher/latest/en/rancher-services/dns-service/
    - http://docs.rancher.com/rancher/latest/en/rancher-services/metadata-service/
    - http://docs.rancher.com/rancher/latest/en/kubernetes/
    - http://docs.rancher.com/rancher/latest/en/kubernetes/ingress/
    - http://docs.rancher.com/rancher/latest/en/kubernetes/k8s-internal-dns-service/


## Preparation - Workstation NetworkManager / dnsmasq
```
rothko$ sudo tee /etc/NetworkManager/conf.d/julio.conf <<eof
[main]
dhcp=dhclient
dns=dnsmasq
eof

rothko$ sudo tee /etc/NetworkManager/dnsmasq.d/julio.conf <<eof
address=/.julio/192.168.40.13
host-record=vm1,192.168.40.11
host-record=vm2,192.168.40.12
host-record=vm3,192.168.40.13
host-record=vm4,192.168.40.14
eof

rothko$ sudo /bin/systemctl restart NetworkManager.service
rothko$ sudo /bin/systemctl status NetworkManager.service
```


Check:
```
rothko$ cat /etc/resolv.conf
# Generated by NetworkManager
search home.mecachis.net
nameserver 127.0.0.1

rothko$ dig +short hello.julio
192.168.40.13

rothko$ dig +short vm1
192.168.40.11

rothko$ dig +short vm2
192.168.40.12

rothko$ dig +short vm3
192.168.40.13

rothko$ dig +short vm4
192.168.40.14

rothko$ dig +short www.microsoft.com
www.microsoft.com-c-2.edgekey.net.
www.microsoft.com-c-2.edgekey.net.globalredir.akadns.net.
e2847.dspb.akamaiedge.net.
23.195.11.179
```


## Preparation - Workstation iptables

Needed to support Vagrant synced folders via nfs:
```
rothko$ sudo /sbin/iptables -I INPUT 1 -i virbr+ -j ACCEPT
```


## Virtual machines

rothko$ sudo -v
rothko$ vagrant up


## Build templater utility

The Kubernetes kubectl utility doesn't currently have any templating support
so the following utility is here as a workaround.
```
vm1$ sudo yum -y install golang

vm1$ cd /vagrant/templater
vm1$ make
```


## Install Docker
```
rothko$ vagrant ssh vm1

vm1234$ sudo yum -y install docker
vm1234$ sudo /bin/systemctl enable docker.service
vm1234$ sudo /bin/systemctl start docker.service
vm1234$ sudo /bin/systemctl status docker.service
```


## Configure Docker
```
vm1$ cd /vagrant/ssl
vm1$ make

vm1234$ sudo mkdir -p /etc/docker/certs.d/vm1:5000
vm1234$ sudo cp -i \
        /vagrant/ssl/vm1.crt \
        /etc/docker/certs.d/vm1:5000/ca.crt
vm1234$ sudo /bin/systemctl restart docker.service
vm1234$ sudo /bin/systemctl status docker.service
```


## Start registry
```
vm1$ sudo docker run -d -i -t \
        --name registry.docker --hostname registry.docker \
        -p 5000:5000 \
        -v /vagrant/ssl:/certs \
        -e REGISTRY_HTTP_TLS_CERTIFICATE=/certs/vm1.crt \
        -e REGISTRY_HTTP_TLS_KEY=/certs/vm1.key \
        registry:2

vm1$ sudo docker logs -f registry.docker

                - Have removed
                  -v /home/vagrant/registry:/var/lib/registry \
                  as was breaking docker push, not sure why.
```


## Install Rancher with Kubernetes
```
vm1$ sudo docker run \
	--name rancherserver.docker --hostname rancherserver.docker \
	-d --restart=always -p 8080:8080 rancher/server

vm1$ sudo docker logs -f rancherserver.docker
```


```
>> http://vm1:8080/

>> ADMIN >> Access Control
>> LOCAL
>> Login Username:    admin
   Full Name:         Admin
   Password:          admin
   Confirm Password:  admin
>> Enable Local Auth


>> Default >> Manage Environments
>> Add Environment
>> Kubernetes
>> Name:                     k8s
   Description:              k8s
   Virtual Machine Support:  Disabled
>> Create

>> Default >> k8s


>> Infrastructure >> Hosts
>> Add Host
>> "What base URL should hosts use to connect to the Rancher API?"
   This site's address:  http://vm1:8080
>> Save
>> Custom
   "sudo docker run -d --privileged -v /var/run/docker.sock:/var/run/docker.sock -v /var/lib/rancher:/var/lib/rancher rancher/agent:v1.0.2 http://vm1:8080/v1/scripts/7E0A24D8E4941C836B9B:1469966400000:INJcBMdDZ332NWuftKE5cTf8U8"
>> Close
```


```
vm2$ sudo docker run \
	-d --privileged \
	-v /var/run/docker.sock:/var/run/docker.sock \
	-v /var/lib/rancher:/var/lib/rancher \
	rancher/agent:v1.0.2 http://vm1:8080/v1/scripts/7E0A24D8E4941C836B9B:1469966400000:INJcBMdDZ332NWuftKE5cTf8U8

vm3$ sudo docker run \
	-d --privileged \
	-v /var/run/docker.sock:/var/run/docker.sock \
	-v /var/lib/rancher:/var/lib/rancher \
	rancher/agent:v1.0.2 http://vm1:8080/v1/scripts/7E0A24D8E4941C836B9B:1469966400000:INJcBMdDZ332NWuftKE5cTf8U8

vm4$ sudo docker run \
	-d --privileged \
	-v /var/run/docker.sock:/var/run/docker.sock \
	-v /var/lib/rancher:/var/lib/rancher \
	rancher/agent:v1.0.2 http://vm1:8080/v1/scripts/7E0A24D8E4941C836B9B:1469966400000:INJcBMdDZ332NWuftKE5cTf8U8
```


```
>> http://vm1:8080/
>> Default >> k8s
>> KUBERNETES

	- And wait... and wait and wait...
```


```
>> KUBERNETES >> Kubectl
>> Generate Config
```


```
vm1$ mkdir ~/.kube
vm1$ cat ~/.kube/config
apiVersion: v1
kind: Config
clusters:
- cluster:
    api-version: v1
    insecure-skip-tls-verify: true
    server: "https://vm1:8080/r/projects/1a7/kubernetes"
  name: "k8s"
contexts:
- context:
    cluster: "k8s"
    user: "k8s"
  name: "k8s"
current-context: "k8s"
users:
- name: "k8s"
  user:
    username: "3CE7AF813355BC6183D9"
    password: "VgzRimkiLGyizSi5cweUGmP18JczwderqFLdVdmF"

vm1$ sudo yum -y install kubernetes-client


vm1$ kubectl get nodes
NAME      STATUS    AGE
vm2       Ready     53m
vm3       Ready     31m
vm4       Ready     20m
```


Run something:
```
vm1$ kubectl run nginx --image=nginx --port=80

        - Wait...

vm1$ kubectl get pods -o wide
NAME                    READY     STATUS    RESTARTS   AGE       NODE
nginx-198147104-lxm6b   1/1       Running   0          6m        vm4

vm1$ kubectl expose deployment nginx --port=80 --type NodePort

vm1$ kubectl get svc nginx --template='{{(index .spec.ports 0).nodePort}}' \
        ; echo
31352

rothko$ curl http://vm4:31352/
```


Check DNS behaviour:
```
vm1$ kubectl run -i --tty chris --image=centos:7 --restart=Never -- /bin/bash
.
.
chris-v8fdz#
.
.
^p^q


vm1$ kubectl attach chris-v8fdz -i -t
.
.

chris-v8fdz# yum -y install bind-utils
chris-v8fdz# dig +short nginx.default.svc.cluster.local
10.43.204.184
^p^q
```


## Building and publishing microservices

```
vm1$ sudo yum -y install golang

vm1$ cd /vagrant/images/date/
vm1$ make clean ; make build VERSION=1.0.0
vm1$ sudo docker tag julio/date:1.0.0 vm1:5000/julio/date:1.0.0
vm1$ sudo docker push vm1:5000/julio/date:1.0.0

vm1$ cd /vagrant/images/time
vm1$ make clean ; make build VERSION=1.0.0
vm1$ sudo docker tag julio/time:1.0.0 vm1:5000/julio/time:1.0.0
vm1$ sudo docker push vm1:5000/julio/time:1.0.0

vm1$ cd /vagrant/images/web
vm1$ make clean ; make build VERSION=1.0.0
vm1$ sudo docker tag julio/web:1.0.0 vm1:5000/julio/web:1.0.0
vm1$ sudo docker push vm1:5000/julio/web:1.0.0

vm1$ sudo docker images
```


## Deploying microservices
```
vm1$ /vagrant/templater/templater \
        VERSION=1.0.0 \
        /vagrant/yaml/date.yaml.template | kubectl create -f /dev/stdin
You have exposed your service on an external port on all nodes in your
cluster.  If you want to expose this service to the external internet, you may
need to set up firewall rules for the service port(s) (tcp:30155) to serve traffic.

See http://releases.k8s.io/release-1.2/docs/user-guide/services-firewalls.md for more details.
service "date-1-0-0" created
replicationcontroller "date-1-0-0" created


vm1$ /vagrant/templater/templater \
        VERSION=1.0.0 \
        /vagrant/yaml/time.yaml.template | kubectl create -f /dev/stdin
You have exposed your service on an external port on all nodes in your
cluster.  If you want to expose this service to the external internet, you may
need to set up firewall rules for the service port(s) (tcp:30848) to serve traffic.

See http://releases.k8s.io/release-1.2/docs/user-guide/services-firewalls.md for more details.
service "time-1-0-0" created
replicationcontroller "time-1-0-0" created


vm1$ /vagrant/templater/templater \
        VERSION=1.0.0 \
        /vagrant/yaml/web.yaml.template | kubectl create -f /dev/stdin
You have exposed your service on an external port on all nodes in your
cluster.  If you want to expose this service to the external internet, you may
need to set up firewall rules for the service port(s) (tcp:30373) to serve traffic.

See http://releases.k8s.io/release-1.2/docs/user-guide/services-firewalls.md for more details.
service "web-1-0-0" created
replicationcontroller "web-1-0-0" created
```


Check:
```
vm1$ kubectl get svc
NAME         CLUSTER-IP      EXTERNAL-IP   PORT(S)    AGE
date-1-0-0   10.43.31.95     nodes         7001/TCP   3m
kubernetes   10.43.0.1       <none>        443/TCP    57m
nginx        10.43.204.184   nodes         80/TCP     31m
time-1-0-0   10.43.56.206    nodes         7002/TCP   3m
web-1-0-0    10.43.37.178    nodes         7000/TCP   3m

vm1$ kubectl get rc
NAME         DESIRED   CURRENT   AGE
date-1-0-0   3         3         3m
time-1-0-0   3         3         3m
web-1-0-0    3         3         3m

vm1$ kubectl get pods -o wide
NAME                    READY     STATUS    RESTARTS   AGE       NODE
chris-cuaz8             1/1       Running   0          29m       vm2
date-1-0-0-do2od        1/1       Running   0          3m        vm2
date-1-0-0-fm8fz        1/1       Running   0          3m        vm4
date-1-0-0-i8q52        1/1       Running   0          3m        vm3
nginx-198147104-lxm6b   1/1       Running   0          36m       vm4
time-1-0-0-2c856        1/1       Running   0          3m        vm3
time-1-0-0-hhnmh        1/1       Running   0          3m        vm4
time-1-0-0-nbi5e        1/1       Running   0          3m        vm2
web-1-0-0-7l0zs         1/1       Running   0          3m        vm3
web-1-0-0-9blh9         1/1       Running   0          3m        vm4
web-1-0-0-pznxo         1/1       Running   0          3m        vm2


vm1$ kubectl attach chris-v8fdz -i -t

chris-v8fdz# dig +short date-1-0-0.default.svc.cluster.local
10.43.31.95
chris-v8fdz# dig +short time-1-0-0.default.svc.cluster.local
10.43.56.206
chris-v8fdz# dig +short web-1-0-0.default.svc.cluster.local
10.43.37.178
^p^q
```


## Configure endpoint services
```
vm1$ /vagrant/templater/templater \
        VERSION=1.0.0 \
        /vagrant/yaml/date-endpoint.yaml.template | kubectl create -f /dev/stdin

vm1$ /vagrant/templater/templater \
        VERSION=1.0.0 \
        /vagrant/yaml/time-endpoint.yaml.template | kubectl create -f /dev/stdin

vm1$ /vagrant/templater/templater \
        VERSION=1.0.0 \
        /vagrant/yaml/web-endpoint.yaml.template | kubectl create -f /dev/stdin
```


## Configure ingress controller
```
vm1$ kubectl create -f /vagrant/yaml/web-ingress.yaml 
ingress "web" created
```


To view configuration:
```
>> KUBERNETES >> System
>> "kubernetes-ingress-lbs"
>> default-web
```


Check:
```
rothko$ curl http://vm2/
web-1-0-0-9h46l 1.0.0:
20160731 - date-1-0-0-1av8a 1.0.0
18:59:04 - time-1-0-0-6yeox 1.0.0
rothko$ curl http://vm3/
web-1-0-0-9h46l 1.0.0:
20160731 - date-1-0-0-1av8a 1.0.0
18:59:10 - time-1-0-0-6yeox 1.0.0
rothko$ curl http://vm4/
web-1-0-0-9h46l 1.0.0:
20160731 - date-1-0-0-1av8a 1.0.0
18:59:15 - time-1-0-0-6yeox 1.0.0
```


## Building and publishing new version of date microservice
```
vm1$ cd /vagrant/images/date/
vm1$ rm -f build bin/date ; make build VERSION=1.0.1
vm1$ sudo docker tag julio/date:1.0.1 vm1:5000/julio/date:1.0.1
vm1$ sudo docker push vm1:5000/julio/date:1.0.1
```


## Deploying new version of date microservice
```
vm1$ /vagrant/templater/templater \
        VERSION=1.0.1 \
        /vagrant/yaml/date.yaml.template | kubectl create -f /dev/stdin
You have exposed your service on an external port on all nodes in your
cluster.  If you want to expose this service to the external internet, you may
need to set up firewall rules for the service port(s) (tcp:32180) to serve traffic.

See http://releases.k8s.io/release-1.2/docs/user-guide/services-firewalls.md for more details.
service "date-1-0-1" created
replicationcontroller "date-1-0-1" created
```


## Update endpoint service
```
vm1$ kubectl delete svc date
vm1$ /vagrant/templater/templater \
        VERSION=1.0.1 \
        /vagrant/yaml/date-endpoint.yaml.template | kubectl create -f /dev/stdin
```

Check:
```
rothko$ curl http://vm2/
web-1-0-0-bn2x1 1.0.0:
20160731 - date-1-0-1-gy7nc 1.0.1
19:03:22 - time-1-0-0-6yeox 1.0.0
rothko$ curl http://vm3/
web-1-0-0-9h46l 1.0.0:
20160731 - date-1-0-1-6g2yb 1.0.1
19:03:23 - time-1-0-0-vvnyb 1.0.0
rothko$ curl http://vm4/
web-1-0-0-9h46l 1.0.0:
20160731 - date-1-0-1-6g2yb 1.0.1
19:03:25 - time-1-0-0-vvnyb 1.0.0
```


## To do
- Replace / update resources with kubectl.
- Confirm Rancher DNS (internal-dns - https://github.com/rancher/rancher-dns,
  external-dns - https://github.com/rancher/external-dns) doesn't support SRV.
- Test green deployments before making live.
  - Ideas:
    - Could do this with Ingress with host rules though would need to be updated
      for every service.
    - Could do this with Ingress in conjunection with proxy service which
      proxies to actual value of Host header.
