package interface_repo

type InterfaceEventManager interface {
	StorageEvent(name string) (string, error)
}
