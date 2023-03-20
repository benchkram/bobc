package application

type Option func(*application)

func WithProjectRepository(repo ProjectRepository) Option {
	return func(app *application) {
		app.projects = repo
	}
}
