package services

import (
	"backend/database"
	"backend/models"
	"context"
	"errors"
	"math"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserStatsService 用户统计服务结构体 - 处理用户统计相关的业务逻辑
type UserStatsService struct {
	collection *mongo.Collection // MongoDB集合引用
}

// NewUserStatsService 创建新的用户统计服务实例
func NewUserStatsService() *UserStatsService {
	return &UserStatsService{
		collection: database.GetCollection(database.UserStatsCollection),
	}
}

// GetUserStatsByUserID 根据用户ID查询用户统计记录
func (s *UserStatsService) GetUserStatsByUserID(userID primitive.ObjectID) (*models.UserStats, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var userStats models.UserStats
	// 使用userID作为进行查询
	err := s.collection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&userStats)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("User statistics not found")
		}
		return nil, errors.New("Database query error")
	}

	return &userStats, nil
}

// CreateNewUserStats 创建新的用户统计记录(在用户注册时调用)
func (s *UserStatsService) CreateNewUserStats(userID primitive.ObjectID) (*models.UserStats, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 创建新的用户统计记录
	userStats := models.NewUserStats(userID)

	// 设置AccuracyRate的初值 - 初始化7天的数据，每天0%
	today := time.Now().Truncate(24 * time.Hour)
	userStats.AccuracyRate.Data = make([]models.AccuracyRateItem, 7)
	for i := 0; i < 7; i++ {
		date := today.AddDate(0, 0, -6+i) // 从7天前到今天
		userStats.AccuracyRate.Data[i] = models.AccuracyRateItem{
			Date:  date,
			Value: 0.0, // 初始准确率为0%
		}
	}

	// 插入到数据库
	_, err := s.collection.InsertOne(ctx, userStats)
	if err != nil {
		return nil, errors.New("Failed to create user statistics")
	}

	return userStats, nil

}

// UpdateUserStats 更新用户统计记录（在用户完成quiz后调用）
func (s *UserStatsService) UpdateUserStats(userID primitive.ObjectID, quiz *models.Quiz, quizService *QuizService) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 获取用户当前的统计记录
	oldStats, err := s.GetUserStatsByUserID(userID)
	if err != nil {
		return errors.New("Failed to get user stats")
	}
	// 创建新的统计记录
	newStats := *oldStats // 解引用创建副本

	// 1.Performance 部分
	s.updatePerformance(&newStats, quiz)

	// 2.AccuracyRate 部分
	s.updateAccuracyRate(&newStats, quizService)

	// 3.ErrorDistribution 部分
	s.updateErrorDistribution(&newStats, quiz)

	// 4.保存更新后的统计数据到数据库
	filter := bson.M{"user_id": userID}
	update := bson.M{"$set": bson.M{
		"performance":        newStats.Performance,
		"accuracy_rate":      newStats.AccuracyRate,
		"error_distribution": newStats.ErrorDistribution,
	}}
	_, err = s.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.New("Failed to update user statistics")
	}

	return nil

}

// updatePerformance 更新用户表现数据
func (s *UserStatsService) updatePerformance(userStats *models.UserStats, quiz *models.Quiz) {
	// 1)成功完成的task num
	userStats.Performance.TaskNum += quiz.CorrectQuestionsNum

	// 2)总得分
	score := 0
	// 遍历quiz中回答正确的问题，根据难度计算得分
	for _, question := range quiz.Questions {
		if question.IsCorrect {
			switch question.Question.Difficulty {
			case models.QuestionDifficultyEasy:
				score += 1
			case models.QuestionDifficultyMedium:
				score += 2
			case models.QuestionDifficultyHard:
				score += 3
			}
		}
	}
	userStats.Performance.Score += score

	// 3)Avg. Time
	avgTime := quiz.CompletionTime / len(quiz.Questions)
	userStats.Performance.AvgTime = float64(avgTime)
}

// updateAccuracyRate 更新准确率数据
func (s *UserStatsService) updateAccuracyRate(userStats *models.UserStats, quizService *QuizService) {
	today := time.Now().Truncate(24 * time.Hour) // 获取今天的日期（去除时分秒）

	// 获取当天所有已完成的quiz（包括当前刚添加的quiz）
	todayQuizzes, err := quizService.GetUserTodayQuizzes(userStats.UserID)
	if err != nil {
		// 如果获取失败，跳过更新
		return
	}

	// 计算当天累计准确率：当天总答对题数 / 当天总题数
	totalCorrectCount := 0
	totalQuestionCount := 0
	for _, dayQuiz := range todayQuizzes {
		totalCorrectCount += dayQuiz.CorrectQuestionsNum
		totalQuestionCount += len(dayQuiz.Questions)
	}

	// 防止除零错误
	if totalQuestionCount == 0 {
		return
	}

	todayAccuracy := float64(totalCorrectCount) / float64(totalQuestionCount)

	// 检查data数组中是否已存在今天的记录
	foundToday := false
	for i := range userStats.AccuracyRate.Data {
		// 比较日期（只比较年月日）
		if userStats.AccuracyRate.Data[i].Date.Truncate(24 * time.Hour).Equal(today) {
			// 存在：更新今天的累计准确率
			userStats.AccuracyRate.Data[i].Value = todayAccuracy
			foundToday = true
			break
		}
	}

	// 不存在：新增今天的记录
	if !foundToday {
		newRecord := models.AccuracyRateItem{
			Date:  today,
			Value: todayAccuracy,
		}
		userStats.AccuracyRate.Data = append(userStats.AccuracyRate.Data, newRecord)
	}

	// 按日期排序（升序）
	s.sortAccuracyDataByDate(userStats)

	// 维护最近7天数据：如果超过7条记录，删除最旧的
	if len(userStats.AccuracyRate.Data) > 7 {
		// 保留最新的7条记录（由于已按日期排序，删除前面的旧记录）
		userStats.AccuracyRate.Data = userStats.AccuracyRate.Data[len(userStats.AccuracyRate.Data)-7:]
	}
}

// sortAccuracyDataByDate 按日期对准确率数据进行排序（升序）
func (s *UserStatsService) sortAccuracyDataByDate(userStats *models.UserStats) {
	data := userStats.AccuracyRate.Data

	// 简单的冒泡排序（按日期升序）
	for i := 0; i < len(data)-1; i++ {
		for j := 0; j < len(data)-1-i; j++ {
			if data[j].Date.After(data[j+1].Date) {
				data[j], data[j+1] = data[j+1], data[j]
			}
		}
	}
}

// updateErrorDistribution 更新错误分布数据
func (s *UserStatsService) updateErrorDistribution(userStats *models.UserStats, quiz *models.Quiz) {
	// 1. 统计当次答错的题目按分类和难度累加到历史错题统计中
	for _, question := range quiz.Questions {
		if !question.IsCorrect {
			// 更新按分类的错误统计
			s.updateErrorByCategory(userStats, question.Question.Category)
			// 更新按难度的错误统计
			s.updateErrorByDifficulty(userStats, question.Question.Difficulty)
		}
	}

	// 2. 重新计算百分比
	s.recalculateErrorPercentages(userStats)
}

// updateErrorByCategory 更新按分类的错误统计
func (s *UserStatsService) updateErrorByCategory(userStats *models.UserStats, category models.QuestionCategory) {
	for i := range userStats.ErrorDistribution.DataByCategory {
		if userStats.ErrorDistribution.DataByCategory[i].Type == string(category) {
			userStats.ErrorDistribution.DataByCategory[i].Count++
			return
		}
	}
}

// updateErrorByDifficulty 更新按难度的错误统计
func (s *UserStatsService) updateErrorByDifficulty(userStats *models.UserStats, difficulty models.QuestionDifficulty) {
	for i := range userStats.ErrorDistribution.DataByDifficulty {
		if userStats.ErrorDistribution.DataByDifficulty[i].Type == string(difficulty) {
			userStats.ErrorDistribution.DataByDifficulty[i].Count++
			return
		}
	}
}

// recalculateErrorPercentages 重新计算错误分布百分比
func (s *UserStatsService) recalculateErrorPercentages(userStats *models.UserStats) {
	// 计算按分类的百分比
	s.calculateCategoryPercentages(userStats)
	// 计算按难度的百分比
	s.calculateDifficultyPercentages(userStats)
}

// calculateCategoryPercentages 计算按分类的错误百分比
func (s *UserStatsService) calculateCategoryPercentages(userStats *models.UserStats) {
	// 计算总错题数
	totalCategoryErrors := 0
	for _, item := range userStats.ErrorDistribution.DataByCategory {
		totalCategoryErrors += item.Count
	}

	// 防止除零错误
	if totalCategoryErrors == 0 {
		for i := range userStats.ErrorDistribution.DataByCategory {
			userStats.ErrorDistribution.DataByCategory[i].Value = 0
		}
		return
	}

	// 计算各分类的百分比
	for i := range userStats.ErrorDistribution.DataByCategory {
		percentage := math.Round(float64(userStats.ErrorDistribution.DataByCategory[i].Count)/float64(totalCategoryErrors)*100*100) / 100
		userStats.ErrorDistribution.DataByCategory[i].Value = percentage
	}
}

// calculateDifficultyPercentages 计算按难度的错误百分比
func (s *UserStatsService) calculateDifficultyPercentages(userStats *models.UserStats) {
	// 计算总错题数
	totalDifficultyErrors := 0
	for _, item := range userStats.ErrorDistribution.DataByDifficulty {
		totalDifficultyErrors += item.Count
	}

	// 防止除零错误
	if totalDifficultyErrors == 0 {
		for i := range userStats.ErrorDistribution.DataByDifficulty {
			userStats.ErrorDistribution.DataByDifficulty[i].Value = 0
		}
		return
	}

	// 计算各难度的百分比
	for i := range userStats.ErrorDistribution.DataByDifficulty {
		percentage := math.Round(float64(userStats.ErrorDistribution.DataByDifficulty[i].Count)/float64(totalDifficultyErrors)*100*100) / 100
		userStats.ErrorDistribution.DataByDifficulty[i].Value = percentage
	}
}
