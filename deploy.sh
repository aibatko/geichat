#!/bin/bash

set -e

echo "🚀 Starting Terraform Docker deployment..."

if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker and try again."
    exit 1
fi

if ! command -v terraform &> /dev/null; then
    echo "❌ Terraform is not installed. Please install Terraform and try again."
    exit 1
fi

echo "📦 Initializing Terraform..."
terraform init

echo "✅ Validating Terraform configuration..."
terraform validate

#echo "📋 Planning Terraform deployment..."
#terraform plan

echo "🔧 Applying Terraform deployment..."
terraform apply -auto-approve

echo ""
echo "🎉 Deployment completed successfully!"
echo ""
echo "📱 Application URLs:"
echo "   Frontend: http://localhost:3000"
echo "   Backend:  http://localhost:8080"
echo "   Database: localhost:5432"
echo ""
echo "📊 Check deployment status:"
echo "   terraform show"
echo ""
echo "🔄 To update the deployment:"
echo "   terraform apply"
echo ""
echo "🗑️  To destroy the deployment:"
echo "   terraform destroy"