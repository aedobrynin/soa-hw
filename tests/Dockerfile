FROM python:3.10-alpine

WORKDIR /tests

COPY requirements.txt .
RUN pip install -r requirements.txt

COPY *.py .

ENTRYPOINT pytest -vs --log-level=DEBUG
