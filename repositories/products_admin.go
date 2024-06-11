package repositories

import (
	"context"
	"kreasi-nusantara-api/entities"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductAdminRepository interface {
	// Product
	CreateProduct(ctx context.Context, product *entities.Products) error
	GetAllProduct(ctx context.Context, page, limit int) ([]*entities.Products, error)
	GetProduct(ctx context.Context, product *entities.Products) (*entities.Products, error)
	GetProductByID(ctx context.Context, productID uuid.UUID) (*entities.Products, error) 
	SearchProductByName(ctx context.Context, name string, page, limit int) ([]*entities.Products, error)
	DeleteProduct(ctx context.Context, productID uuid.UUID) error 
	UpdateProduct(ctx context.Context, productID uuid.UUID, product *entities.Products) error 
	// Category
	CreateCategory(ctx context.Context, category *entities.ProductCategory) error
	DeleteCategory(ctx context.Context, id int) error
	UpdateCategory(ctx context.Context, category *entities.ProductCategory) error
	GetCategory(ctx context.Context, category *entities.ProductCategory) (*entities.ProductCategory, error)
	GetAllCategory(ctx context.Context) ([]*entities.ProductCategory, error)
	GetCategoryByID(ctx context.Context, id int) (*entities.ProductCategory, error)
}

type productAdminRepository struct {
	DB *gorm.DB
}

func NewProductAdminRepository(db *gorm.DB) *productAdminRepository {
	return &productAdminRepository{
		DB: db,
	}
}

// Product

func (pr *productAdminRepository) CreateProduct(ctx context.Context, product *entities.Products) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return pr.DB.Create(product).Error
}

func (pr *productAdminRepository) GetAllProduct(ctx context.Context, page, limit int) ([]*entities.Products, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	var products []*entities.Products
	if page > 0 && limit > 0 {
		offset := (page - 1) * limit
		if err := pr.DB.Preload("ProductPricing").Preload("ProductVariants").Preload("ProductImages").Preload("ProductVideos").Limit(limit).Offset(offset).Find(&products).Error; err != nil {
			return nil, err
		}
	} else {
		if err := pr.DB.Preload("ProductPricing").Preload("ProductVariants").Preload("ProductImages").Preload("ProductVideos").Find(&products).Error; err != nil {
			return nil, err
		}
	}
	return products, nil
}

func (pr *productAdminRepository) GetProduct(ctx context.Context, product *entities.Products) (*entities.Products, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if err := pr.DB.Where(product).First(product).Error; err != nil {
		return nil, err
	}
	return product, nil
}

func (pr *productAdminRepository) SearchProductByName(ctx context.Context, name string, page, limit int) ([]*entities.Products, error) {
    if err := ctx.Err(); err != nil {
        return nil, err
    }
    
    var products []*entities.Products
    query := pr.DB.Preload("ProductPricing").Preload("ProductVariants").Preload("ProductImages").Preload("ProductVideos").Where("name ILIKE ?", "%"+name+"%")
    
    if page > 0 && limit > 0 {
        offset := (page - 1) * limit
        query = query.Limit(limit).Offset(offset)
    }
    
    if err := query.Find(&products).Error; err != nil {
        return nil, err
    }
    
    return products, nil
}


// Categories

func (pr *productAdminRepository) CreateCategory(ctx context.Context, category *entities.ProductCategory) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return pr.DB.Create(category).Error
}

func (pr *productAdminRepository) GetCategory(ctx context.Context, category *entities.ProductCategory) (*entities.ProductCategory, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if err := pr.DB.Where(category).First(category).Error; err != nil {
		return nil, err
	}
	return category, nil
}

func (pr *productAdminRepository) GetAllCategory(ctx context.Context) ([]*entities.ProductCategory, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	var categories []*entities.ProductCategory
	if err := pr.DB.Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (pr *productAdminRepository) GetCategoryByID(ctx context.Context, id int) (*entities.ProductCategory, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	var category entities.ProductCategory
	if err := pr.DB.Where("id = ?", id).First(&category).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (pr *productAdminRepository) UpdateCategory(ctx context.Context, category *entities.ProductCategory) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return pr.DB.Save(category).Error
}

func (pr *productAdminRepository) DeleteCategory(ctx context.Context, id int) error {

	if err := pr.DB.WithContext(ctx).Where("id = ?", id).Delete(&entities.ProductCategory{}).Error; err != nil {
		return err
	}
	return nil
}

func (pr *productAdminRepository) UpdateProduct(ctx context.Context, productID uuid.UUID, product *entities.Products) error {
    if err := ctx.Err(); err != nil {
        return err
    }
    return pr.DB.Model(&entities.Products{}).Where("id = ?", productID).Updates(product).Error
}

func (pr *productAdminRepository) GetProductByID(ctx context.Context, productID uuid.UUID) (*entities.Products, error) {
    if err := ctx.Err(); err != nil {
        return nil, err
    }
    var product entities.Products
    if err := pr.DB.Preload("ProductPricing").Preload("ProductVariants").Preload("ProductImages").Preload("ProductVideos").First(&product, "id = ?", productID).Error; err != nil {
        return nil, err
    }
    return &product, nil
}



func (pr *productAdminRepository) DeleteProduct(ctx context.Context, productID uuid.UUID) error {
    if err := ctx.Err(); err != nil {
        return err
    }

    // Start a transaction
    tx := pr.DB.Begin()
    if tx.Error != nil {
        return tx.Error
    }

    // Delete product pricing
    if err := tx.Where("product_id = ?", productID).Delete(&entities.ProductPricing{}).Error; err != nil {
        tx.Rollback()
        return err
    }

    // Delete product variants
    if err := tx.Where("product_id = ?", productID).Delete(&entities.ProductVariants{}).Error; err != nil {
        tx.Rollback()
        return err
    }

    // Delete product images
    if err := tx.Where("product_id = ?", productID).Delete(&entities.ProductImages{}).Error; err != nil {
        tx.Rollback()
        return err
    }

    // Delete product videos
    if err := tx.Where("product_id = ?", productID).Delete(&entities.ProductVideos{}).Error; err != nil {
        tx.Rollback()
        return err
    }

    // Delete main product
    if err := tx.Where("id = ?", productID).Delete(&entities.Products{}).Error; err != nil {
        tx.Rollback()
        return err
    }

    // Commit the transaction
    return tx.Commit().Error
}

