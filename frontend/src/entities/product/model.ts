export interface Product {
  id: string;
  external_id: string;
  sku: string;
  name: string;
  original_price: number;
  price: number;
  image_url: string;
  product_url: string;
  brand: string;
  category_id: string;
  created_at: string;
  updated_at: string;
}

export interface ProductList {
  products: Product[] | null;
  total: number;
  page: number;
  page_size: number;
}

export interface ProductFilter {
  page?: number;
  page_size?: number;
  sort_fields?: string;
  category_id?: string;
  brand?: string;
  min_price?: number;
  max_price?: number;
  search?: string;
}
