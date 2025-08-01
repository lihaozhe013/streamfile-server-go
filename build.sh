cd react
npm install
npm run build
cd ../go
go build

cd ..
mkdir -p dist
cp go/simple-server dist/
[ -f go/config.yaml ] && cp go/config.yaml dist/
cp go/config.yaml.example dist/
cp -r go/public dist/
