
FROM node:alpine
WORKDIR /usr/app/diplom
EXPOSE 5173
COPY ./ ./
WORKDIR /usr/app/diplom/frontend
RUN npm install
CMD ["npm", "start"]