import os


class Config:
    MQ_URL: str = os.getenv("MQ_URL", "amqp://rabbitmq:rabbitmq@localhost:5672")
