#
# Dockerizing CentOS7: Dockerfile for building CentOS images
#
FROM centos:7
 
#设置一个时区的环境变量
ENV TZ = "Asia/Shanghai"\
LANG = zh_CN.UTF-8\
#虚拟终端
container = docker\
DISTCC_HOSTS="127.0.0.1 3"
#dockerfile中有2条命令可以复制文件，1.COPY 2.ADD， ADD比COPY多2个功能，可以写成连接 直接COPY到container，如果是压缩文件，add能自动解压
RUN yum install -y curl wget tar bzip2 unzip vim-enhanced passwd sudo yum-utils hostname net-tools rsync man && \
    yum install -y git make  logrotate glibc-langpack-en.x86_64 gcc gcc-c++ && \
    yum -y install epel-release && yum -y install distcc screen cmake3 python-pip python-devel mariadb-devel&&\
    pip install MySQL-python&&\
    yum clean all;\
#创建admin用户
useradd -u0 -pcaveonline -o cave2;\
#修改admin用户密码
echo -e "caveonline\ncaveonline" | passwd cave2;\
#USER cave2
#WORKDIR /home/cave2 
git clone https://github.com/boyxiaolong/mygccrpm.git;\
cd mygccrpm;\
rpm -iva gcc-5.4.0-1.el7.x86_64.rpm;
WORKDIR /usr/bin
RUN rm -rf cc c++;\
ln -s /usr/local/bin/gcc cc;\
ln -s /usr/local/bin/g++ c++;\
rm -rf /home/cave2/mygccrpm 
EXPOSE 22
VOLUME [ "/sys/fs/cgroup" ]
CMD ["/usr/sbin/init"]