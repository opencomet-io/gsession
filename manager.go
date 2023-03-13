package gsession

type Manager struct {
	TempStore    Store
	RefreshStore Store
	Codec        Codec
	TokenLength  int
}
