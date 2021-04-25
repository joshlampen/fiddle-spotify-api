package action

import "context"

// Run calls Fetch, Execute, and Save
func Run(ctx context.Context, a Action) error {
	if err := a.Fetch(ctx); err != nil {
		return err
	}
	if err := a.Execute(ctx); err != nil {
		return err
	}
	if err := a.Save(ctx); err != nil {
		return err
	}

	return nil
}
