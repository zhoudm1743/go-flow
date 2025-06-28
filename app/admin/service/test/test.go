package test

import "github.com/gin-gonic/gin"

type TestService interface {
	Test() gin.H
}

type testService struct {
	// db database.Database
}

func NewTestService() TestService {
	return &testService{
		// db: db,
	}
}

func (s *testService) Test() gin.H {
	return gin.H{
		"message": "test response from service",
		"status":  "success",
		"data":    "hello from TestService",
	}
}
