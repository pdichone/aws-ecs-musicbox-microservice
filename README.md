## simple microservice example incl AWS XRay tracing

Microservice consisting of 3 different services:
* jukebox => frontend, where clients are talking to, distributes requests to either metalapp or popapp based on request URI
  * metalapp => backend, returns a list of Metal artists
  * popapp => backend, returns a list of Pop artists

Each service also has a ```/ping``` endpoint, for e.g. container healthchecks.

The list of artists to be returned can be specified by ENV VAR **ARTISTS=** for both, metalapp as well as popapp.  
Also as ENV VAR you can specify the port on which the service listens on, **PORT=** (default values: jukeboxapp: 9000, metalapp: 9001, popapp: 9002)

## Containers

You can find the docker containers for the services on DockerHub:

* [jukeboxapp container](https://hub.docker.com/r/gkoenig/jukeboxapp)
* [metalapp container](https://hub.docker.com/r/gkoenig/metalapp)
* [popapp container](https://hub.docker.com/r/gkoenig/popapp)

## Usage

If you offer the jukeboxapp directly to the internet, then just call ```http://<jukeboxapp-public-ip>:9000/ping``` to check if it is responding.  
To request artists from the backend services, call
* ```http://<jukeboxapp-public-ip>:9000/metal``` for a list of Metal artists
* ```http://<jukeboxapp-public-ip>:9000/pop``` for a list of Pop artists

If you place a loadbalancer in front of the jukeboxapp, the URL to talk to jukeboxapp depends on the (port) forwarding rules of your loadbalancer. 
E.g. if your LB listens to port **80** and is forwarding traffic to the jukeboxapp on port **9000**, then you'd call:
```http://<loadbalancer-public-dns>/ping``` to call the jukeboxapp healthcheck endpoint, and
* ```http://<loadbalancer-public-dns>/pop``` to request a list of Pop artists
* ```http://<loadbalancer-public-dns>/metal``` to request a list of Metal artists
