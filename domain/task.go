package domain

type Task interface {
	Run()
	AsJSON() TaskJSON
}

type TaskJSON interface{}
