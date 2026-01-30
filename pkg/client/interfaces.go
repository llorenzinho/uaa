package client

type keyStoreGetter interface {
	Get(string) (*JwkKey, error)
}

type keyStoreSetter interface {
	Set(*JwkKey) error
}

type keyStoredeleter interface {
	Pop(string) error
}

type keyStoreCleaner interface {
	Clean()
}

type keyStoreChecker interface {
	Exist(string) bool
}

type keyStoreActive interface {
	Active() *JwkKey
}

type KeyStore interface {
	keyStoreGetter
	keyStoreSetter
	keyStoredeleter
	keyStoreCleaner
	keyStoreChecker
	keyStoreActive
}
