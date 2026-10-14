'use client';

import axios from 'axios';
import Link from 'next/link';
import { useEffect, useState } from 'react';
import DashboardProductCard from '@/components/dashboard/product-card';
import type { Product } from '@/types/product';

const API = 'http://localhost:8080';

type WishlistItem = {
  product_id: number;
  product_name: string;
  category_name: string;
  product_description?: string;
  description?: string;
  image_url?: string;
  product_image_url?: string;
  price: number;
  stock: number;
  popular: boolean;
  size: string;
  discount_percentage?: number;
  start_date?: string;
  end_date?: string;
};

const mapWishlistItem = (item: WishlistItem): Product => ({
  id: item.product_id,
  created_at: new Date().toISOString(),
  category_id: 0,
  category_name: item.category_name,
  product_name: item.product_name,
  product_description: item.product_description || item.description || '',
  product_image_url: item.product_image_url || item.image_url || '',
  price: Number(item.price || 0),
  stock: Number(item.stock || 0),
  popular: Boolean(item.popular),
  size: item.size || '',
  discount_percentage: Number(item.discount_percentage || 0),
  start_date: item.start_date || '',
  end_date: item.end_date || '',
});

export default function WishlistPage() {
  const [wishlistProducts, setWishlistProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState(true);

  const syncProducts = async () => {
    try {
      const response = await axios.get(`${API}/wishlist`, { withCredentials: true });
      const items = Array.isArray(response.data?.data) ? response.data.data : [];
      setWishlistProducts(items.map(mapWishlistItem));
    } catch (error) {
      console.error('Unable to load wishlist products', error);
      setWishlistProducts([]);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    syncProducts();
  }, []);

  const removeFromWishlist = async (productId: number) => {
    try {
      await axios.delete(`${API}/wishlist`, {
        withCredentials: true,
        data: { product_id: productId },
      });

      setWishlistProducts((currentProducts) => currentProducts.filter((product) => product.id !== productId));
    } catch (error) {
      console.error('Failed to remove product from wishlist', error);
    }
  };

  return (
    <div className="min-h-screen bg-slate-950 text-white">
      <header className="sticky top-0 z-20 border-b border-white/10 bg-slate-950/80 backdrop-blur-xl">
        <div className="mx-auto flex max-w-7xl items-center justify-between px-6 py-4">
          <div className="flex items-center gap-4">
            <Link href="/dashboard" className="text-lg font-semibold text-white hover:text-cyan-400">
              GenZone
            </Link>
          </div>

          <div className="flex items-center gap-3">
            <Link href="/dashboard" className="rounded-xl border border-white/10 bg-white/5 px-4 py-2 text-sm font-medium text-slate-100 hover:border-cyan-400/40 hover:bg-cyan-500/10">
              Back to Dashboard
            </Link>
          </div>
        </div>
      </header>

      <main className="mx-auto max-w-7xl px-6 py-10">
        <div className="mb-8 flex items-center justify-between gap-4">
          <div>
            <p className="text-sm uppercase tracking-[0.25em] text-red-300">Saved items</p>
            <h1 className="mt-2 text-3xl font-bold">My Wishlist</h1>
          </div>
          <div className="rounded-full border border-red-500/30 bg-red-500/10 px-4 py-2 text-sm text-red-200">
            {wishlistProducts.length} item{wishlistProducts.length === 1 ? '' : 's'}
          </div>
        </div>

        {loading ? (
          <div className="rounded-2xl border border-white/10 bg-white/5 p-10 text-slate-300">Loading wishlist...</div>
        ) : wishlistProducts.length === 0 ? (
          <div className="rounded-2xl border border-dashed border-slate-600 bg-white/5 p-12 text-center">
            <div className="mb-4 flex justify-center">
              <svg viewBox="0 0 24 24" fill="currentColor" className="h-12 w-12 text-red-400">
                <path d="M12 21.35 10.55 20C5.4 15.36 2 12.28 2 8.5A4.5 4.5 0 0 1 6.5 4c1.74 0 3.41.81 4.5 2.09A6.07 6.07 0 0 1 15.5 4 4.5 4.5 0 0 1 20 8.5c0 3.78-3.4 6.86-8.55 11.5L12 21.35Z" />
              </svg>
            </div>
            <h2 className="text-2xl font-semibold">Your wishlist is empty</h2>
            <p className="mt-3 text-slate-400">Tap the red heart on any product from the dashboard to save it here.</p>
            <Link href="/dashboard" className="mt-6 inline-flex rounded-xl bg-cyan-500 px-5 py-3 text-sm font-medium text-slate-950 hover:bg-cyan-400">
              Explore products
            </Link>
          </div>
        ) : (
          <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-4">
            {wishlistProducts.map((product) => (
              <div key={product.id} className="flex flex-col gap-3">
                <DashboardProductCard
                  product={product}
                  onAdded={() => {}}
                  onWishlistChange={(productId, isWishlisted) => {
                    if (!isWishlisted) {
                      removeFromWishlist(productId);
                    }
                  }}
                />

                <div className="flex justify-end">
                  <button
                    type="button"
                    onClick={() => removeFromWishlist(product.id)}
                    className="rounded-xl border border-red-500/40 bg-red-500/10 px-3 py-2 text-xs font-semibold uppercase tracking-[0.15em] text-red-200 transition hover:bg-red-500/20"
                  >
                    Remove
                  </button>
                </div>
              </div>
            ))}
          </div>
        )}
      </main>
    </div>
  );
}
