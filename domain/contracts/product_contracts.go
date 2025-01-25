package contracts

type ProductRepository interface {
	CreateNewProduct()
	GetProduct()
	UpdateProductById()
	DeleteProductById()
}

type ProductService interface {
	CreateNewProduct()
	GetProduct()
	UpdateProductById()
	DeleteProductById()
}
