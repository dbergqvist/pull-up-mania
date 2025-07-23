#!/bin/bash

# Deploy to Google Cloud Functions
# Make sure you have gcloud CLI installed and authenticated

PROJECT_ID=${1:-"your-project-id"}
REGION=${2:-"us-central1"}

echo "Deploying Pull-Up Mania to Google Cloud Functions..."
echo "Project: $PROJECT_ID"
echo "Region: $REGION"

# Deploy the function
gcloud functions deploy pull-up-mania \
    --gen2 \
    --runtime=go121 \
    --region=$REGION \
    --source=. \
    --entry-point=PullUpMania \
    --trigger=http \
    --allow-unauthenticated \
    --memory=256MB \
    --timeout=60s \
    --max-instances=10 \
    --project=$PROJECT_ID

if [ $? -eq 0 ]; then
    echo ""
    echo "✅ Deployment successful!"
    echo "Your app is available at:"
    gcloud functions describe pull-up-mania --region=$REGION --project=$PROJECT_ID --format="value(serviceConfig.uri)"
else
    echo "❌ Deployment failed!"
    exit 1
fi 