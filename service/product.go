// service/product.go
package service

import (
	"context"
	"douyin/repository/cache"
	"douyin/repository/db/dao"
	"douyin/repository/db/model"
	"douyin/types"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"
)

func CreateProduct(ctx context.Context, userID uint32, product *types.Product) error {
	// 1. 验证用户身份是否有效，或者是否有权限创建商品
	if userID == 0 {
		return errors.New("用户身份无效")
	}

	// 2. 将 types.Product 转换为 model.Product
	modelProduct := &model.Product{
		Name:        product.Name,
		Description: product.Description,
		Picture:     product.Picture,
		Price:       product.Price,
		// Stock and Version will be set by GORM default or manually if needed
	}

	// 3. 调用DAO层，进行商品创建
	err := dao.CreateProduct(modelProduct) // Assuming DAO CreateProduct handles context if necessary, or pass ctx
	if err != nil {
		log.Printf("创建商品失败：%v", err)
		return err
	}
	fmt.Println("商品创建成功")
	// Consider invalidating product list caches here
	// For now, log or skip as per subtask instructions for list invalidation.
	log.Println("Product list cache invalidation would be needed after creating a product.")
	return nil
}

func GetProductByID(ctx context.Context, id uint32) (*types.Product, error) {
	key := cache.ProductDetailKey(uint(id))
	// Try to get from cache
	cachedData, err := cache.RedisClient.Get(ctx, key).Result()
	if err == nil {
		var productModel model.Product
		if errUnmarshal := json.Unmarshal([]byte(cachedData), &productModel); errUnmarshal == nil {
			// Cache hit and unmarshal success
			return &types.Product{
				ID:          uint32(productModel.ID),
				Name:        productModel.Name,
				Description: productModel.Description,
				Picture:     productModel.Picture,
				Price:       productModel.Price,
				Stock:       productModel.Stock, // Assuming types.Product also has Stock and Version
				Version:     productModel.Version,
			}, nil
		}
		// Unmarshal failed, treat as cache miss and delete potentially corrupt cache entry
		log.Printf("Error unmarshalling cached product for ID %d: %v. Deleting cache entry.", id, errUnmarshal)
		cache.RedisClient.Del(ctx, key) // Fire and forget deletion
	} else if err != cache.Nil {
		// Some other Redis error
		log.Printf("Redis error getting product for ID %d: %v", id, err)
	}

	// Cache miss or error, fetch from DB
	product, dbErr := dao.GetProduct(id) // Assuming DAO GetProduct doesn't need context or handles it
	if dbErr != nil {
		log.Printf("查询商品失败：%v", dbErr)
		return nil, dbErr
	}

	// Marshal model.Product to JSON for caching
	jsonData, marshalErr := json.Marshal(product)
	if marshalErr != nil {
		log.Printf("Error marshalling product for ID %d for cache: %v", id, marshalErr)
		// Proceed without caching if marshalling fails
	} else {
		expireDuration := time.Minute*10 + time.Duration(rand.Intn(300))*time.Second
		if setErr := cache.RedisClient.Set(ctx, key, jsonData, expireDuration).Err(); setErr != nil {
			log.Printf("Error setting product cache for ID %d: %v", id, setErr)
		}
	}

	// Convert model.Product to types.Product
	typesProduct := &types.Product{
		ID:          uint32(product.ID),
		Name:        product.Name,
		Description: product.Description,
		Picture:     product.Picture,
		Price:       product.Price,
		Stock:       product.Stock,
		Version:     product.Version,
	}

	return typesProduct, nil
}

func ListProducts(ctx context.Context, pageNum, pageSize int) ([]types.Product, int64, error) {
	key := cache.ProductListKey(pageNum, pageSize)
	// Try to get from cache
	cachedData, err := cache.RedisClient.Get(ctx, key).Result()
	if err == nil {
		var res struct {
			Total    int64            `json:"total"`
			Products []model.Product `json:"products"`
		}
		if errUnmarshal := json.Unmarshal([]byte(cachedData), &res); errUnmarshal == nil {
			// Cache hit and unmarshal success
			var typesProducts []types.Product
			for _, pModel := range res.Products {
				typesProducts = append(typesProducts, types.Product{
					ID:          uint32(pModel.ID),
					Name:        pModel.Name,
					Description: pModel.Description,
					Picture:     pModel.Picture,
					Price:       pModel.Price,
					Stock:       pModel.Stock,
					Version:     pModel.Version,
				})
			}
			return typesProducts, res.Total, nil
		}
		// Unmarshal failed
		log.Printf("Error unmarshalling cached product list for page %d, size %d: %v. Deleting cache entry.", pageNum, pageSize, errUnmarshal)
		cache.RedisClient.Del(ctx, key)
	} else if err != cache.Nil {
		log.Printf("Redis error getting product list for page %d, size %d: %v", pageNum, pageSize, err)
	}

	// Cache miss or error, fetch from DB
	productsFromDB, total, dbErr := dao.ListProducts(pageNum, pageSize) // Assuming DAO ListProducts doesn't need context
	if dbErr != nil {
		log.Printf("获取商品列表失败：%v", dbErr)
		return nil, 0, dbErr
	}

	// Prepare data for caching
	cachePayload := map[string]interface{}{
		"total":    total,
		"products": productsFromDB, // Cache the model.Product slice
	}
	jsonData, marshalErr := json.Marshal(cachePayload)
	if marshalErr != nil {
		log.Printf("Error marshalling product list for page %d, size %d for cache: %v", pageNum, pageSize, marshalErr)
	} else {
		expireDuration := time.Minute*5 + time.Duration(rand.Intn(120))*time.Second
		if setErr := cache.RedisClient.Set(ctx, key, jsonData, expireDuration).Err(); setErr != nil {
			log.Printf("Error setting product list cache for page %d, size %d: %v", pageNum, pageSize, setErr)
		}
	}

	// Convert model.Product to types.Product for return
	var result []types.Product
	for _, p := range productsFromDB {
		result = append(result, types.Product{
			ID:          uint32(p.ID),
			Name:        p.Name,
			Description: p.Description,
			Picture:     p.Picture,
			Price:       p.Price,
			Stock:       p.Stock,
			Version:     p.Version,
		})
	}
	return result, total, nil
}

func UpdateProduct(ctx context.Context, userID uint32, product *types.Product) error {
	// 验证用户身份或权限
	if userID == 0 {
		return errors.New("用户身份无效")
	}

	// 2. 将 types.Product 转换为 model.Product
	modelProduct := &model.Product{
		ID:          uint(product.ID),
		Name:        product.Name,
		Description: product.Description,
		Picture:     product.Picture,
		Price:       product.Price,
		Stock:       product.Stock,   // Assuming types.Product carries this for updates
		Version:     product.Version, // And version for optimistic locking if applicable at DAO
	}

	// 调用DAO层修改商品
	err := dao.UpdateProduct(modelProduct) // Pass ctx if DAO method is updated
	if err != nil {
		log.Printf("修改商品失败：%v", err)
		return err
	}

	fmt.Println("商品信息修改成功")
	// Invalidate product detail cache
	detailKey := cache.ProductDetailKey(uint(product.ID))
	if delErr := cache.RedisClient.Del(ctx, detailKey).Err(); delErr != nil {
		log.Printf("Error deleting product detail cache for ID %d: %v", product.ID, delErr)
	}
	// List cache invalidation - log or skip as per subtask
	log.Println("Product list cache invalidation would be needed after updating a product.")
	return nil
}

func DeleteProduct(ctx context.Context, userID uint32, productID uint32) error {
	// 验证用户身份或权限
	if userID == 0 {
		return errors.New("用户身份无效")
	}

	// 调用DAO层删除商品
	err := dao.DeleteProduct(productID) // Pass ctx if DAO method is updated
	if err != nil {
		log.Printf("删除商品失败：%v", err)
		return err
	}

	fmt.Println("商品删除成功")
	// Invalidate product detail cache
	detailKey := cache.ProductDetailKey(uint(productID))
	if delErr := cache.RedisClient.Del(ctx, detailKey).Err(); delErr != nil {
		log.Printf("Error deleting product detail cache for ID %d: %v", productID, delErr)
	}
	// List cache invalidation - log or skip as per subtask
	log.Println("Product list cache invalidation would be needed after deleting a product.")
	return nil
}
