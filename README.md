## Getting Start

### Before install

- Install golang
- Config gcloud
- Config goapp
- Config node.js with npm
- Install dep `go get -u github.com/golang/dep/cmd/dep`


### Install foreman and start

```
npm install -g foreman
cp .env.sample .env

# check configs by running server only
cd server
bin/update # update vendering
bin/run # run server
# access to http://localhost:9999
# if work fine, stop by Ctrl + C

cd ..
nf start # frontend and backend servers start
```

access to
http://localhost:5000
