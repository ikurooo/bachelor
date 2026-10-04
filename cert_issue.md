Troubleshooting & Common Gotchas
1. TLS Certificate Error (x509: certificate signed by unknown authority) After Reinstalling K3s

Every time K3s is uninstalled and reinstalled, it generates a brand-new internal Certificate Authority (CA). Because of this, your local ~/.kube/config file retains the old certificates, causing kubectl to fail with a TLS verification error.

The Fix

Refresh your user's kubeconfig file by copying the newly generated one and fixing its permissions:

# Point your environment to the correct config path
export KUBECONFIG=~/.kube/config

# Copy the fresh K3s config and update ownership/permissions
sudo cp /etc/rancher/k3s/k3s.yaml ~/.kube/config
sudo chown ivan:ivan ~/.kube/config
chmod 600 ~/.kube/config

# Verify it works
kubectl get pods -o wide


Tip: To make the KUBECONFIG variable permanent, add it to your shell profile:

echo 'export KUBECONFIG=~/.kube/config' >> ~/.bashrc

2. K3s Startup Failure / Certificate Error Due to WSL Time Drift

If your WSL clock drifts or desynchronizes from your Windows host, K3s can generate certificates with invalid time bounds, causing it to fail immediately upon startup.

The Fix: Sync Clock & Restart WSL

Run these commands inside WSL to sync the system time:

sudo timedatectl set-ntp true
sudo systemctl restart systemd-timesyncd
timedatectl


Then shut down your WSL instance from Windows PowerShell or Command Prompt — not from inside WSL:

wsl --shutdown


Re-open your WSL terminal and proceed with your K3s setup.