#
# Dockerizing CentOS7: Dockerfile for building CentOS images
#
FROM centos:7
 
#设置一个时区的环境变量
ENV TZ "Asia/Shanghai"
 
#虚拟终端
ENV container docker
#dockerfile中有2条命令可以复制文件，1.COPY 2.ADD， ADD比COPY多2个功能，可以写成连接 直接COPY到container，如果是压缩文件，add能自动解压
#ADD aliyun-mirror.repo /etc/yum.repos.d/CentOS-Base.repo
#ADD aliyun-epel.repo /etc/yum.repos.d/epel.repo
 
RUN yum install -y curl wget tar bzip2 unzip vim-enhanced passwd sudo yum-utils hostname net-tools rsync man && \
    yum install -y gcc gcc-c++ git make automake cmake3 patch logrotate glibc-langpack-en.x86_64 git distcc golang&& \
    yum clean all;\
#创建admin用户
useradd -u0 -pcaveonline -o cave2;\
#修改admin用户密码
echo -e "caveonline\ncaveonline" | passwd cave2;
#USER cave2
WORKDIR /home/cave2 
RUN git clone https://github.com/boyxiaolong/mygccrpm.git;\
cd mygccrpm;\
rpm -iva gcc-5.4.0-1.el7.x86_64.rpm;\
rm -rf /home/cave2/mygccrpm
EXPOSE 22
VOLUME [ "/sys/fs/cgroup" ]
CMD ["/usr/sbin/init"]