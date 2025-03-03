FROM python:3

COPY . /home/app
WORKDIR /home/app/

RUN pip install -r requirements.txt

EXPOSE 3000

CMD python3 manage.py runserver 0.0.0.0:3000