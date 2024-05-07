version=`date +%y%m%d%H%M`
git pull
docker build -t octahub.8lab.cn:5000/oscro/octa-sourceserver:v${version} .
docker push octahub.8lab.cn:5000/oscro/octa-sourceserver:v${version}

echo "push Successful - v${version}"