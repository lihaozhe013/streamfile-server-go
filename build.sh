cd react
npm install
npm run build
rm -r ../public/markdown-viewer
mkdir -p ../public/markdown-viewer
cp -r dist/* ../public/markdown-viewer/

cd ..
go build
mkdir -p dist
cp simple-server dist/
[ -f config.yaml ] && cp config.yaml dist/
cp config.yaml.example dist/
cp -r public dist/
