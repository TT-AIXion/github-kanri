#!/usr/bin/env node
"use strict";

const fs = require("fs");
const path = require("path");
const { spawn } = require("child_process");

const binPath = path.join(__dirname, "gkn");
if (!fs.existsSync(binPath)) {
  console.error("gkn binary not found. Reinstall with npm to fetch it.");
  process.exit(1);
}

const child = spawn(binPath, process.argv.slice(2), { stdio: "inherit" });
child.on("exit", (code, signal) => {
  if (signal) {
    process.exit(1);
  }
  process.exit(code == null ? 1 : code);
});
child.on("error", (err) => {
  console.error(`failed to run gkn: ${err.message}`);
  process.exit(1);
});
