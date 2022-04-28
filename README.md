# Tinyseed

## TinySeed is a fork of Binary Holding's Tenderseed, which is a fork of Polychain's Tenderseed

This tool runs a seed node for any tendermint based blockchain, crawls the network and generates a map with the geolocalisation of the peers.

It is used to pinpoint centralization of network in common infrastructure hosts like AWS, GCP etc, as well as to globally lanunch new chains from Miredo-enabled, edge-of-network write-through-cache powered raspberry pi devices that provide validation at a fraction of the original cost. 
 

### Configuration

```bash
git clone https://github.com/notional-labs/tinyseed
cd tinyseed
go mod tidy
go install .
tinyseed
```

Then you'll become a seed node on Osmosis-1. Let's do Cosmoshub-4, shall we? We've made Osmosis zeroconf, but hey this
thing here reads 2 env vars!

```bash
export ID=cosmoshub-4
export SEEDS=bf8328b66dceb4987e5cd94430af66045e59899f@public-seed.cosmos.vitwit.com:26656,cfd785a4224c7940e9a10f6c1ab24c343e923bec@164.68.107.188:26656,d72b3011ed46d783e369fdf8ae2055b99a1e5074@173.249.50.25:26656,ba3bacc714817218562f743178228f23678b2873@public-seed-node.cosmoshub.certus.one:26656,3c7cad4154967a294b3ba1cc752e40e8779640ad@84.201.128.115:26656,366ac852255c3ac8de17e11ae9ec814b8c68bddb@51.15.94.196:26656
tinyseed
```



## Docker

This streamlined program can be used with `docker` to easily deploy a seed node from within preconfigured `docker-compose.yml` file.

### Docker Setup
```bash
# Remove any currently installed Docker for clean install
sudo apt-get remove docker docker-engine docker.io containerd runc

# Download latest Package Information
sudo apt-get update

# Set Up Repository
sudo apt-get install \
    ca-certificates \
    curl \
    gnupg \
    lsb-release

# Get Docker Official GPG Key
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg

# Get Docker Stable
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu \
  $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

# Install Docker Engine
 sudo apt-get update
 sudo apt-get install docker-ce docker-ce-cli containerd.io docker-compose-plugin -y
 ```

Test to see if `Docker` is properly installed:
 ```bash
 sudo docker run hello-world
 ```
 
Run whichever chain you'd like in a docker container with this command below:

```bash
docker run -e ID=cosmoshub-4 -e SEEDS=bf8328b66dceb4987e5cd94430af66045e59899f@public-seed.cosmos.vitwit.com:26656,cfd785a4224c7940e9a10f6c1ab24c343e923bec@164.68.107.188:26656,d72b3011ed46d783e369fdf8ae2055b99a1e5074@173.249.50.25:26656,ba3bacc714817218562f743178228f23678b2873@public-seed-node.cosmoshub.certus.one:26656,3c7cad4154967a294b3ba1cc752e40e8779640ad@84.201.128.115:26656,366ac852255c3ac8de17e11ae9ec814b8c68bddb@51.15.94.196:26656 ghcr.io/notional-labs/tinyseed
```


## Akash

```yaml
---
version: "2.0"

services:
  cosmos:
    image: ghcr.io/notional-labs/tinyseed:latest
    env:
      - ID=cosmoshub-4
      - SEEDS=bf8328b66dceb4987e5cd94430af66045e59899f@public-seed.cosmos.vitwit.com:26656,cfd785a4224c7940e9a10f6c1ab24c343e923bec@164.68.107.188:26656,d72b3011ed46d783e369fdf8ae2055b99a1e5074@173.249.50.25:26656,ba3bacc714817218562f743178228f23678b2873@public-seed-node.cosmoshub.certus.one:26656,3c7cad4154967a294b3ba1cc752e40e8779640ad@84.201.128.115:26656,366ac852255c3ac8de17e11ae9ec814b8c68bddb@51.15.94.196:26656
    expose:
      - port: 8080
        as: 80
        to:
          - global: true
      - port: 6969
        to:
          - global: true
profiles:
  compute:
    cosmos:
      resources:
        cpu:
          units: 1
        memory:
          size: 200Mi
        storage:
          size: 200Mi
  placement:
    dcloud:
      pricing:
        cosmos:
          denom: uakt
          amount: 10
deployment:
  cosmos:
    dcloud:
      profile: cosmos
      count: 1
```

## Service File

If you'd like to run as service file, here is a template. TinySeed limits to 2 variables, displayed in example:

```
sudo nano /etc/systemd/system/tinyseed.service
```

```
[Unit]
Description=tinyseed
After=network-online.target

[Service]
User=<username>
ExecStart=/home/<username>/go/bin/tinyseed tinyseed
Restart=always
RestartSec=3
Environment="ID=<chain-id>"
Environment="SEEDS=<SEEDS>"
LimitNOFILE=4096

[Install]
WantedBy=multi-user.target
```

## License

[Blue Oak Model License 1.0.0](https://blueoakcouncil.org/license/1.0.0)
