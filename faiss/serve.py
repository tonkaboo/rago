import faiss
import pika
import json
import numpy as np
from concurrent import futures
import grpc
from proto import faiss_pb2, faiss_pb2_grpc
import threading

chunk_contents = [] # stores chunks' contents
index = faiss.IndexFlatL2(1024)

def callback(ch, method, properties, body):
    print("[x] Received a chunk")
    chunk = json.loads(body)
    chunk_contents.append(chunk["Content"])
    numpy_array = np.array(chunk["Embedding"], dtype=np.float32)
    index.add(numpy_array.reshape(1, 1024))

    print("[x] Sent an ack")
    ch.basic_publish(exchange="", routing_key="acks", body="ack")


def connect_mq():
    # Set the connection and communication channel
    connection = pika.BlockingConnection(pika.ConnectionParameters(host="localhost"))
    channel = connection.channel()

    # Declare the queue
    channel.queue_declare(queue="chunks", durable=True)
    channel.queue_declare(queue="acks", durable=True)

    # Start to consume the messages
    channel.basic_consume(queue="chunks", on_message_callback=callback, auto_ack=False)
    print("Waiting for messages. To exit press CTRL+C")
    channel.start_consuming()


class FaissService(faiss_pb2_grpc.FaissServiceServicer):
    def FindSimilarChunks(self, request, context):
        embedding = np.array(request.embedding, dtype=np.float32).reshape(1, 1024)
        # Perform the search
        _, indices = index.search(embedding, k=3)  # Let's assume we want the top 3 nearest neighbors
        for idx in indices[0]:
            print(idx)
        results = [chunk_contents[idx] for idx in indices[0]]
        print(results)
        return faiss_pb2.SimilarityReply(chunk_contents=results)


def start_rpc_service():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    faiss_pb2_grpc.add_FaissServiceServicer_to_server(FaissService(), server)
    server.add_insecure_port('[::]:50051')
    server.start()
    server.wait_for_termination()


def main():
    mq_thread = threading.Thread(target=connect_mq)
    rpc_thread = threading.Thread(target=start_rpc_service)

    mq_thread.start()
    rpc_thread.start()

    mq_thread.join()
    rpc_thread.join()

if __name__ == "__main__":
    main()