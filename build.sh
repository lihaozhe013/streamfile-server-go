cd react
npm install
npm run build
rm -r ../public/markdown-viewer
mkdir -p ../public/markdown-viewer
cp -r dist/* ../public/markdown-viewer/

cd ..
go build
mkdir -p dist
[ -f simple-server ] && cp simple-server dist/
[ -f simple-server.exe ] && cp simple-server.exe dist/
[ -f config.yaml ] && cp config.yaml dist/
[ -f config.yaml.example ] && cp config.yaml.example dist/
cp -r public dist/
