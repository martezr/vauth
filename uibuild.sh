cd frontend
yarn build
rm -Rf ../cmd/webui
mv dist ../cmd/webui
cd ..
go build