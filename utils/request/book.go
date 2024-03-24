package request

type GetBooks struct {
	Offset int `query:"offset"`
	Limit  int `query:"limit"`
}
