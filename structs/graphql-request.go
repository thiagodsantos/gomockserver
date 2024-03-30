package structs

type GraphQLRequest struct {
	Query    string `json:"query"`
	Mutation string `json:"mutation"`
}
