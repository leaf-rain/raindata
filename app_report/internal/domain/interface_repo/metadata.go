package interface_repo

type InterfaceMetadataRepo interface {
	MetadataPut(keys []string) (map[string]string, error)
}
