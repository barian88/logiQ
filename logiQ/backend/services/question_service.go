package services

import (
	"backend/database"
	gengerationService "backend/generation/service"
	"backend/models"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type QuestionService struct {
	collection *mongo.Collection
}

func NewQuestionService() *QuestionService {
	return &QuestionService{
		collection: database.GetCollection(database.QuestionsCollection),
	}
}

// GenerateQuestion 从外部API生成题目并存入数据库，返回插入的题目数量
func (s *QuestionService) GenerateQuestion(req *models.GenerateQuestionRequest) ([]models.Question, error) {
	if req.Number == 0 {
		req.Number = 1
	}
	genService := gengerationService.NewService()
	questionList, err := genService.GenerateQuestion(req.Number, req.Category, req.Difficulty, req.Type)
	if err != nil {
		return nil, err
	}
	// 批量插入题目
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	docs := make([]interface{}, len(questionList))
	for i, q := range questionList {
		docs[i] = q
	}
	_, err = s.collection.InsertMany(ctx, docs)
	if err != nil {
		return nil, err
	}
	return questionList, nil
}

// GetQuestionByID 根据ID获取题目
func (s *QuestionService) GetQuestionByID(questionID primitive.ObjectID) (*models.Question, error) {
	// TODO: 根据ID获取题目
	return nil, nil
}

// GetQuestionList 获取题目列表，支持分页和筛选
func (s *QuestionService) GetQuestionList(category models.QuestionCategory, difficulty models.QuestionDifficulty, qType models.QuestionType, page, pageSize int) ([]models.QuestionResponseForAdmin, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 5
	}
	skip := (page - 1) * pageSize

	filter := bson.M{"is_active": true}
	if category != "" {
		filter["category"] = string(category)
	}
	if difficulty != "" {
		filter["difficulty"] = string(difficulty)
	}
	if qType != "" {
		filter["type"] = string(qType)
	}
	// 查询符合条件的总数
	count, err := s.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: filter}},
	}

	if skip > 0 {
		pipeline = append(pipeline, bson.D{{Key: "$skip", Value: int64(skip)}})
	}
	if pageSize > 0 {
		pipeline = append(pipeline, bson.D{{Key: "$limit", Value: int64(pageSize)}})
	}

	pipeline = append(pipeline,
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "question_stats"},
			{Key: "localField", Value: "_id"},
			{Key: "foreignField", Value: "question_id"},
			{Key: "as", Value: "stats"},
		}}},
		bson.D{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$stats"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}}},
		bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "total_answers", Value: bson.D{{Key: "$ifNull", Value: bson.A{"$stats.total_answers", 0}}}},
			{Key: "correct_answers", Value: bson.D{{Key: "$ifNull", Value: bson.A{"$stats.correct_answers", 0}}}},
		}}},
		bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "accuracy_rate", Value: bson.D{{Key: "$cond", Value: bson.D{
				{Key: "if", Value: bson.D{{Key: "$gt", Value: bson.A{"$total_answers", 0}}}},
				{Key: "then", Value: bson.D{{Key: "$round", Value: bson.A{
					bson.D{{Key: "$divide", Value: bson.A{"$correct_answers", "$total_answers"}}},
					3,
				}}}},
				{Key: "else", Value: 0},
			}}}},
		}}},
		bson.D{{Key: "$project", Value: bson.D{{Key: "stats", Value: 0}}}},
	)

	cursor, err := s.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var questions []models.QuestionResponseForAdmin
	if err = cursor.All(ctx, &questions); err != nil {
		return nil, 0, err
	}

	return questions, int(count), nil
}

// GetRandomQuestions 根据category difficulty获取10个题目，用来创建quiz
func (s *QuestionService) GetRandomQuestions(category models.QuestionCategory, difficulty models.QuestionDifficulty, count int) ([]models.Question, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 设置count默认值
	if count <= 0 {
		count = 10
	}

	// 构建筛选条件
	filter := bson.M{"is_active": true}

	if category != "" {
		filter["category"] = string(category)
	}

	if difficulty != "" {
		filter["difficulty"] = string(difficulty)
	}

	// 使用MongoDB的$sample进行随机抽样
	pipeline := []bson.M{
		{"$match": filter},
		{"$sample": bson.M{"size": count}},
	}

	cursor, err := s.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var questions []models.Question
	if err = cursor.All(ctx, &questions); err != nil {
		return nil, err
	}

	return questions, nil
}

// SoftDeleteQuestionsByIDs 软删除题目
func (s *QuestionService) SoftDeleteQuestionsByIDs(questionIDs []primitive.ObjectID) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": bson.M{"$in": questionIDs}}
	update := bson.M{"$set": bson.M{"is_active": false}}

	result, err := s.collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return 0, err
	}

	return result.ModifiedCount, nil
}
