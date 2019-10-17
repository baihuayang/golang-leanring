$GOPATH/bin/protoc --go_out=./proto/pbgo ./proto/*.proto 
cp -rf ./proto/pbgo/proto/* ./proto/pbgo
rm -rf ./proto/pbgo/proto