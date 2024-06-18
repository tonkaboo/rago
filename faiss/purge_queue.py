import pika

def connect_rabbitmq():
    """连接到 RabbitMQ 服务器"""
    return pika.BlockingConnection(pika.ConnectionParameters('localhost'))

def get_queue_message_count(channel, queue_name):
    """获取指定队列的消息数量"""
    queue = channel.queue_declare(queue=queue_name, passive=True)
    return queue.method.message_count

def purge_queue(channel, queue_name):
    """清空指定队列"""
    channel.queue_purge(queue=queue_name)

def main():
    connection = connect_rabbitmq()
    channel = connection.channel()

    queues = ["acks", "chunks"]
    for q in queues:
    # 检查清空前的消息数量
        message_count_before = get_queue_message_count(channel, q)
        print(f"Messages in queue {q} before purge: {message_count_before}")

        # 清空队列
        purge_queue(channel, q)

        # 检查清空后的消息数量
        message_count_after = get_queue_message_count(channel, q)
        print(f"Messages in queue {q} after purge: {message_count_after}")

    connection.close()

if __name__ == "__main__":
    main()
