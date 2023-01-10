docker build -t watchman-prd .

# Get repo URI from WATCHMAN_REPO in .env
repo=$(grep WATCHMAN_REPO .env | cut -d '=' -f 2)

aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin $repo

docker tag watchman-prd:latest $repo/watchman-prd:latest

docker push $repo/watchman-prd:latest
