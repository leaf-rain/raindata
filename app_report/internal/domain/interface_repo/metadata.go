package interface_repo

type InterfaceMetadataRepo interface {
	MetadataPut(keys []string) (map[string]string, error)
}

type DefaultMetadata struct{}

func NewMetadata() *DefaultMetadata {
	return &DefaultMetadata{}
}

func (d DefaultMetadata) MetadataPut(keys []string) (map[string]string, error) {
	var result = make(map[string]string)
	for i := range keys {
		result[keys[i]] = keys[i]
	}
	return result, nil
}
