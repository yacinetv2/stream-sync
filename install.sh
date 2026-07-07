#!/bin/bash
set -e

echo "=== Stream Sync Installer ==="

read -p "Access Key: " ACCESS_KEY
read -p "Secret Key: " SECRET_KEY
read -p "Endpoint: " ENDPOINT
read -p "Bucket: " BUCKET
read -p "Public URL مثل https://cdn.ycncdn.online: " PUBLIC_URL
read -p "Watch Path مثل /home/pankaj-dev/stream/tmp/live: " WATCH_PATH
read -p "Workers [64]: " WORKERS
WORKERS=${WORKERS:-64}

PROJECT_DIR="/root/stream-sync"
REPO_URL="https://github.com/yacinetv2/stream-sync.git"

apt update
apt install -y curl wget git build-essential

rm -rf /usr/local/go
wget -q https://go.dev/dl/go1.24.6.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.24.6.linux-amd64.tar.gz
rm -f go1.24.6.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
echo 'export PATH=$PATH:/usr/local/go/bin' >> /root/.bashrc

rm -f /usr/bin/rclone
curl https://rclone.org/install.sh | bash

mkdir -p /root/.config/rclone

cat > /root/.config/rclone/rclone.conf <<EOF
[r2]
type = s3
provider = Cloudflare
access_key_id = $ACCESS_KEY
secret_access_key = $SECRET_KEY
endpoint = $ENDPOINT
region = auto
acl =
env_auth = false
EOF

rm -rf "$PROJECT_DIR"
git clone "$REPO_URL" "$PROJECT_DIR"

cat > "$PROJECT_DIR/config.yaml" <<EOF
workers: $WORKERS
watch_path: "$WATCH_PATH"

r2:
  endpoint: "$ENDPOINT"
  bucket: "$BUCKET"
  access_key: "$ACCESS_KEY"
  secret_key: "$SECRET_KEY"
  public_url: "$PUBLIC_URL"
EOF

cd "$PROJECT_DIR"
go mod tidy
go build -o stream-sync ./cmd/stream-sync
chmod +x stream-sync

cat > /etc/systemd/system/stream-sync.service <<EOF
[Unit]
Description=Stream Sync R2 Uploader
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
User=root
WorkingDirectory=$PROJECT_DIR
ExecStart=$PROJECT_DIR/stream-sync
Restart=always
RestartSec=2
LimitNOFILE=100000

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl enable stream-sync
systemctl restart stream-sync

echo "=== Installed successfully ==="
systemctl status stream-sync --no-pager
