# Pinger (IN PROGRESS)

This program will ping specified hosts and notify you if some go down

# Installation

First install required packages

```
go get github.com/tatsushid/go-fastping
```

Then create remotes.txt in the root of project. The file should contain list of hosts (one per line).

# Usage

```
go build pinger.go
sudo ./pinger
```

`sudo` is required because of ICMP packages