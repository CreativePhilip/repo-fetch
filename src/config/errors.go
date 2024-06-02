package config

type CouldNotFindConfigError struct {
}

func (e CouldNotFindConfigError) Error() string {
	return "Couldn't find config, if this is a first time using the program run `rfetch init` first"
}
