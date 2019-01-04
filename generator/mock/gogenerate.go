package mock

//go:generate mockgen -destination mock.go -package mock github.com/deadcheat/goblet/generator UseCase,RegexpRepository,PathMatcherRepository
