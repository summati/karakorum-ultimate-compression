<h1 align="center">
  <img src="https://karakorum.onrender.com/assets/og-image.jpg" alt="Karakorum" width="100%">
  <br>
  Karakorum Pro CLI
</h1>

<p align="center">
  <b>The Ultimate Zero-Dependency, In-Memory Data Compression Engine.</b><br>
  <i>Engineered for maximum throughput and 100% cryptographic data integrity.</i>
</p>

<p align="center">
  <a href="https://karakorum.onrender.com/">🌐 Live Demo</a> •
  <a href="https://karakoram.gumroad.com/l/karakoram">💳 Get CLI License</a> •
  <a href="#-why-karakorum">Features</a>
</p>

---

## ⚡ Why Karakorum?

Standard ZIP tools rely heavily on disk I/O, creating massive bottlenecks when dealing with multi-gigabyte datasets. **Karakorum** was built from the ground up in pure Go to solve this. It streams, chunks, and compresses data **entirely in memory**, delivering blistering speeds.

* **Pure Speed:** No disk I/O bottlenecks. Engineered for enterprise environments.
* **100% Lossless Integrity:** Every byte is verified using SHA-256 checksums. Your data is exactly as you left it.
* **Zero Dependencies:** A single, lightweight static binary. No bloated runtimes or complex setups.
* **Cross-Platform:** Works natively across Linux, macOS, and Windows.

## 🚀 Quick Start

### Decompression (Free)
Decompressing `.kuc` files is completely free and requires no license.
```bash
./kuc decompress --file data.kuc
```

### Compression (Pro License Required)
To compress files, you need to activate your CLI using a license key. You can get yours from our [Gumroad Page](https://karakoram.gumroad.com/l/karakoram).
```bash
# 1. Activate your license
./kuc auth --key YOUR-GUMROAD-LICENSE-KEY

# 2. Compress your files instantly
./kuc compress --file massive_dataset.json
```

## 🌐 Web Interface & API

Not ready for the CLI? You can test the compression engine for free directly in your browser or integrate it into your backend using our REST API.

* **Interactive Web App:** [karakorum.onrender.com](https://karakorum.onrender.com/)
* **Developer API:**
  ```bash
  curl -X POST \
    https://karakorum.onrender.com/api/v1/compress \
    -F "file=@data.json"
  ```

---

<p align="center">
  Built with ❤️ by Summati
</p>
