package domain

type Task interface {
	Run() error
	AsJSON() TaskJSON
}

type TaskJSON interface{}
