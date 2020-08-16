package app

// MockServerDB is used in tests which require a mocked db.
// MockServerDB implements the Database interface.
type MockServerDB struct {
	closeMock        func() error
	checkHealthMock  func() error
	createTablesMock func()

	// Articles
	getArticlesMock       func(q, pubDate string) ([]*Article, bool)
	getRecentArticlesMock func() []*Article
	insertArticleMock     func(article *Article) (int, error)
	insertArticlesMock    func(articles []*Article) []int

	// TaskLog
	getRecentTaskLogMock func(string) *TaskLog
	insertTaskLogMock    func(tasklog *TaskLog) (int, error)
	recordTaskMock       func(task string, manual bool)
}

// Close is exported
func (mc *MockServerDB) Close() error {
	return mc.closeMock()
}

// CreateTables is exported
func (mc *MockServerDB) CreateTables() {
	// We don't need to return anything since there's no return value.
}

// GetArticles is exported
func (mc *MockServerDB) GetArticles(q, pubDate string) ([]*Article, bool) {
	return mc.getArticlesMock(q, pubDate)
}

// GetRecentArticles is exported
func (mc *MockServerDB) GetRecentArticles() []*Article {
	return mc.getRecentArticlesMock()
}

// InsertArticle is exported
func (mc *MockServerDB) InsertArticle(article *Article) (int, error) {
	return mc.insertArticleMock(article)
}

// InsertArticles is exported
func (mc *MockServerDB) InsertArticles(articles []*Article) []int {
	return mc.insertArticlesMock(articles)
}

// GetRecentTaskLog is exported
func (mc *MockServerDB) GetRecentTaskLog(task string) *TaskLog {
	return mc.getRecentTaskLogMock(task)
}

// InsertTaskLog is exported
func (mc *MockServerDB) InsertTaskLog(tasklog *TaskLog) (int, error) {
	return mc.insertTaskLogMock(tasklog)
}

// RecordTask is exported
func (mc *MockServerDB) RecordTask(task string, manual bool) {
	// We don't need to return anything since there's no return value.
}
