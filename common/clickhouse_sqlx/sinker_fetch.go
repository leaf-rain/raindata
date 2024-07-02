package clickhouse_sqlx

type Fetch interface {
	GetData() []string
	GetCallback() func()
}

type FetchSingle struct {
	Data     string
	Callback func()
}

func (f FetchSingle) GetData() []string {
	return []string{f.Data}
}

func (f FetchSingle) GetCallback() func() {
	return f.Callback
}
