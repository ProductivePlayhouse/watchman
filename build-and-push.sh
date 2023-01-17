# Retrieve values from .env file
# NOTE: .env file must already exist
CONTAINER_REPO=$(grep CONTAINER_REPO .env | cut -d '=' -f 2)
AWS_REGION=$(grep AWS_REGION .env | cut -d '=' -f 2)

docker system prune -f

docker build -t watchman-prd . --no-cache --build-arg CONTAINER_REPO=$CONTAINER_REPO

aws ecr get-login-password --region $AWS_REGION | docker login --username AWS --password-stdin $CONTAINER_REPO

docker tag watchman-prd:latest $CONTAINER_REPO/watchman-prd:latest

docker push $CONTAINER_REPO/watchman-prd:latest
