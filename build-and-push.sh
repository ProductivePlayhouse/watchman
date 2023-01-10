docker build -t watchman-prd . --no-cache

# Get repo URI from WATCHMAN_REPO in .env
repo=$(grep WATCHMAN_REPO .env | cut -d '=' -f 2)
region=$(grep AWS_REGION .env | cut -d '=' -f 2)

aws ecr get-login-password --region $region | docker login --username AWS --password-stdin $repo

docker tag watchman-prd:latest $repo/watchman-prd:latest

docker push $repo/watchman-prd:latest
