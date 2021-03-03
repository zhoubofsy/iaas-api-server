FROM reg.iaas/baseimg/centos7.4:v1
COPY iaas-api-server /home/app/iaas-api-server
COPY conf /home/app/conf
#COPY entrypoint.sh /home/app/entrypoint.sh
USER root
RUN chmod 755 /home/app/iaas-api-server
USER app
ENTRYPOINT [ "/home/app/iaas-api-server", "--conf", "/home/app/conf/iaas.conf" ]
