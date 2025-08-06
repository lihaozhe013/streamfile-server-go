cd react
npm install
npm run build


cd ..
go build
mkdir -p dist
cp simple-server dist/
[ -f config.yaml ] && cp config.yaml dist/
cp config.yaml.example dist/
cp -r public dist/
