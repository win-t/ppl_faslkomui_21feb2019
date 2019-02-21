FROM nginx:alpine
ADD /docker-entrypoint.sh /
ENTRYPOINT ["/docker-entrypoint.sh"]
ADD /src/ /usr/share/nginx/html/
