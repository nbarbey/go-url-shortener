package urlshortener

type InfraStructure struct {
	store      Storer
	countStore CountStorer
}

func NewInMemoryInfrastructure() *InfraStructure {
	return &InfraStructure{
		store:      NewInMemorySqlite(),
		countStore: NewInMemoryCountStore(),
	}
}

func NewPGInfrastructure() *InfraStructure {
	return &InfraStructure{
		store:      NewPG(),
		countStore: NewPGCountStore(),
	}
}
