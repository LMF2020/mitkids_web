#!/bin/bash
cd /opt/nginxdocker/mulkids-cms-pro
yarnbash="yarn run build"
if [ "$1" == "all" ]; then
  yarnbash="yarn && yarn run build"
fi
git pull
rm -Rf /opt/nginxdocker/muitkid-web-app/*
cd muitkid-cms && eval ${yarnbash}
cp -Rf /opt/nginxdocker/mulkids-cms-pro/muitkid-cms/dist /opt/nginxdocker/muitkid-web-app/mkcms

cd ..
cd muitkid-stu/ && eval ${yarnbash}
cp -Rf /opt/nginxdocker/mulkids-cms-pro/muitkid-stu/dist /opt/nginxdocker/muitkid-web-app/mkstu

cd ..
cd muitkid-tea/ && eval ${yarnbash}
cp -Rf /opt/nginxdocker/mulkids-cms-pro/muitkid-tea/dist /opt/nginxdocker/muitkid-web-app/mktea

cd ..
cd muitkid-portal/ && yarn run build
cp -Rf /opt/nginxdocker/mulkids-cms-pro/muitkid-portal/dist /opt/nginxdocker/muitkid-web-app/muitkid-portal