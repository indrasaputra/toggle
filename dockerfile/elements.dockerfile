FROM nginx

COPY ../index.html /usr/share/nginx/html/index.html
COPY ../openapiv2/proto/indrasaputra/toggle/v1/toggle.swagger.json /usr/share/nginx/html/api/toggle.swagger.json
