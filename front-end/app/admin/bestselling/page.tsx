'use client';

import axios from 'axios';
import Link from 'next/link';
import { useEffect, useState } from 'react';
import { toast } from 'sonner';

const API = 'http://localhost:8080';

type ProductRecord = {
  id?: number;
  product_id?: number;
  product_name?: string;
  category_name?: string;
  price?: number;
  stock?: number;
  category?: string;
};

const money = (value: number = 0) =>
  new Intl.NumberFormat('en-IN', {
    style: 'currency',
    currency: 'INR',
    maximumFractionDigits: 2,
  }).format(value);

export default function BestSellingPage() {
  const [products, setProducts] = useState<ProductRecord[]>([]);
  const [categories, setCategories] = useState<string[]>([]);
  const [filters, setFilters] = useState({ day: '', month: '', year: '' });
  const [loading, setLoading] = useState(false);

  const loadBestSelling = async (nextFilters = filters) => {
    setLoading(true);

    try {
      const params = {
        ...(nextFilters.day ? { day: nextFilters.day } : {}),
        ...(nextFilters.month ? { month: nextFilters.month } : {}),
        ...(nextFilters.year ? { year: nextFilters.year } : {}),
      };

      const response = await axios.get(`${API}/admin/bestselling`, {
        params,
        withCredentials: true,
      });

      const bestProducts = Array.isArray(response.data?.best_selling_product) ? response.data.best_selling_product : [];
      const bestCategories = Array.isArray(response.data?.best_selling_category) ? response.data.best_selling_category : [];
      setProducts(bestProducts);
      setCategories(bestCategories);
    } catch (error) {
      console.error('Failed to load best selling data:', error);
      setProducts([]);
      setCategories([]);
      toast.error('Unable to load best selling products');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadBestSelling();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const handleApply = async () => {
    await loadBestSelling(filters);
  };

  const handleReset = async () => {
    const resetValues = { day: '', month: '', year: '' };
    setFilters(resetValues);
    await loadBestSelling(resetValues);
  };

  return (
    <div className="min-h-screen bg-slate-950 px-4 py-8 text-white sm:px-6 lg:px-8">
      <div className="mx-auto max-w-7xl">
        <header className="mb-6 flex flex-col gap-4 rounded-3xl border border-white/10 bg-white/5 p-5 backdrop-blur-xl sm:flex-row sm:items-center sm:justify-between">
          <div className="flex items-center gap-4">
            <Link href="/admin" className="rounded-full border border-white/10 bg-slate-900/60 px-3 py-2 text-sm font-medium text-slate-200 transition hover:bg-slate-900">
              ← Back
            </Link>
            <div>
              <p className="text-sm uppercase tracking-[0.2em] text-yellow-400">Insights</p>
              <h1 className="text-2xl font-bold">Best Selling</h1>
            </div>
          </div>
        </header>

        <section className="mb-6 grid gap-4 rounded-3xl border border-white/10 bg-slate-900/60 p-5 shadow-xl shadow-slate-950/30 md:grid-cols-4">
          <label className="text-sm text-slate-300">
            Day
            <input
              type="date"
              value={filters.day}
              onChange={(e) => setFilters((prev) => ({ ...prev, day: e.target.value }))}
              className="mt-2 w-full rounded-xl border border-white/10 bg-slate-950 px-3 py-2 text-white outline-none"
            />
          </label>
          <label className="text-sm text-slate-300">
            Month
            <input
              type="number"
              min={1}
              max={12}
              value={filters.month}
              onChange={(e) => setFilters((prev) => ({ ...prev, month: e.target.value, year: e.target.value ? '' : prev.year }))}
              placeholder="1-12"
              className="mt-2 w-full rounded-xl border border-white/10 bg-slate-950 px-3 py-2 text-white outline-none"
            />
          </label>
          <label className="text-sm text-slate-300">
            Year
            <input
              type="number"
              value={filters.year}
              onChange={(e) => setFilters((prev) => ({ ...prev, year: e.target.value, month: e.target.value ? '' : prev.month }))}
              placeholder="2026"
              className="mt-2 w-full rounded-xl border border-white/10 bg-slate-950 px-3 py-2 text-white outline-none"
            />
          </label>
          <div className="flex items-end gap-2">
            <button
              type="button"
              onClick={handleApply}
              className="flex-1 rounded-xl bg-white/10 px-4 py-2 text-sm font-medium text-white transition hover:bg-white/15"
            >
              Apply
            </button>
            <button
              type="button"
              onClick={handleReset}
              className="flex-1 rounded-xl border border-white/10 bg-slate-900 px-4 py-2 text-sm font-medium text-slate-200 transition hover:bg-slate-800"
            >
              Reset
            </button>
          </div>
        </section>

        <section className="mb-6 grid gap-4 md:grid-cols-2">
          <div className="rounded-3xl border border-white/10 bg-slate-900/60 p-5 shadow-xl shadow-slate-950/30">
            <h2 className="mb-4 text-lg font-semibold text-yellow-300">Top products</h2>
            {loading ? (
              <p className="text-sm text-slate-400">Loading...</p>
            ) : products.length === 0 ? (
              <p className="text-sm text-slate-400">No product data available.</p>
            ) : (
              <div className="space-y-3">
                {products.map((product, index) => (
                  <div key={`${product.id ?? product.product_id ?? 'product'}-${index}`} className="flex items-center justify-between rounded-xl border border-white/10 bg-slate-950/50 p-3">
                    <div>
                      <p className="font-medium text-slate-100">{index + 1}. {product.product_name ?? 'Unnamed product'}</p>
                      <p className="text-xs text-slate-400">{product.category_name ?? product.category ?? 'Uncategorized'}</p>
                    </div>
                    <div className="text-right">
                      <p className="font-semibold text-emerald-300">{money(Number(product.price ?? 0))}</p>
                      <p className="text-xs text-slate-400">Stock: {product.stock ?? 0}</p>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>

          <div className="rounded-3xl border border-white/10 bg-slate-900/60 p-5 shadow-xl shadow-slate-950/30">
            <h2 className="mb-4 text-lg font-semibold text-cyan-300">Top categories</h2>
            {categories.length === 0 ? (
              <p className="text-sm text-slate-400">No category data available.</p>
            ) : (
              <div className="space-y-3">
                {categories.map((category, index) => (
                  <div key={`${category}-${index}`} className="flex items-center justify-between rounded-xl border border-white/10 bg-slate-950/50 p-3">
                    <div className="font-medium text-slate-100">{index + 1}. {category || 'Unknown category'}</div>
                    <div className="rounded-full border border-yellow-400/30 bg-yellow-500/10 px-2 py-1 text-xs text-yellow-100">Top seller</div>
                  </div>
                ))}
              </div>
            )}
          </div>
        </section>
      </div>
    </div>
  );
}
