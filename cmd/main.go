package main

import (
	"bufio"
	"flag"
	"fmt"
	"gems/pkg/chunk"
	"gems/pkg/chunks"
	"gems/pkg/message"
	"gems/pkg/model"
	"gems/pkg/repo"
	"log"
	"os"
)

func main() {
	// instantiate a customized model
	myModel, err := model.NewModel()
	if err != nil {
		log.Fatal("Error creating model: ", err)
	}

	// the file flag determines if a new document would be processed
	fileName := flag.String("file", "", "The document to be indexed.")
	// the repo flag indicates the underlying embedding storage
	repoType := flag.String("repo", "MatrixOne", "The embedding storage to be used. You may choose from FAISS and MatrixOne.")
	flag.Parse()

	// connect to the database
	myRepo, err := repo.NewRepo(*repoType, *fileName)
	if err != nil {
		log.Fatal("Error creating the repo: ", err)
	}
	defer myRepo.CloseConn()

	if !myRepo.CheckIndex() {
		// read the document you want to query about
		docChunks, err := chunk.SplitDoc(*fileName, 10)
		if err != nil {
			log.Fatal("Error spliting the document: ", err)
		}

		// embed the chunks and store them into the database
		myChunks := chunks.Chunks(docChunks)
		myChunks.EmbedAndInsert(myModel, myRepo)
	}

	// embed your question
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("🐱 Please input your question (input 'quit' to exit the program):")
	for scanner.Scan() {
		myQuestion := scanner.Text()
		if myQuestion == "quit" {
			break
		}

		embeddingQuestion, err := myModel.Embed(myQuestion)
		if err != nil {
			log.Fatal("Error embedding the question: ", err)
		}
		relatedChunks, err := myRepo.RetrieveRelatedChunks(embeddingQuestion)
		if err != nil {
			log.Fatal("Error finding related chunks: ", err)
		}

		context := "Context:\n"
		for i, chunk := range relatedChunks {
			context += fmt.Sprintf("\nContext%s:\n%s", fmt.Sprint(i), chunk)
		}

		messages := make([]message.Message, 0)
		systemMessage := message.Message{Role: "system", Content: "You are a student taking an literature exam. Answer the question briefly based on the given context."}
		myMessage := message.Message{Role: "user", Content: fmt.Sprintf("%s\nQuestion:\n%s", context, myQuestion)}
		messages = append(messages, systemMessage)
		messages = append(messages, myMessage)
		response, err := myModel.Chat(messages)
		if err != nil {
			log.Fatal("Error generating response: ", err)
		}
		fmt.Println("🤖 The answer is: ")
		fmt.Println(response)

		fmt.Println("🐱 Please input your questions (input 'quit' to exit the program):")
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("Error reading the questions: ", err)
	}
}
