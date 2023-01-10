docker system prune -f

docker login registry1.dso.mil

docker build -t watchman-prd . --no-cache

# Get repo URI from WATCHMAN_REPO in .env
WATCHMAN_REPO=$(grep WATCHMAN_REPO .env | cut -d '=' -f 2)
AWS_REGION=$(grep AWS_REGION .env | cut -d '=' -f 2)

aws ecr get-login-password --region $AWS_REGION | docker login --username AWS --password-stdin $WATCHMAN_REPO

docker tag watchman-prd:latest $WATCHMAN_REPO/watchman-prd:latest

docker push $WATCHMAN_REPO/watchman-prd:latest
