"Ubuntu 24.04.5 LTS"
Kernel: Linux DESKTOP-JHRC3D7 6.18.40.1-microsoft-standard-WSL2 #1 SMP PREEMPT_DYNAMIC Fri Jul 31 22:12:15 UTC 2026 x86_64 x86_64 x86_64 GNU/Linux

WSL version: 3.0.1.0
Kernel version: 6.18.40.1-1
WSLg version: 1.0.79
MSRDC version: 1.2.7214
Direct3D version: 1.611.1-81528511
DXCore version: 10.0.26100.1-240331-1435.ge-release
Windows version: 10.0.19045.6466

Client:
Version:           29.1.3
API version:       1.52
Go version:        go1.24.4
Git commit:        29.1.3-0ubuntu3~24.04.2
Built:             Wed Apr 29 16:41:06 2026
OS/Arch:           linux/amd64
Context:           default

Server:
Engine:
Version:          29.1.3
API version:      1.52 (minimum version 1.44)
Go version:       go1.24.4
Git commit:       29.1.3-0ubuntu3~24.04.2
Built:            Wed Apr 29 16:41:06 2026
OS/Arch:          linux/amd64
Experimental:     false
containerd:
Version:          2.2.1
GitCommit:
runc:
Version:          1.3.4-0ubuntu1~24.04.1
GitCommit:
docker-init:
Version:          0.19.0
GitCommit:

github.com/docker/buildx v0.37.1 0b265a9f62db554fa9aba6dd19e1bd5704bc7d8a
NAME/NODE     DRIVER/ENDPOINT   STATUS    BUILDKIT   PLATFORMS
default*      docker
\_ default    \_ default       running   v0.26.2    linux/amd64 (+3)

k3s version v1.36.5+k3s1 (3dd98cc5)
go version go1.26.8

kubectl: Client Version: v1.36.5+k3s1
Kustomize Version: v5.8.1
Server Version: v1.36.5+k3s1

go version go 1.27.1 linux/amd64
tinygo version 0.42.0 linux/amd64 (using go version go1.27.1 and LLVM version 22.1.4)
wasm-tools 1.260.0
wasmtime 49.0.0 (17830bd3c 2026-09-21)
containerd github.com/containerd/containerd/v2 2.2.1
-rwxr-xr-x 1 user user 45M Sep 27 20:29 /usr/local/bin/containerd-shim-wasmtime-v1

# 0. (Crucial for WSL) Ensure time is synced and restart WSL if needed (--shutdown from Windows)
sudo timedatectl set-ntp true
sudo systemctl restart systemd-timesyncd
timedatectl

# Docker
sudo apt install -y docker.io
sudo usermod -aG docker $USER
sudo service docker start

# Rust
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
source "$HOME/.cargo/env"

# Build dependencies
sudo apt update
sudo apt install -y build-essential protobuf-compiler libseccomp-dev

# Build Wasmtime shim
git clone https://github.com/containerd/runwasi.git
cd ~/runwasi
make build-wasmtime

# Install the built shim
sudo install -m 0755 \
target/x86_64-unknown-linux-gnu/debug/containerd-shim-wasmtime-v1 \
/usr/local/bin/containerd-shim-wasmtime-v1

# Verify
ls -lh /usr/local/bin/containerd-shim-wasmtime-v1
command -v containerd-shim-wasmtime-v1

# Install K3s
curl -sfL https://get.k3s.io | sh -

# Install Go
## 1. Download the latest Go tarball (adjust version if needed)
curl -LO https://go.dev/dl/go1.27.1.linux-amd64.tar.gz

## 2. Extract it to /usr/local (removes any previous installation first)
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.27.1.linux-amd64.tar.gz

## 3. Add Go to your PATH (add this to ~/.bashrc or ~/.profile if you haven't already)
export PATH=$PATH:/usr/local/go/bin

## 4. Verify the installation
go version

cd /path/to/your/go-project

# Add componentize-go as a tool dependency to your module
go get -tool github.com/bytecodealliance/componentize-go@latest
go mod tidy

# Build the WebAssembly component
go tool componentize-go build -o app.wasm
docker build -t wasm-serverless:latest .
docker save wasm-serverless:latest -o /tmp/wasm-serverless.tar
sudo k3s ctr images import /tmp/wasm-serverless.tar
kubectl apply -f app.yaml
kubectl get pods -o wide
curl http://:8080

