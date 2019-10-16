#!/bin/bash
cd /opt/nginxdocker/mulkids-cms-pro
git pull
cd muitkid-cms && yarn run build 
rm -Rf /opt/nginxdocker/muitkid-web-app/mkcms
cp -Rf /opt/nginxdocker/mulkids-cms-pro/muitkid-cms/dist /opt/nginxdocker/muitkid-web-app/mkcms

cd ..
cd muitkid-stu/ && yarn run build 
rm -Rf /opt/nginxdocker/muitkid-web-app/mkstu
cp -Rf /opt/nginxdocker/mulkids-cms-pro/muitkid-stu/dist /opt/nginxdocker/muitkid-web-app/mkstu

cd ..
cd muitkid-tea/ && yarn run build
rm -Rf /opt/nginxdocker/muitkid-web-app/mktea
cp -Rf /opt/nginxdocker/mulkids-cms-pro/muitkid-tea/dist /opt/nginxdocker/muitkid-web-app/mktea

