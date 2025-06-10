#!/bin/bash

set -e

echo "🗑️  Destroying Terraform Docker deployment..."

if [ ! -d ".terraform" ]; then
    echo "❌ Terraform not initialized. Run 'terraform init' first."
    exit 1
fi

echo "📋 Planning destruction..."
terraform plan -destroy

read -p "⚠️  Are you sure you want to destroy all resources? (y/N): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "❌ Destruction cancelled."
    exit 1
fi

echo "🔥 Destroying Terraform deployment..."
terraform destroy -auto-approve

echo "🧹 Cleaning up remaining Docker resources..."
docker container prune -f
docker image prune -f
docker volume prune -f
docker network prune -f

echo ""
echo "✅ Cleanup complete!"
echo "🏁 All resources have been destroyed."