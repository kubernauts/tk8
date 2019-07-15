package server

type AddonRequest struct {
	// Name of the addon
	Name string
	// Namespace of the addon
	Scope string
}

type AddonResponse struct {
	// Status
	Status string
	// Error
	Error string
}
