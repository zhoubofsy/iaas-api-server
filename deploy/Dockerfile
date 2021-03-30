FROM reg.iaas/baseimg/centos8.3:v1
COPY iaas-api-server /home/app/iaas-api-server
COPY conf /home/app/conf
COPY entrypoint.sh /home/app/entrypoint.sh
USER root
RUN mkdir -p /home/app/log && chmod 755 /home/app/iaas-api-server /home/app/entrypoint.sh /home/app/log
USER app
ENTRYPOINT [ "/home/app/entrypoint.sh" ]
