FROM hashicorp/terraform:light
MAINTAINER "Kerry Wilson <kwilson@goodercode.com>"

# Copy sources
ADD tfwatch /bin/tfwatch

ADD site/dist /opt/site

ENTRYPOINT [ "/bin/tfwatch", "--site-dir", "/opt/site", "/usr/src" ]