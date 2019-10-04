#!/bin/bash
cd /opt/nginxdocker/mulkids-cms-pro
yarnbash="yarn run build"
if [ $1 == "all" ]; then
  yarnbash="yarn &&yarn run build"
fi
git pull
cd muitkid-cms && ${yarnbash}
rm -Rf /opt/nginxdocker/muitkid-web-app/mkcms
cp -Rf /opt/nginxdocker/mulkids-cms-pro/muitkid-cms/dist /opt/nginxdocker/muitkid-web-app/mkcms

cd ..
cd muitkid-stu/ && ${yarnbash}
rm -Rf /opt/nginxdocker/muitkid-web-app/mkstu
cp -Rf /opt/nginxdocker/mulkids-cms-pro/muitkid-stu/dist /opt/nginxdocker/muitkid-web-app/mkstu

cd ..
cd muitkid-tea/ && ${yarnbash}
rm -Rf /opt/nginxdocker/muitkid-web-app/mktea
cp -Rf /opt/nginxdocker/mulkids-cms-pro/muitkid-tea/dist /opt/nginxdocker/muitkid-web-app/mktea

