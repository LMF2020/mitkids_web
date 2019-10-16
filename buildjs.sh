#!/bin/bash
cd /opt/nginxdocker/mulkids-cms-pro
git pull
cd muitkid-cms && yarn run build 
#rm -Rf /opt/nginxdocker/muitkid-web-app/mkcms
#cp -Rf /opt/nginxdocker/mulkids-cms-pro/muitkid-cms/dist /opt/nginxdocker/muitkid-web-app/mkcms
cd ..
cd muitkid-stu/ && yarn run build 
cd ..
cd muitkid-tea/ && yarn run build
cd ..
cd muitkid-portal/ && yarn run build

