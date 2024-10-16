package interfaces

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"go-todo-app/entities"
	"io"
	"net/http"
	"os"
	"time"
)

type JwtTodoRepository struct {
	Client *http.Client
	Jwt    string
}

func NewJwtTodoRepository() *JwtTodoRepository {
	return &JwtTodoRepository{
		Client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (repo *JwtTodoRepository) Insert(todo entities.Todo) error {
	// get token from cookie

	fmt.Println("Insert todo: ", todo)
	jsonData, err := json.Marshal(todo)
	if err != nil {
		return err
	}

	// 創建請求到 todo-app 的 /get-all API, 把 jwt 放在 header 裡
	req, err := http.NewRequest("POST", os.Getenv("TODO_APP_URL")+"/add", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	// 發送請求
	req.Header.Set("Authorization", "Bearer "+repo.Jwt)
	resp, err := repo.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fmt.Println("resp.StatusCode: ", resp.StatusCode)
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return errors.New("新增待辦事項失敗: " + resp.Status)
	}

	return nil
}

func (repo *JwtTodoRepository) GetAll() ([]entities.Todo, error) {
	// 創建請求到 todo-app 的 /get-all API
	req, err := http.NewRequest("GET", os.Getenv("TODO_APP_URL")+"/get-all", nil)
	if err != nil {
		return nil, err
	}

	// 發送請求
	fmt.Println("repo.Jwt: ", repo.Jwt)
	req.Header.Set("Authorization", "Bearer "+repo.Jwt)
	resp, err := repo.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 檢查回應狀態碼
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("取得待辦事項列表失敗: " + resp.Status)
	}

	// 打印回應內容，便於調試
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	fmt.Println("Response Body: ", string(bodyBytes))

	// 將回應內容轉換為待辦事項列表
	var todoResp entities.TodoResponse
	if err := json.Unmarshal(bodyBytes, &todoResp); err != nil {
		return nil, err
	}

	// 打印解析後的結果
	fmt.Println("todos: ", todoResp)

	return todoResp.Tasks, nil
}

func (repo *JwtTodoRepository) Delete(id string) error {
	req, err := http.NewRequest(http.MethodDelete, os.Getenv("TODO_APP_URL")+"/delete/"+id, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+repo.Jwt)
	resp, err := repo.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("刪除待辦事項失敗: " + resp.Status)
	}

	return nil
}

func (repo *JwtTodoRepository) SetJWT(jwt string) {
	fmt.Println("SetJwt: ", jwt)
	repo.Jwt = jwt
}
