'use client';

import axios from 'axios';
import Link from 'next/link';
import { useEffect, useState } from 'react';
import DashboardProductCard from '@/components/dashboard/product-card'

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
  const [categorySearch, setCategorySearch] = useState('');
  const [categoryNameSort, setCategoryNameSort] = useState('');
  const [productSearch, setProductSearch] = useState('');
  const [productNameSort, setProductNameSort] = useState('');
  const [productPriceSort, setProductPriceSort] = useState('');
  const [loadingCategories, setLoadingCategories] = useState(true);
  const [loadingProducts, setLoadingProducts] = useState(true);
  const [cartCount, setCartCount] = useState(0);

  useEffect(() => {
    const fetchCategories = async () => {
      try {
        setLoadingCategories(true);
        const response = await axios.get(`${API}/public/category`, {
          withCredentials: true,
          params: {
            search: categorySearch || undefined,
            name_sort: categoryNameSort || undefined,
          },
        });

        if (response.data?.status && Array.isArray(response.data?.data?.categories)) {
          setCategories(response.data.data.categories);
        } else {
          setCategories([]);
        }
      } catch (error) {
        console.error('Failed to load categories', error);
        setCategories([]);
      } finally {
        setLoadingCategories(false);
      }
    };

    fetchCategories();
  }, [categorySearch, categoryNameSort]);

  useEffect(() => {
    const fetchProducts = async () => {
      try {
        setLoadingProducts(true);
        const response = await axios.get(`${API}/public/product`, {
          withCredentials: true,
          params: {
            search: productSearch || undefined,
            name_sort: productNameSort || undefined,
            price_sort: productPriceSort || undefined,
          },
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

    fetchProducts();
  }, [productSearch, productNameSort, productPriceSort]);

  useEffect(() => {
    const fetchCart = async () => {
      try {
        const response = await axios.get(`${API}/cart-total-quantity`, { withCredentials: true });
        const totalQuantity = response.data?.total_quantity ?? 0;
        setCartCount(totalQuantity);
      } catch (error) {
        console.error('Unable to load cart', error);
        setCartCount(0);
      }
    };

    fetchCart();
  }, []);

  const addToCart = async (productId: number) => {
    try {
      await axios.post(
        `${API}/cart`,
        { product_id: String(productId) },
        { withCredentials: true }
      );

      const response = await axios.get(`${API}/cart-total-quantity`, { withCredentials: true });
      console.log('Cart total quantity response:', response.data);
      const totalQuantity = response.data?.total_quantity ?? 0;
      setCartCount(totalQuantity);
      alert('Product added to cart');
    } catch (error) {
      console.error('Failed to add to cart', error);
      alert('Unable to add product to cart');
    }
  };

  return (
    <div className="min-h-screen bg-slate-950 text-white">
      <header className="sticky top-0 z-20 border-b border-white/10 bg-slate-950/80 backdrop-blur-xl">
        <div className="mx-auto flex max-w-7xl items-center justify-between px-6 py-4">
          <div>
            <h1 className="text-2xl font-bold">GenZone</h1>
          </div>

          <div className="flex items-center gap-4">
            <div className="relative">
              <span className="text-sm">Cart</span>
              <span className="ml-2 rounded-full bg-white px-2 py-1 text-xs font-bold text-slate-950">
                {cartCount}
              </span>
            </div>
            <Link
              href="/cart"
              className="rounded-xl border border-white/10 bg-white/5 px-4 py-2 text-sm font-medium text-slate-100 transition hover:border-cyan-400/40 hover:bg-cyan-500/10"
            >
              View Cart
            </Link>
          </div>
        </div>
      </header>

      <main className="mx-auto max-w-7xl px-6 py-8">
        <section className="mb-8 rounded-2xl bg-gradient-to-br from-cyan-500/20 to-violet-500/20 p-8 text-white">
          <p className="mb-2 text-sm uppercase tracking-[0.2em] text-cyan-300">Welcome back</p>
          <h2 className="text-3xl font-bold">Shop the latest products</h2>
        </section>

        <section className="mb-8 rounded-2xl bg-white/5 p-4 shadow-sm">
          <div className="grid gap-4 md:grid-cols-[1.5fr_0.8fr]">
            <div>
              <label className="mb-2 block text-xs font-semibold uppercase tracking-[0.2em] text-slate-500">Search categories</label>
              <input
                type="text"
                value={categorySearch}
                onChange={(e) => setCategorySearch(e.target.value)}
                placeholder="Search categories..."
                className="w-full rounded-lg border border-white/10 bg-slate-900/80 px-4 py-2 text-white placeholder:text-slate-400 outline-none ring-0"
              />
            </div>

            <div>
              <label className="mb-2 block text-xs font-semibold uppercase tracking-[0.2em] text-slate-500">Category sort</label>
              <select
                value={categoryNameSort}
                onChange={(e) => setCategoryNameSort(e.target.value)}
                className="w-full rounded-lg border border-white/10 bg-slate-900/80 px-4 py-2 text-white outline-none ring-0"
              >
                <option value="">Default</option>
                <option value="aA-zZ">A to Z</option>
                <option value="zZ-aA">Z to A</option>
              </select>
            </div>
          </div>
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
                <div key={category.id} className="overflow-hidden rounded-xl bg-white/5 shadow-sm">
                  <img src={category.category_image_url} alt={category.category_name} className="h-40 w-full object-cover" />
                  <div className="p-4">
                    <h4 className="text-lg font-semibold">{category.category_name}</h4>
                    <p className="mt-2 text-sm text-slate-400">{category.category_description}</p>
                  </div>
                </div>
              ))}
            </div>
          )}
        </section>

        <section className="rounded-2xl bg-white/5 p-4 shadow-sm">
          <div className="grid gap-4 md:grid-cols-[1.3fr_0.8fr_0.8fr]">
            <div>
              <label className="mb-2 block text-xs font-semibold uppercase tracking-[0.2em] text-slate-500">Search products</label>
              <input
                type="text"
                value={productSearch}
                onChange={(e) => setProductSearch(e.target.value)}
                placeholder="Search products..."
                className="w-full rounded-lg border border-white/10 bg-slate-900/80 px-4 py-2 text-white placeholder:text-slate-400 outline-none ring-0"
              />
            </div>

            <div>
              <label className="mb-2 block text-xs font-semibold uppercase tracking-[0.2em] text-slate-500">Name sort</label>
              <select
                value={productNameSort}
                onChange={(e) => setProductNameSort(e.target.value)}
                className="w-full rounded-lg border border-white/10 bg-slate-900/80 px-4 py-2 text-white outline-none ring-0"
              >
                <option value="">Default</option>
                <option value="aA-zZ">A to Z</option>
                <option value="zZ-aA">Z to A</option>
              </select>
            </div>

            <div>
              <label className="mb-2 block text-xs font-semibold uppercase tracking-[0.2em] text-slate-500">Price sort</label>
              <select
                value={productPriceSort}
                onChange={(e) => setProductPriceSort(e.target.value)}
                className="w-full rounded-lg border border-white/10 bg-slate-900/80 px-4 py-2 text-white outline-none ring-0"
              >
                <option value="">Default</option>
                <option value="low-high">Low to High</option>
                <option value="high-low">High to Low</option>
              </select>
            </div>
          </div>
        </section>

        <section className="mt-8">
          <div className="mb-5 flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
            <h3 className="text-2xl font-bold">Products</h3>
          </div>

          {loadingProducts ? (
            <div className="text-slate-500">Loading products...</div>
          ) : products.length === 0 ? (
            <div className="rounded-xl border border-dashed border-slate-300 bg-white p-12 text-center text-slate-500">
              No matching products found.
            </div>
          ) : (
            <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-4">
              {products.map((product) => (
                <DashboardProductCard key={product.id} product={product as any} onAdded={async (success) => {
                    if (success) {
                        try {
                            const response = await axios.get(`${API}/cart-total-quantity`, { withCredentials: true })
                            const totalQuantity = response.data?.total_quantity ?? 0
                            setCartCount(totalQuantity)
                        } catch (err) {
                            console.error('Unable to refresh cart', err)
                        }
                    }
                }} />
              ))}
            </div>
          )}
        </section>
      </main>
    </div>
  );
}
