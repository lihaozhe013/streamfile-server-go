const fs = require("fs-extra");
const path = require("path");
const root = path.resolve(__dirname, "..");
const r = (...p) => path.resolve(root, ...p);

// Rename dist to markdown-viewer
fs.renameSync(r("src/frontend/markdown-viewer/dist"), r("src/frontend/markdown-viewer/markdown-viewer"));

// Move to public directory
fs.moveSync(r("src/frontend/markdown-viewer/markdown-viewer"), r("public/markdown-viewer"), { overwrite: true });

console.log("Markdown viewer build artifacts moved to public/markdown-viewer");