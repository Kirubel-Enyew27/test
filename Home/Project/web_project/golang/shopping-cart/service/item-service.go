package service


type Cart struct {
    Items    map[string]*Item
    Discount float64
}

func NewCart() *Cart {
    return &Cart{Items: make(map[string]*Item)}
}