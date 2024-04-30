package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Post struct {
	ID          int    `bson:"id" json:"id"`
	IsFavorited bool   `bson:"isFavorited" json:"isFavorited"`
	Title       string `bson:"title" json:"title"`
	Content     string `bson:"content" json:"content"`
}

var client *mongo.Client
var postCollection *mongo.Collection

func connectDB() *mongo.Client {
	uri := "mongodb://max:123456@localhost:27017" // Update with your connection string
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}

	return client
}

func main() {
	// Initialize the Gin router
	router := gin.Default()

	// Connect to MongoDB
	client := connectDB()
	postCollection = client.Database("blog").Collection("posts")

	// Define routes
	router.GET("/posts", getPosts)
	router.POST("/posts", createPost) // 添加这行

	// Start serving the application
	router.Run(":8080")
}

func getPosts(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var posts []Post
	results, err := postCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Printf("Error finding posts: %v", err) // 日志记录错误详情
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching posts"})
		return
	}
	defer results.Close(ctx)

	for results.Next(ctx) {
		var post Post
		if err := results.Decode(&post); err != nil {
			log.Printf("Error decoding post: %v", err) // 日志记录错误详情
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding post"})
			return
		}
		posts = append(posts, post)
	}

	c.JSON(http.StatusOK, posts)
}

func createPost(c *gin.Context) {
	var newPost Post
	// 尝试从请求体中解析 JSON
	if err := c.ShouldBindJSON(&newPost); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 创建上下文，设置超时
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 插入新文章到数据库
	result, err := postCollection.InsertOne(ctx, newPost)
	if err != nil {
		log.Printf("Error decoding post: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}

	// 返回成功消息和创建的文章的 ID
	c.JSON(http.StatusOK, gin.H{"message": "Post created successfully", "post_id": result.InsertedID})
}

func deletePost(c *gin.Context) {

}
