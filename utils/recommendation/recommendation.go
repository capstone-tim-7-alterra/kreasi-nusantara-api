package recommendation

import (
	"fmt"
	"kreasi-nusantara-api/entities"
)

func GetProductRecommendationInstruction() string {
	return "Kamu adalah seorang konselor yang bisa memberikan rekomendasi produk yang menarik dan belum pernah dimasukkan ke cart oleh user. \n\nCatatan:\n- data 'sudah pernah masuk cart' yang diberikan user berupa list nama product menggunakan format csv (uid, nama product)\n- respon berupa 4 UUID SAJA yang berdasarkan pada data 'belum pernah masuk cart' yang diberikan user, TANPA kata pengantar maupun penutup\n- jika data 'belum pernah masuk cart' kosong, berikan rekomendasi produk yang yang menurutmu relevan\n\nContoh respon:\n9764e579-8200-4171-8a1c-76e879590757\n9764e579-8200-4171-8a1c-76e879590757"
}

func ToRecommendationPrompt(enteredCart *[]entities.CartItems, unEnteredCart *[]entities.Products) string {
	var (
		enteredPrompt = "sudah pernah masuk cart:\n"
		unEnteredPrompt = "belum pernah masuk cart:\n"
	)

	enteredPrompt += toCartPrompt(enteredCart)
	unEnteredPrompt += toProductPrompt(unEnteredCart)

	prompt := enteredPrompt + unEnteredPrompt
	return prompt
}

func toProductPrompt(products *[]entities.Products) string {
	var prompt string
	for _, product := range *products {
		prompt += fmt.Sprintf("- %s,%s\n", product.ID, product.Name)
	}
	return prompt
}

func toCartPrompt(cartItems *[]entities.CartItems) string {
	var prompt string
	for _, cartItem := range *cartItems {
		prompt += fmt.Sprintf("- %s,%s\n", cartItem.ProductVariant.ProductID, cartItem.ProductVariant.Products.Name)
	}
	return prompt
}
