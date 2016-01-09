#
# Official docker file for building HuntJS powered applications
#

FROM fedora:23

# Upgrade dependencies
RUN dnf upgrade -y

# Clear cache
RUN dnf clean all

# Listen on 22 port
ENV GOSSHA_PORT=22
EXPOSE 22

# Create home directory
ENV GOSSHA_HOMEDIR=/root/.gossha

# Inject code of your application
ADD build/gossha /usr/bin/gossha

# Create root user
RUN /usr/bin/gossha passwd root root

# Create first ordinary user
RUN /usr/bin/gossha passwd user1 user1

# Create second ordinary user
RUN /usr/bin/gossha passwd user2 user2

# Inject SSHD server keys
ADD build/id_rsa /root/.ssh/

# Inject SSHD server keys
ADD build/id_rsa.pub /root/.ssh/

# Run the image process. Point second argument to your entry point of application
CMD ["/usr/bin/gossha"]