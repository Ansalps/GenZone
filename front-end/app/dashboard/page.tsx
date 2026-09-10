'use client';

import axios from 'axios';
import Link from 'next/link';
import { useEffect, useMemo, useState } from 'react';

interface Category {
  id: number;
  category_name: string;
  category_description: string;
  category_image_url: string;
}

interface Product {
  id: number;
  category_id: number;
  category_name: string;
  product_name: string;
  product_description: string;
  product_image_url: string;
  price: number;
  stock: number;
  popular: boolean;
  size: string;
  discount_percentage?: number;
}

const API = 'http://localhost:8080';

export default function DashboardPage() {
  const [categories, setCategories] = useState<Category[]>([]);
  const [products, setProducts] = useState<Product[]>([]);
  const [search, setSearch] = useState('');
  const [loadingCategories, setLoadingCategories] = useState(true);
  const [loadingProducts, setLoadingProducts] = useState(true);
  const [cartCount, setCartCount] = useState(0);

  useEffect(() => {
    const fetchCategories = async () => {
      try {
        const response = await axios.get(`${API}/public/category`, { withCredentials: true });
        if (response.data?.status && Array.isArray(response.data?.data?.categories)) {
          setCategories(response.data.data.categories);
        }
      } catch (error) {
        console.error('Failed to load categories', error);
      } finally {
        setLoadingCategories(false);
      }
    };

    const fetchProducts = async () => {
      try {
        setLoadingProducts(true);
        const response = await axios.get(`${API}/public/product`, {
          withCredentials: true,
          params: { search: search || undefined },
        });

        if (response.data?.status && Array.isArray(response.data?.data?.products)) {
          setProducts(response.data.data.products);
        } else {
          setProducts([]);
        }
      } catch (error) {
        console.error('Failed to load products', error);
        setProducts([]);
      } finally {
        setLoadingProducts(false);
      }
    };

    fetchCategories();
    fetchProducts();
  }, [search]);

  useEffect(() => {
    const fetchCart = async () => {
      try {
        const response = await axios.get(`${API}/cart`, { withCredentials: true });
        const cartItems = response.data?.data?.cart_items ?? [];
        setCartCount(Array.isArray(cartItems) ? cartItems.length : 0);
      } catch (error) {
        console.error('Unable to load cart', error);
        setCartCount(0);
      }
    };

    fetchCart();
  }, []);

  const filteredProducts = useMemo(() => {
    const term = search.trim().toLowerCase();
    if (!term) return products;

    return products.filter((product) => {
      return (
        product.product_name.toLowerCase().includes(term) ||
        product.product_description.toLowerCase().includes(term) ||
        product.category_name.toLowerCase().includes(term)
      );
    });
  }, [products, search]);

  const addToCart = async (productId: number) => {
    try {
      await axios.post(
        `${API}/cart`,
        { product_id: String(productId) },
        { withCredentials: true }
      );

      const response = await axios.get(`${API}/cart`, { withCredentials: true });
      const cartItems = response.data?.data?.cart_items ?? [];
      setCartCount(Array.isArray(cartItems) ? cartItems.length : 0);
      alert('Product added to cart');
    } catch (error) {
      console.error('Failed to add to cart', error);
      alert('Unable to add product to cart');
    }
  };

  return (
    <div className="min-h-screen bg-slate-100 text-slate-800">
      <header className="bg-green-800 text-white shadow-md">
        <div className="mx-auto flex max-w-7xl items-center justify-between px-6 py-4">
          <div>
            <h1 className="text-2xl font-bold">GenZone</h1>
          </div>

          <div className="flex items-center gap-4">
            <div className="relative">
              <span className="text-sm">Cart</span>
              <span className="ml-2 rounded-full bg-white px-2 py-1 text-xs font-bold text-green-800">
                {cartCount}
              </span>
            </div>
            <Link
              href="/cart"
              className="rounded-lg border border-white px-4 py-2 text-sm font-medium hover:bg-white hover:text-green-800"
            >
              View Cart
            </Link>
          </div>
        </div>
      </header>

      <main className="mx-auto max-w-7xl px-6 py-8">
        <section className="mb-8 rounded-2xl bg-gradient-to-r from-green-700 to-emerald-600 p-8 text-white">
          <p className="mb-2 text-sm uppercase tracking-[0.2em] text-green-100">Welcome back</p>
          <h2 className="text-3xl font-bold">Shop the latest products</h2>
        </section>

        <section className="mb-8">
          <h3 className="mb-4 text-2xl font-bold">Categories</h3>
          {loadingCategories ? (
            <div className="text-slate-500">Loading categories...</div>
          ) : categories.length === 0 ? (
            <div className="text-slate-500">No categories available.</div>
          ) : (
            <div className="grid gap-5 sm:grid-cols-2 lg:grid-cols-4">
              {categories.map((category) => (
                <div key={category.id} className="overflow-hidden rounded-xl bg-white shadow-sm">
                  <img src={category.category_image_url} alt={category.category_name} className="h-40 w-full object-cover" />
                  <div className="p-4">
                    <h4 className="text-lg font-semibold">{category.category_name}</h4>
                    <p className="mt-2 text-sm text-slate-600">{category.category_description}</p>
                  </div>
                </div>
              ))}
            </div>
          )}
        </section>

        <section>
          <div className="mb-5 flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
            <h3 className="text-2xl font-bold">Products</h3>
            <input
              type="text"
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              placeholder="Search products..."
              className="w-full rounded-lg border border-slate-300 bg-white px-4 py-2 outline-none ring-0 md:max-w-md"
            />
          </div>

          {loadingProducts ? (
            <div className="text-slate-500">Loading products...</div>
          ) : filteredProducts.length === 0 ? (
            <div className="rounded-xl border border-dashed border-slate-300 bg-white p-12 text-center text-slate-500">
              No matching products found.
            </div>
          ) : (
            <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-4">
              {filteredProducts.map((product) => {
                const discountPercentage = Number(product.discount_percentage) || 0;
                const hasDiscount = discountPercentage > 0;
                const discountAmount = hasDiscount ? product.price * (discountPercentage / 100) : 0;
                const finalPrice = product.price - discountAmount;

                return (
                  <div key={product.id} className="overflow-hidden rounded-xl bg-white shadow-sm">
                    <img src={product.product_image_url} alt={product.product_name} className="h-52 w-full object-cover" />
                    <div className="p-4">
                      <p className="text-xs uppercase tracking-wide text-slate-500">{product.category_name}</p>
                      <h4 className="mt-2 text-lg font-semibold">{product.product_name}</h4>
                      <p className="mt-2 line-clamp-3 text-sm text-slate-600">{product.product_description}</p>

                      <div className="mt-4 flex items-center gap-2">
                        {hasDiscount ? (
                          <>
                            <span className="text-xl font-bold text-green-700">₹{finalPrice.toFixed(2)}</span>
                            <span className="text-sm text-slate-400 line-through">₹{product.price.toFixed(2)}</span>
                          </>
                        ) : (
                          <span className="text-xl font-bold text-green-700">₹{product.price.toFixed(2)}</span>
                        )}
                      </div>

                      <div className="mt-3 text-sm text-slate-600">Stock: {product.stock}</div>

                      <button
                        type="button"
                        onClick={() => addToCart(product.id)}
                        className="mt-5 w-full rounded-lg bg-green-700 px-4 py-2 font-medium text-white hover:bg-green-800"
                      >
                        Add to cart
                      </button>
                    </div>
                  </div>
                );
              })}
            </div>
          )}
        </section>
      </main>
    </div>
  );
}
