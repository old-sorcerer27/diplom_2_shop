
FROM node:alpine
WORKDIR /usr/app/diplom
EXPOSE 3000
COPY ./ ./
WORKDIR /usr/app/diplom/frontend
RUN npm install
CMD ["npm", "start"]