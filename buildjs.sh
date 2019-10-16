#!/bin/bash
cd /opt/nginxdocker/mulkids-cms-pro
yarnbash="yarn run build"
if [ $1 = "all" ]; then
  yarnbash="yarn && yarn run build"
fi
git pull

cd muitkid-cms  && eval ${yarnbash}
#rm -Rf /opt/nginxdocker/muitkid-web-app/mkcms
#cp -Rf /opt/nginxdocker/mulkids-cms-pro/muitkid-cms/dist /opt/nginxdocker/muitkid-web-app/mkcms
cd ..
cd muitkid-stu/  && eval ${yarnbash}
cd ..
cd muitkid-tea/  && eval ${yarnbash}
cd ..
cd muitkid-portal/  && eval ${yarnbash}


