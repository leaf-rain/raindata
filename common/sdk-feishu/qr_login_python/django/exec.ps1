docker build -t p1 .
docker run --env-file .env -it -p 3000:3000 --name container1 p1