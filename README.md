
## Threat Cloud reputation API wrapper

API service to consume Threat Cloud reputation databases, but easier to use compared to official API.


### Pre-requisites

You need to obtain ThreatCloud reputation service API key from Check Point employee. We assume it will be set in env variable as below.

```
export TC_API_KEY=bring_your_own_api_key
```


### Usage

```
# start service in background using Docker
$ docker run --name tcapi -d -p 8080:8080 -e PORT=8080 -e TC_API_KEY=$TC_API_KEY chkpmkoldov/threatcloudapi
8d37d497e6a85692c54f82f593967413d6958f36abbb56db88c040b86c244c1d

# check usage instructions
$ docker logs tcapi
Listening on port 8080
 usage http://127.0.0.1:8080/query/resource/(ip|hash|domain)

# also homepage has usage instructions
$ curl localhost:8080
Usage: /query/resource/(ip|hash|domain)

# your first query
$ curl localhost:8080/query/resource/1.1.1.1
{"response":[{"status":{"code":2001,"label":"SUCCESS","message":"Succeeded to generate reputation"},"resource":"1.1.1.1","reputation":{"classification":"Benign","severity":"N/A","confidence":"Low"},"risk":0,"context":{"location":{"countryCode":"AU","countryName":"Australia","region":"07","city":"Research","postalCode":"3095","latitude":-37.699997,"longitude":145.18329,"dma_code":0,"area_code":0,"metro_code":0},"asn":13335,"as_owner":"Cloudflare Inc"}}]}
```


### Docker image

Image is published in Docker Hub as [chkpmkoldov/threatcloudapi](https://hub.docker.com/repository/docker/chkpmkoldov/threatcloudapi)