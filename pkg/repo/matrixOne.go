package repo

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"gems/pkg/chunk"
	"io"
	"os"
	"strings"

	"github.com/google/uuid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type matrixOne struct {
	tableName  string
	dbConn     *gorm.DB
	indexBuilt bool
}

// the default table name will be chunks
// but actually one repo corresponds to a table
// so we need to adjust the table name accordingly
type Chunk struct {
	Content   string
	Embedding string
}

// return the repo struct and tell the user if the specified repo already existed
func NewMatrixOne(fileName string) (*matrixOne, error) {
	// connect to matrix one database
	dsn := "root:111@tcp(127.0.0.1:6001)/test?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("an error occured when connecting to the database: %s", err)
	}
	fmt.Println("database connection succeed")

	// find the table corresponding to the file
	tableName, indexBuilt, err := findTableName(fileName, db)
	if err != nil {
		return nil, fmt.Errorf("failed to determine whether there is a table for the file: %s", err)
	}

	return &matrixOne{dbConn: db, tableName: tableName, indexBuilt: indexBuilt}, nil
}

type Meta struct {
	gorm.Model
	FileName  string
	Md5       string
	TableName string
}

// the first return value is the table name, the second id is whether the table existed long ago
func findTableName(fileName string, db *gorm.DB) (string, bool, error) {
	// open the doc and calculate the md5 hash
	file, err := os.Open(fmt.Sprintf("../../docs/%s", fileName))
	if err != nil {
		return "", false, fmt.Errorf("failed to open the file: %s", err)
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", false, fmt.Errorf("failed to hash the document: %s", err)
	}
	fileMd5 := fmt.Sprintf("%x", hash.Sum(nil))

	// check the meta table to see whether the doc has already had a corresponding table
	// note that the existence should be determined by the md5
	db.AutoMigrate(&Meta{})
	var fileMeta Meta
	result := db.First(&fileMeta, "md5 = ?", fileMd5)
	if result.Error != nil {
		// create a meta table record if none could be found, then create the corresponding table
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			fmt.Println("no existing table could be found")
			fileMeta = Meta{FileName: fileName, Md5: fileMd5, TableName: formatUuid(uuid.New())}
			db.Create(&fileMeta)

			createTableSQL := fmt.Sprintf("CREATE TABLE %s (content text, embedding vecf32(1024));", fileMeta.TableName)
			db.Exec(createTableSQL)
			fmt.Println("table created successfully")

			return fileMeta.TableName, false, nil
		} else {
			// if a record could be found, just return the corresponding table's name
			fmt.Println()
			return "", false, fmt.Errorf("an error occured when looking for meta data: %s", err)
		}
	}
	return fileMeta.TableName, true, nil
}

func formatUuid(u uuid.UUID) string {
	return strings.ReplaceAll(u.String(), "-", "")
}

func (mo *matrixOne) InsertChunk(appChunk *chunk.Chunk) error {
	strEmbedding, err := stringifyEmbedding(appChunk.Embedding)
	if err != nil {
		return fmt.Errorf("cannot convert the chunk embedding to string: %s", err)
	}

	mo.dbConn.Table(mo.tableName).AutoMigrate(&Chunk{})
	dbChunk := Chunk{
		Content:   appChunk.Content,
		Embedding: strEmbedding,
	}
	mo.dbConn.Table(mo.tableName).Create(&dbChunk)
	return nil
}

func (mo *matrixOne) RetrieveRelatedChunks(questionEmbedding []float32) ([]string, error) {
	strEmbedding, err := stringifyEmbedding(questionEmbedding)
	if err != nil {
		return nil, fmt.Errorf("failed to convert the embedding to a string: %s", err)
	}

	var relatedChunks []Chunk
	query := fmt.Sprintf("SELECT content FROM %s ORDER BY l2_distance(embedding, '%s') ASC LIMIT 3", mo.tableName, strEmbedding)
	mo.dbConn.Raw(query).Scan(&relatedChunks)

	var relatedContent []string
	for _, chunk := range relatedChunks {
		relatedContent = append(relatedContent, chunk.Content)
	}

	return relatedContent, nil
}

func stringifyEmbedding(embedding []float32) (string, error) {
	byteEmbedding, err := json.Marshal(embedding)
	if err != nil {
		return "", fmt.Errorf("failed to convert the embedding to bytes: %s", err)
	}
	strEmbedding := string(byteEmbedding)
	return strEmbedding, nil
}

func (mo *matrixOne) CloseConn() {
	db, _ := mo.dbConn.DB()
	db.Close()
}

func (mo *matrixOne) CheckIndex() bool {
	return mo.indexBuilt
}
