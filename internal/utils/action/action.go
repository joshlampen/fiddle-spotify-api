package action

import "context"

// Fetcher is for fetching the inputs required to perform the action
type Fetcher interface {
	Fetch(ctx context.Context) error
}

// Executor is for executing the action and storing the output
type Executor interface {
	Execute(ctx context.Context) error
}

// Saver is for saving the output to the database
type Saver interface {
	Save(ctx context.Context) error
}


// Action is composed of three interfaces that help enforce order of operations
// and separation of concerns when performing an action
type Action interface {
	Fetcher
	Executor
	Saver
}
