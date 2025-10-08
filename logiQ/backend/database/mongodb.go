package database

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 全局数据库实例 - 整个应用程序共享这个数据库连接
var DB *mongo.Database

// ConnectMongoDB 连接到MongoDB数据库
// 这个函数在应用启动时调用，建立与MongoDB的连接
func ConnectMongoDB() {
	// MongoDB连接字符串 - 默认连接到本地MongoDB实例
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI is not set")
	}

	dbName := os.Getenv("MONGO_DB_NAME")
	if dbName == "" {
		log.Fatal("MONGO_DB_NAME is not set")
	}
	// 设置MongoDB客户端选项
	clientOptions := options.Client().ApplyURI(mongoURI)

	// 创建上下文，设置10秒超时
	// context.WithTimeout创建一个带超时的上下文，防止连接操作无限等待
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() // 确保在函数结束时取消上下文，释放资源

	// 尝试连接到MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		// 如果连接失败，记录错误并退出程序
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	// 测试连接是否正常工作
	// Ping操作向数据库发送一个简单的命令来验证连接
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
	}

	// 选择数据库（如果数据库不存在，MongoDB会在第一次写入数据时自动创建）
	DB = client.Database(dbName)

	log.Println("Successfully connected to MongoDB!")
}

// GetCollection 获取指定名称的集合
// 这是一个辅助函数，让我们更容易获取MongoDB集合
// 集合类似于关系型数据库中的表
func GetCollection(collectionName string) *mongo.Collection {
	return DB.Collection(collectionName)
}

// 定义所有集合名称常量
// 使用常量可以避免在代码中硬编码集合名称，减少拼写错误
const (
	UsersCollection         = "users"                 // 用户集合 - 存储用户账户信息
	QuestionsCollection     = "questions"             // 题目集合 - 存储测验题目
	QuestionStatsCollection = "question_stats"        // 题目统计集合 - 存储题目使用和正确率统计
	QuizzesCollection       = "quizzes"               // 测验集合 - 存储测验记录和结果
	PendingUsersCollection  = "pending_registrations" // 待注册用户集合 - 存储未完成注册的用户信息
	UserStatsCollection     = "user_stats"            // 用户统计集合 - 存储用户统计数据
)
