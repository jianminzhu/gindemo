@rem go get github.com/go-bindata/go-bindata/...
@rem go get github.com/elazarl/go-bindata-assetfs/...

@rem go-bindata -debug -o=asset/asset.go -pkg=asset static/...

go-bindata -o=asset/asset.go -pkg=asset static/...