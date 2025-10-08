package services

import (
	"backend/database"
	"backend/models"
	"context"
	"math"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// QuestionStatsService handles database operations for question statistics.
type QuestionStatsService struct {
	questionCollection *mongo.Collection
	statsCollection    *mongo.Collection
}

// NewQuestionStatsService creates a new QuestionStatsService.
func NewQuestionStatsService() *QuestionStatsService {
	return &QuestionStatsService{
		questionCollection: database.GetCollection(database.QuestionsCollection),
		statsCollection:    database.GetCollection(database.QuestionStatsCollection),
	}
}

// UpdateStats updates the total and correct answer counts for a given question.
// It uses an upsert operation to create the document if it doesn't exist.
func (s *QuestionStatsService) UpdateStats(questionID primitive.ObjectID, isCorrect bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"question_id": questionID}

	update := bson.M{
		"$inc": bson.M{"total_answers": 1},
	}

	if isCorrect {
		update["$inc"].(bson.M)["correct_answers"] = 1
	}

	opts := options.Update().SetUpsert(true)

	_, err := s.statsCollection.UpdateOne(ctx, filter, update, opts)
	return err
}

func (s *QuestionStatsService) CalcQuestionDistributions() (*models.QuestionDistributions, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 使用 $facet 聚合管道，同时统计难度、类型、分类和总数
	pipeline := mongo.Pipeline{
		{
			{"$facet", bson.M{
				"difficulty": []bson.M{
					// 按难度分组并统计数量
					{"$group": bson.M{"_id": "$difficulty", "count": bson.M{"$sum": 1}}},
				},
				"type": []bson.M{
					// 按类型分组并统计数量
					{"$group": bson.M{"_id": "$type", "count": bson.M{"$sum": 1}}},
				},
				"category": []bson.M{
					// 按分类分组并统计数量
					{"$group": bson.M{"_id": "$category", "count": bson.M{"$sum": 1}}},
				},
				"total": []bson.M{
					// 统计总题目数量
					{"$count": "count"},
				},
			}},
		},
	}

	// 执行聚合查询
	cursor, err := s.questionCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// 用于接收 $facet 聚合结果的临时结构体
	var results []struct {
		Difficulty []struct {
			ID    string `bson:"_id"`
			Count int64  `bson:"count"`
		} `bson:"difficulty"`
		Type []struct {
			ID    string `bson:"_id"`
			Count int64  `bson:"count"`
		} `bson:"type"`
		Category []struct {
			ID    string `bson:"_id"`
			Count int64  `bson:"count"`
		} `bson:"category"`
		Total []struct {
			Count int64 `bson:"count"`
		} `bson:"total"`
	}

	// 解析聚合结果到临时结构体
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	// 检查是否有结果可处理
	if len(results) == 0 || len(results[0].Total) == 0 || results[0].Total[0].Count == 0 {
		// 没有题目时返回空结构体（非 nil）
		return &models.QuestionDistributions{}, nil
	}

	// 总题目数量
	total := float64(results[0].Total[0].Count)
	// 初始化分布结构体
	distributions := &models.QuestionDistributions{
		Difficulty: make([]models.DimensionPortion, len(results[0].Difficulty)),
		Type:       make([]models.DimensionPortion, len(results[0].Type)),
		Category:   make([]models.DimensionPortion, len(results[0].Category)),
	}

	// 辅助函数：将原始数据转换为 DimensionPortion
	transform := func(item struct {
		ID    string `bson:"_id"`
		Count int64  `bson:"count"`
	}) models.DimensionPortion {
		return models.DimensionPortion{
			Value:   item.ID,
			Count:   item.Count,
			Portion: math.Round(float64(item.Count)/total*100) / 100, // 保留两位小数
		}
	}

	// 填充难度分布
	for i, item := range results[0].Difficulty {
		distributions.Difficulty[i] = transform(item)
	}
	// 填充类型分布
	for i, item := range results[0].Type {
		distributions.Type[i] = transform(item)
	}
	// 填充分类分布
	for i, item := range results[0].Category {
		distributions.Category[i] = transform(item)
	}

	return distributions, nil
}

// CalcAccuracyByDimension 计算各维度（难度、类型、分类）的答题准确率
func (s *QuestionStatsService) CalcAccuracyByDimension() (*models.AccuracyDistributions, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pipeline := mongo.Pipeline{
		// Join question_stats with questions to get dimensions
		{{Key: "$lookup", Value: bson.M{
			"from":         "questions",
			"localField":   "question_id",
			"foreignField": "_id",
			"as":           "questionInfo",
		}}},
		// Deconstruct the questionInfo array
		{{Key: "$unwind", Value: "$questionInfo"}},
		// Facet for parallel aggregation
		{{Key: "$facet", Value: bson.M{
			"difficulty": s.buildAccuracyPipeline("$questionInfo.difficulty"),
			"type":       s.buildAccuracyPipeline("$questionInfo.type"),
			"category":   s.buildAccuracyPipeline("$questionInfo.category"),
		}}},
	}

	cursor, err := s.statsCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Temporary struct to decode the facet result
	var results []models.AccuracyDistributions
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return &models.AccuracyDistributions{}, nil
	}

	return &results[0], nil
}

// buildAccuracyPipeline 构建计算准确率的Pipeline
func (s *QuestionStatsService) buildAccuracyPipeline(dimensionField string) mongo.Pipeline {
	return mongo.Pipeline{
		{{Key: "$group", Value: bson.M{
			"_id":             dimensionField,
			"total_answers":   bson.M{"$sum": "$total_answers"},
			"correct_answers": bson.M{"$sum": "$correct_answers"},
		}}},
		{{Key: "$match", Value: bson.M{"_id": bson.M{"$ne": nil}}}},
		{{Key: "$project", Value: bson.M{
			"_id":             0,
			"value":           "$_id",
			"total_answers":   "$total_answers",
			"correct_answers": "$correct_answers",
			"accuracy": bson.M{
				"$round": bson.A{
					bson.M{
						"$cond": bson.A{
							bson.M{"$gt": bson.A{"$total_answers", 0}},
							bson.M{"$divide": bson.A{"$correct_answers", "$total_answers"}},
							0,
						},
					},
					2,
				},
			},
		}}},
		{{Key: "$sort", Value: bson.M{"value": 1}}},
	}
}
