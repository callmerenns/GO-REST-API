package config

const (
	// Routing Group
	ApiGroup = "/api/v1"

	// Routing Products
	GetProductsList     = "/products"
	GetProducts         = "/products/:id"
	GetProductsByStocks = "/products/stock/:stock"
	PostProducts        = "/products"
	PutProducts         = "/products/:id"
	DelProducts         = "/products/:id"

	// Routing Users
	GetUsersList = "/profiles"
	GetUsers     = "/profiles/:id"

	// Routing Auth
	PostRegister = "/auth/register"
	PostLogin    = "/auth/login"
	GetLogout    = "/auth/logout"
)
