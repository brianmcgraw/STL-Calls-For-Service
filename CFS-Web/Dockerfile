# build stage
FROM node:12.16.1-alpine As builder


# RUN npm ci && mkdir /bdm && mv ./node_modules ./bdm

WORKDIR /usr/src/app
COPY package.json package-lock.json ./
RUN npm install
COPY . .
# RUN npm rebuild node-sass
RUN npm run ng build --prod

# production stage

FROM nginx:stable-alpine as production-stage
RUN rm -rf /usr/share/nginx/html/*

COPY ./default.conf /etc/nginx/conf.d/default.conf
COPY --from=builder /usr/src/app/dist/CFS-Web/ /usr/share/nginx/html


EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]