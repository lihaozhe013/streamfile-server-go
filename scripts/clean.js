const fs = require("fs-extra");
const path = require("path");
const root = path.resolve(__dirname, "..");
const r = (...p) => path.resolve(root, ...p);

try {
    fs.rmSync(r("simple-server"));
    console.log("removed /simple-server");
}
catch {}

try {
    fs.rmSync(r("simple-server.exe"));
    console.log("removed /simple-server.exe");
}
catch {}

try {
    fs.rmSync(r("public/markdown-viewer"), { recursive: true });
    console.log("removed public/markdown-viewer");
}
catch {}

try {
    fs.rmSync(r("public/styles.css"));
    console.log("removed public/styles.css");
}
catch {}

