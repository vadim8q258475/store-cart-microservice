package repo

type CartProduct struct {
	Id        uint32 `db:"id"`
	ProductId uint32 `db:"product_id"`
	Qty       uint32 `db:"qty"`
}

type Cart struct {
	Id     uint32 `db:"id"`
	UserId uint32 `db:"user_id"`
}
