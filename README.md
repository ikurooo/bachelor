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

go version go1.27.1 linux/amd64
tinygo version 0.42.0 linux/amd64 (using go version go1.27.1 and LLVM version 22.1.4)
wasm-tools 1.260.0
wasmtime 49.0.0 (17830bd3c 2026-09-21)
containerd github.com/containerd/containerd/v2 2.2.1
-rwxr-xr-x 1 user user 45M Sep 27 20:29 /usr/local/bin/containerd-shim-wasmtime-v1

git clone https://github.com/containerd/runwasi.git
cd runwasi
make build-wasmtime
sudo make install
