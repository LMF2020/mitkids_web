#!/bin/bash
cd /opt/nginxdocker/mulkids-cms-pro
yarnbash="yarn run build"
if [ "$1" == "all" ]; then
  yarnbash="yarn && yarn run build"
fi
git pull
cd muitkid-cms && eval ${yarnbash}
rm -Rf /opt/nginxdocker/muitkid-web-app/mkcms
cp -Rf /opt/nginxdocker/mulkids-cms-pro/muitkid-cms/dist /opt/nginxdocker/muitkid-web-app/mkcms

cd ..
cd muitkid-stu/ && eval ${yarnbash}
rm -Rf /opt/nginxdocker/muitkid-web-app/mkstu
cp -Rf /opt/nginxdocker/mulkids-cms-pro/muitkid-stu/dist /opt/nginxdocker/muitkid-web-app/mkstu

cd ..
cd muitkid-tea/ && eval ${yarnbash}
rm -Rf /opt/nginxdocker/muitkid-web-app/mktea
cp -Rf /opt/nginxdocker/mulkids-cms-pro/muitkid-tea/dist /opt/nginxdocker/muitkid-web-app/mktea

cd ..
cd muitkid-portal/ && yarn run build
rm -Rf /opt/nginxdocker/muitkid-web-app/muitkid-portal
cp -Rf /opt/nginxdocker/mulkids-cms-pro/muitkid-portal/dist /opt/nginxdocker/muitkid-web-app/muitkid-portal