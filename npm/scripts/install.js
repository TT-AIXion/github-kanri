#!/usr/bin/env node
"use strict";

const crypto = require("crypto");
const fs = require("fs");
const https = require("https");
const os = require("os");
const path = require("path");
const { execFileSync } = require("child_process");

const pkg = require("../package.json");
const version = process.env.npm_package_version || pkg.version;
const tag = `v${version}`;

const platformMap = {
  darwin: "darwin",
  linux: "linux"
};
const archMap = {
  x64: "amd64",
  arm64: "arm64"
};

const platform = platformMap[process.platform];
const arch = archMap[process.arch];
if (!platform || !arch) {
  console.error(`unsupported platform: ${process.platform} ${process.arch}`);
  process.exit(1);
}

const artifact = `gkn_${version}_${platform}_${arch}.tar.gz`;
const baseUrl = `https://github.com/TT-AIXion/github-kanri/releases/download/${tag}`;
const tarUrl = `${baseUrl}/${artifact}`;
const checksumsUrl = `${baseUrl}/checksums.txt`;

const tmpDir = fs.mkdtempSync(path.join(os.tmpdir(), "gkn-"));
const tarPath = path.join(tmpDir, artifact);
const checksumsPath = path.join(tmpDir, "checksums.txt");
const extractDir = path.join(tmpDir, "extract");

const binDir = path.join(__dirname, "..", "bin");
const destBin = path.join(binDir, "gkn");

function download(url, dest, redirects = 0) {
  return new Promise((resolve, reject) => {
    const file = fs.createWriteStream(dest);
    https
      .get(url, { headers: { "User-Agent": "gkn-npm-installer" } }, (res) => {
        if (res.statusCode && [301, 302, 307, 308].includes(res.statusCode)) {
          if (redirects > 5) {
            reject(new Error("too many redirects"));
            res.resume();
            return;
          }
          const next = res.headers.location;
          if (!next) {
            reject(new Error("redirect missing location"));
            res.resume();
            return;
          }
          res.resume();
          download(next, dest, redirects + 1).then(resolve).catch(reject);
          return;
        }
        if (res.statusCode !== 200) {
          reject(new Error(`download failed: ${url} status=${res.statusCode}`));
          res.resume();
          return;
        }
        res.pipe(file);
        file.on("finish", () => file.close(resolve));
      })
      .on("error", reject);
    file.on("error", reject);
  });
}

function sha256File(filePath) {
  const hash = crypto.createHash("sha256");
  const data = fs.readFileSync(filePath);
  hash.update(data);
  return hash.digest("hex");
}

function getExpectedChecksum(text, filename) {
  const lines = text.split(/\r?\n/);
  for (const line of lines) {
    const parts = line.trim().split(/\s+/);
    if (parts.length >= 2 && parts[1] === filename) {
      return parts[0];
    }
  }
  return "";
}

async function main() {
  fs.mkdirSync(binDir, { recursive: true });
  fs.mkdirSync(extractDir, { recursive: true });

  await download(checksumsUrl, checksumsPath);
  const checksums = fs.readFileSync(checksumsPath, "utf8");
  const expected = getExpectedChecksum(checksums, artifact);
  if (!expected) {
    throw new Error("checksum not found for artifact");
  }

  await download(tarUrl, tarPath);
  const actual = sha256File(tarPath);
  if (actual !== expected) {
    throw new Error(`checksum mismatch: expected=${expected} actual=${actual}`);
  }

  execFileSync("tar", ["-xzf", tarPath, "-C", extractDir]);
  const extracted = path.join(extractDir, "gkn");
  if (!fs.existsSync(extracted)) {
    throw new Error("extracted binary not found");
  }
  fs.copyFileSync(extracted, destBin);
  fs.chmodSync(destBin, 0o755);
}

main().catch((err) => {
  console.error(`install failed: ${err.message}`);
  process.exit(1);
});
