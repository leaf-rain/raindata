package clickhouse_sqlx

type Fetch interface {
	GetData() []string
	GetCallback() []func()
	Copy() Fetch
}

var _ Fetch = (*FetchSingle)(nil)

type FetchSingle struct {
	Data     string
	Callback func()
}

func (f FetchSingle) Copy() Fetch {
	return FetchSingle{
		Data:     f.Data,
		Callback: f.Callback,
	}
}

func (f FetchSingle) GetData() []string {
	return []string{f.Data}
}

func (f FetchSingle) GetCallback() []func() {
	return []func(){f.Callback}
}

var _ Fetch = (*FetchArray)(nil)

type FetchArray struct {
	Data     []string
	Callback []func()
}

func (f FetchArray) Copy() Fetch {
	return FetchArray{
		Data:     f.Data,
		Callback: f.Callback,
	}
}

func (f FetchArray) GetData() []string {
	return f.Data
}

func (f FetchArray) GetCallback() []func() {
	return f.Callback
}
