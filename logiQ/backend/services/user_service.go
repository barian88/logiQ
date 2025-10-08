package services

import (
	"backend/database"
	"backend/models"
	"backend/utils"
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/sas"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserService 用户服务结构体 - 处理所有用户相关的业务逻辑
type UserService struct {
	collection          *mongo.Collection    // MongoDB集合引用
	pendingCollection   *mongo.Collection    // 待注册用户集合引用
	verificationService *VerificationService // 验证码服务
	userStatsService    *UserStatsService    // 用户统计服务
}

// NewUserService 创建新的用户服务实例
// 这是一个构造函数，返回初始化好的UserService
func NewUserService(verificationService *VerificationService, userStatsService *UserStatsService) *UserService {
	return &UserService{
		// 获取users集合的引用，用于数据库操作
		collection:          database.GetCollection(database.UsersCollection),
		pendingCollection:   database.GetCollection(database.PendingUsersCollection),
		verificationService: verificationService,
		userStatsService:    userStatsService,
	}
}

// GetUserByEmail 根据邮箱查找用户 - 主要用于登录验证
func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User
	// 在数据库中查找指定邮箱的用户
	err := s.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("User not found")
		}
		return nil, errors.New("Database query error")
	}

	return &user, nil
}

// Login 用户登录验证
func (s *UserService) Login(req *models.LoginRequest) (*models.User, error) {
	// 第1步：根据邮箱查找用户
	user, err := s.GetUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}

	// 第2步：验证密码
	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, errors.New("Incorrect password")
	}

	// 第3步：返回用户信息（不包含密码）
	return user, nil
}

// GetUserByID 根据用户ID查找用户 - 用于获取用户详细信息
func (s *UserService) GetUserByID(userID primitive.ObjectID) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User
	// 根据用户ID查找用户
	err := s.collection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("User not found")
		}
		return nil, errors.New("Database query error")
	}

	return &user, nil
}

// RegisterRequest 第一步：处理注册请求，保存待注册信息
func (s *UserService) RegisterRequest(req *models.RegisterRequestModel) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 1. 检查邮箱是否已经注册
	existingUser, err := s.GetUserByEmail(req.Email)
	if err == nil && existingUser != nil {
		return errors.New("Email already registered")
	}

	// 2. 加密密码
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return errors.New("Password encryption failed")
	}

	// 3. 删除该邮箱之前的待注册记录
	_, err = s.pendingCollection.DeleteMany(ctx, bson.M{"email": req.Email})
	if err != nil {
		return errors.New("Failed to clean pending registration data")
	}

	// 4. 创建待注册记录
	pendingReg := models.PendingRegistration{
		Username:  req.Username,
		Email:     req.Email,
		Password:  hashedPassword,
		ExpiresAt: time.Now().Add(10 * time.Minute), // 10分钟后过期
		CreatedAt: time.Now(),
	}

	_, err = s.pendingCollection.InsertOne(ctx, pendingReg)
	if err != nil {
		return errors.New("Failed to save pending registration data")
	}

	return nil
}

// CompleteRegistration 完成注册（使用临时token）
func (s *UserService) CompleteRegistration(email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 查找待注册记录
	var pendingReg models.PendingRegistration
	err := s.pendingCollection.FindOne(ctx, bson.M{
		"email":      email,
		"expires_at": bson.M{"$gt": time.Now()},
	}).Decode(&pendingReg)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("Registration request has expired, please re-register")
		}
		return nil, errors.New("Failed to find registration data")
	}

	// 创建正式用户
	user, err := s.createUserDirectly(&pendingReg)
	if err != nil {
		return nil, err
	}

	// 删除待注册记录
	_, err = s.pendingCollection.DeleteOne(ctx, bson.M{"_id": pendingReg.ID})
	if err != nil {
		// 注册已经成功，这个错误不影响结果，只记录日志
		// log.Printf("删除待注册记录失败: %v", err)
	}

	return user, nil
}

// createUserDirectly 直接创建用户（密码已加密）
func (s *UserService) createUserDirectly(pendingReg *models.PendingRegistration) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	now := time.Now()
	newUser := models.User{
		Username:          pendingReg.Username,
		Email:             pendingReg.Email,
		Password:          pendingReg.Password,                                            // 已加密的密码
		ProfilePictureUrl: "https://logiq.blob.core.windows.net/bucket-logiq/default.JPG", // 默认头像
		Role:              "user",
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	// 用户插入到user表
	result, err := s.collection.InsertOne(ctx, newUser)
	if err != nil {
		return nil, errors.New("User creation failed")
	}
	newUser.ID = result.InsertedID.(primitive.ObjectID)
	// 创建用户统计记录
	_, err = s.userStatsService.CreateNewUserStats(newUser.ID)
	if err != nil {
		return nil, err
	}

	return &newUser, nil
}

// UpdatePassword 更新用户密码
func (s *UserService) UpdatePassword(email, newPassword string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 1. 检查用户是否存在
	var user models.User
	err := s.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("User not found")
		}
		return errors.New("Database query error")
	}

	// 2. 加密新密码
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return errors.New("Password encryption failed")
	}

	// 3. 更新密码
	_, err = s.collection.UpdateOne(ctx,
		bson.M{"email": email},
		bson.M{"$set": bson.M{
			"password":   hashedPassword,
			"updated_at": time.Now(),
		}},
	)
	if err != nil {
		return errors.New("Password update failed")
	}

	return nil
}

// UpdateProfile 更新用户资料, 允许更新用户名和头像URL
func (s *UserService) UpdateProfile(userID primitive.ObjectID, req *models.UpdateProfileRequest) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 构建更新字段
	updateFields := bson.M{}
	if req.Username != "" {
		updateFields["username"] = req.Username
	}
	if req.ProfilePictureUrl != "" {
		updateFields["profile_picture_url"] = req.ProfilePictureUrl
	}
	updateFields["updated_at"] = time.Now()

	// 执行更新操作
	_, err := s.collection.UpdateOne(ctx, bson.M{"_id": userID}, bson.M{"$set": updateFields})
	if err != nil {
		return nil, errors.New("failed to update user profile")
	}
	// 返回更新后的用户信息
	updatedUser, err := s.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	return updatedUser, nil

}

// LoginAdmin 管理员登录验证
func (s *UserService) LoginAdmin(req *models.LoginRequest) (*models.User, error) {
	// 第1步：根据邮箱查找用户
	user, err := s.GetUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}

	// 第2步：验证密码
	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, errors.New("Incorrect password")
	}

	// 第3步：检查用户角色是否为admin
	if user.Role != "admin" {
		return nil, errors.New("Access denied: not an admin user")
	}

	// 第4步：返回用户信息（不包含密码）
	return user, nil

}

// PresignUploadURL 生成Azure Blob存储的预签名上传URL
func (s *UserService) PresignUploadURL(blobName string) (string, error) {
	container := os.Getenv("AZURE_BLOB_CONTAINER")
	if container == "" {
		return "", errors.New("AZURE_BLOB_CONTAINER is not set")
	}
	acct := os.Getenv("AZURE_BLOB_ACCOUNT")
	if acct == "" {
		return "", errors.New("AZURE_BLOB_ACCOUNT is not set")
	}
	key := os.Getenv("AZURE_BLOB_KEY")
	if key == "" {
		return "", errors.New("AZURE_BLOB_KEY is not set")
	}

	cred, err := azblob.NewSharedKeyCredential(acct, key)
	if err != nil {
		return "", err
	}

	// 生成写入权限的 SAS（短效 5 分钟）
	now := time.Now().UTC()
	start := now.Add(-2 * time.Minute)
	exp := now.Add(5 * time.Minute)
	perms := sas.BlobPermissions{Write: true, Create: true}
	sigValues := sas.BlobSignatureValues{
		StartTime:     start,
		ExpiryTime:    exp,
		Permissions:   perms.String(),
		ContainerName: container,
		BlobName:      blobName,
	}
	sasQueryParams, err := sigValues.SignWithSharedKey(cred)
	if err != nil {
		return "", err
	}

	sasURL := fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s?%s",
		acct, container, blobName, sasQueryParams.Encode())
	return sasURL, nil
}
