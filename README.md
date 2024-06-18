# RAGO: a GO-based RAG project

## 🦙 Base Model
RAGO leverages the Ollama API to enhance retrieval-augmented generation capabilities. Ensure you have Ollama installed. Pull the chat and embed models of your choice, then update the `model/modelConfig.json` file accordingly. 

## 🏠 Vector Store
RAGO requires a vector store to save and index document embeddings. You have two options depending on your usage needs:
- **MatrixOne**: Recommended for scenarios where documents are accessed repeatedly. Installation of MatrixOne is required.
- **FAISS**: Suitable for one-time usage. Run the FAISS service with `python faiss/serve.py`.

## 📜 Document Preparation
Place the documents you intend to retrieve in the `docs` directory. Currently, RAGO only supports documents in `.txt` format.

## ✨ Running RAGO
Execute RAGO using the following command, specifying the file and repository as required:
```
cd cmd
go run main.go --file pride-and-prejudice.txt --repo MatrixOne
```
