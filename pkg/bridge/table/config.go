package table

// Config is the configuration for a table
type Config struct {
	ID           string
	Headers      []string
	LimitOptions []LabeledValue
	HxPost       string
	HxSwapTarget string
}

// NewConfig provides html table configuration.
// It is recommended to store the id for this table in
// a constant so the htmlx selectors on outer target swap
// don't have race condition on key up debounce. Make any changes
// you want to this before creating a table based off it.
func NewConfig(id string, headers []string, hxPost string) Config {
	return Config{
		ID:           id,
		HxSwapTarget: "#" + id,
		Headers:      headers,
		HxPost:       hxPost,
		LimitOptions: DefaultPaginationLimitOptions,
	}
}
