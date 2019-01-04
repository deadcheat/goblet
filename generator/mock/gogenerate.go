package mock

// go:generate mockgen -destination generator/mock/mock.go -package mock github.com/deadcheat/goblet/generator UseCase,RegexpRepository,PathMatcherRepository
