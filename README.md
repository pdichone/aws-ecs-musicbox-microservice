## simple microservice example incl AWS XRay tracing

Microservice consisting of 3 different services:
* musicbox => frontend, where clients are talking to, distributes requests to either flamencoapp or operaapp based on request URI
  * flamencoapp => backend, returns a list of Flamenco artists
  * operaapp => backend, returns a list of Opera artists

Each service also has a ```/ping``` endpoint, for e.g. container healthchecks.

The list of artists to be returned can be specified by ENV VAR **ARTISTS=** for both, flamencoapp as well as operaapp.  
Also as ENV VAR you can specify the port on which the service listens on, **PORT=** (default values: musicboxapp: 9000, flamencoapp: 9001, operaapp: 9002)

## Containers

You can find the docker containers for the services on DockerHub:

* [musicboxapp container](https://hub.docker.com/repository/docker/pauloclouddev/musicboxapp)
* [flamencoapp container](https://hub.docker.com/repository/docker/pauloclouddev/flamencoapp)
* [operaapp container](https://hub.docker.com/repository/docker/pauloclouddev/operaapp)

## Usage

If you offer the musicboxapp directly to the internet, then just call ```http://<musicboxapp-public-ip>:9000/ping``` to check if it is responding.  
To request artists from the backend services, call
* ```http://<musicboxapp-public-ip>:9000/opera``` for a list of opera artists
* ```http://<musicboxapp-public-ip>:9000/flamenco``` for a list of flamenco artists

If you place a loadbalancer in front of the musicboxapp, the URL to talk to musicboxapp depends on the (port) forwarding rules of your loadbalancer. 
E.g. if your LB listens to port **80** and is forwarding traffic to the musicboxapp on port **9000**, then you'd call:
```http://<loadbalancer-public-dns>/ping``` to call the musicboxapp healthcheck endpoint, and
* ```http://<loadbalancer-public-dns>/flamenco``` to request a list of flamenco artists
* ```http://<loadbalancer-public-dns>/opera``` to request a list of opera artists
