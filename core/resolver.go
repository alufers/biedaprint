package core

import "context"

type Resolver struct {
	App *App
}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}
func (r *Resolver) Subscription() SubscriptionResolver {
	return &subscriptionResolver{r}
}

type mutationResolver struct{ *Resolver }

type queryResolver struct{ *Resolver }

type subscriptionResolver struct{ *Resolver }

func (r *queryResolver) Configuration(ctx context.Context) (*RootConfiguration, error) {
	panic("not implemented")
}
func (r *queryResolver) ConfigurationDescriptor(ctx context.Context) (ConfigurationFieldDescriptor, error) {
	panic("not implemented")
}
