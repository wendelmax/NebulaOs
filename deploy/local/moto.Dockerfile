FROM python:3.12-alpine

RUN pip install --no-cache-dir moto[ec2,s3,server]

EXPOSE 4566

CMD ["moto_server", "-H", "0.0.0.0", "-p", "4566"]
