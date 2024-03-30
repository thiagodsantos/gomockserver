package structs

type GraphQLConfig struct {
	Query    string            `json:"query"`
	Mutation string            `json:"mutation"`
	Fields   map[string]string `json:"fields"`
}