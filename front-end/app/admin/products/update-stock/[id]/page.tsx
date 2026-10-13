'use client';

import axios from 'axios';
import Link from 'next/link';
import { useEffect, useState } from 'react';
import { useParams, useRouter } from 'next/navigation';
import { toast } from 'sonner';

const API = process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8080';

const sizeOrder = ['Small', 'Medium', 'Large'];

export default function UpdateStockPage() {
  const params = useParams();
  const router = useRouter();
  const productId = params?.id as string;

  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [productName, setProductName] = useState('');
  const [inventory, setInventory] = useState<Record<string, number>>({
    Small: 0,
    Medium: 0,
    Large: 0,
  });

  useEffect(() => {
    const fetchProduct = async () => {
      try {
        const { data } = await axios.get(`${API}/admin/product/${productId}`, { withCredentials: true });
        const product = data?.data ?? {};
        const loaded: Record<string, number> = { Small: 0, Medium: 0, Large: 0 };

        (product.inventory ?? []).forEach((variant: { size: string; stock: number }) => {
          if (variant.size) loaded[variant.size] = Number(variant.stock ?? 0);
        });

        setProductName(product.product_name || 'Product');
        setInventory(loaded);
      } catch (error) {
        console.error('Failed to fetch product stock', error);
        toast.error('Failed to load product stock');
      } finally {
        setLoading(false);
      }
    };

    if (productId) {
      void fetchProduct();
    }
  }, [productId]);

  const onStockChange = (size: string, value: number) => {
    setInventory((prev) => ({
      ...prev,
      [size]: Number.isFinite(value) ? Math.max(0, value) : 0,
    }));
  };

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setSaving(true);

    try {
      const payload = {
        inventory: sizeOrder.map((size) => ({
          size,
          stock: inventory[size] ?? 0,
        })),
      };

      await axios.put(`${API}/admin/product/${productId}/stock`, payload, { withCredentials: true });
      toast.success('Stock updated successfully');
      router.push('/admin/products');
    } catch (error: any) {
      console.error('Update stock failed', error);
      toast.error(error?.response?.data?.message || 'Failed to update stock');
    } finally {
      setSaving(false);
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-slate-950 px-6 py-10 text-white">
        <div className="mx-auto max-w-xl rounded-3xl border border-white/10 bg-slate-900/80 p-8 text-center">
          Loading stock information...
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-slate-950 px-4 py-8 text-white sm:px-6 lg:px-8">
      <div className="mx-auto max-w-2xl">
        <header className="mb-6 flex items-center justify-between gap-4 rounded-3xl border border-white/10 bg-white/5 p-5 shadow-2xl shadow-slate-950/30 backdrop-blur-xl">
          <div>
            <p className="text-xs font-semibold uppercase tracking-[0.28em] text-violet-300">Inventory</p>
            <h1 className="mt-2 text-2xl font-bold">Update Stock for {productName}</h1>
          </div>
          <Link href="/admin/products" className="rounded-xl border border-white/10 bg-slate-900/60 px-3 py-2 text-sm font-medium text-slate-200 hover:bg-slate-900 transition">
            Back
          </Link>
        </header>

        <main className="rounded-3xl border border-white/10 bg-slate-900/60 p-6 shadow-2xl shadow-slate-950/30">
          <form onSubmit={handleSubmit} className="space-y-5">
            <div className="grid grid-cols-1 gap-4 md:grid-cols-3">
              {sizeOrder.map((size) => (
                <div key={size} className="rounded-2xl border border-white/10 bg-slate-950/40 p-4">
                  <label className="mb-2 block text-sm font-medium text-slate-200">{size}</label>
                  <input
                    type="number"
                    min={0}
                    value={inventory[size] ?? 0}
                    onChange={(e) => onStockChange(size, Number(e.target.value) || 0)}
                    className="w-full rounded-xl border border-white/10 bg-slate-900/50 px-3 py-2 text-white outline-none"
                  />
                </div>
              ))}
            </div>

            <button
              type="submit"
              disabled={saving}
              className="w-full rounded-xl bg-gradient-to-r from-violet-500 to-purple-500 px-4 py-3 text-sm font-semibold text-white shadow-lg shadow-violet-500/20 transition hover:brightness-110 disabled:opacity-60"
            >
              {saving ? 'Updating...' : 'Update Stock'}
            </button>
          </form>
        </main>
      </div>
    </div>
  );
}
