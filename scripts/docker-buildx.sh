docker buildx build --platform linux/amd64 --build-arg GH_TOKEN=$GH_TOKEN --build-arg DOPPLER_TOKEN=$DOPPLER_TOKEN -t foodquest-api:$VERSION . && \
docker tag foodquest-api:$VERSION 061044801495.dkr.ecr.us-west-2.amazonaws.com/foodquest-api:$VERSION && \
docker push 061044801495.dkr.ecr.us-west-2.amazonaws.com/foodquest-api:$VERSION