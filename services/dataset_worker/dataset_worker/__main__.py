import json
from contextlib import contextmanager

import pika
from config import Config
from dotenv import load_dotenv


@contextmanager
def open_mq_connection(url: str):
    conn = pika.BlockingConnection(pika.URLParameters(url))
    try:
        yield conn
    finally:
        conn.close()


def callback(ch, method, properties, body):
    # 假设发送方发送的是 JSON 格式的数据
    data = json.loads(body)
    print("📦 Received dataset.create message:", data)

    # TODO: 在这里处理你的逻辑，比如写入数据库、文件等

    # 手动确认消息
    ch.basic_ack(delivery_tag=method.delivery_tag)


def main():
    load_dotenv()
    with open_mq_connection(Config.MQ_URL) as conn:
        channel = conn.channel()
        # 声明队列（确保存在，不会报错）
        channel.queue_declare(queue="dataset.create", durable=True)

        # 注册消费回调
        channel.basic_consume(
            queue="dataset.create",
            on_message_callback=callback,
            auto_ack=False,  # 手动 ack，确保可靠性
        )

        print(" [*] Waiting for messages in dataset.create. To exit press CTRL+C")
        channel.start_consuming()
        print(channel)


if __name__ == "__main__":
    main()


if __name__ == "__main__":
    main()
