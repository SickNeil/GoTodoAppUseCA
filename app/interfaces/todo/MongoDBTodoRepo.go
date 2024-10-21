// interfaces/mongo_todo_repository.go
package interfaces

import (
	"context"
	"fmt"
	"go-todo-app/entities"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDBTodoRepo struct {
	Collection *mongo.Collection
}

func NewMongoDBTodoRepository(collection *mongo.Collection) *MongoDBTodoRepo {
	return &MongoDBTodoRepo{Collection: collection}
}

// Insert 用於新增待辦事項
func (m *MongoDBTodoRepo) Insert(todo entities.Todo) error {
	var objID primitive.ObjectID
	var err error

	fmt.Println("Insert todo: ", todo)
	fmt.Println("Insert todo: ", todo.ID, todo.ID == "")
	// 如果 ID 是空的，則讓 MongoDB 自動生成 ObjectID
	if todo.ID == "" {
		objID = primitive.NewObjectID()
		todo.ID = objID.Hex() // 將生成的 ObjectID 設置為字符串形式的 ID
	} else {
		// 將字符串形式的 ID 轉換為 ObjectID
		objID, err = primitive.ObjectIDFromHex(todo.ID)
		if err != nil {
			return err
		}
	}

	// 插入到 MongoDB
	_, err = m.Collection.InsertOne(context.TODO(), bson.M{
		"_id":        objID,     // MongoDB 的 _id 欄位
		"task":       todo.Task, // 其他字段
		"created_at": todo.CreatedAt,
	})

	return err
}

func (m *MongoDBTodoRepo) GetAll() ([]entities.Todo, error) {
	var todos []entities.Todo
	cursor, err := m.Collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var todoDB struct {
			ID        primitive.ObjectID `bson:"_id"`
			Task      string             `bson:"task"`
			CreatedAt time.Time          `bson:"created_at"`
		}
		err := cursor.Decode(&todoDB)
		if err != nil {
			return nil, err
		}

		// 將 ObjectID 轉換為字符串，並創建 Todo 實體
		todo := entities.Todo{
			ID:        todoDB.ID.Hex(), // 將 ObjectID 轉換為字符串
			Task:      todoDB.Task,
			CreatedAt: todoDB.CreatedAt,
		}

		todos = append(todos, todo)
	}

	return todos, nil
}

func (m *MongoDBTodoRepo) Delete(id string) error {
	// 將字符串 ID 轉換為 ObjectID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	// 刪除該 ObjectID 對應的記錄
	_, err = m.Collection.DeleteOne(context.TODO(), bson.M{"_id": objID})
	return err
}
