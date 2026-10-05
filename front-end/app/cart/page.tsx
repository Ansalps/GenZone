'use client';

import axios from 'axios';
import Link from 'next/link';
import { useEffect, useState } from 'react';

interface CartItem {
  user_id?: number | string;
  product_id?: number | string;
  product_name?: string;
  total_amount?: number; // original total (price * qty)
  qty?: number;
  price?: number; // unit price (after discount if applied)
  discount?: number; // total discount amount for this line
  final_amount?: number; // final total after discount
}

export default function CartPage() {
  const [items, setItems] = useState<CartItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [removing, setRemoving] = useState<string | null>(null);

  useEffect(() => {
    loadCart();
  }, []);

  const subtotal = items.reduce((sum, item) => sum + Number(item.final_amount ?? item.total_amount ?? 0), 0);

  const handleRemove = async (productId?: number | string, qty?: number) => {
    if (!productId) return;
    const should = window.confirm('Remove this item from cart?');
    if (!should) return;
    try {
      setRemoving(String(productId));
      await axios.delete('http://localhost:8080/cart', {
        data: { product_id: productId, quantity: qty ?? 0 },
        withCredentials: true,
      });
      await loadCart();
    } catch (err) {
      console.error('Failed to remove item', err);
    } finally {
      setRemoving(null);
    }
  };

  // extracted so it can be reused after mutations
  const loadCart = async () => {
    setLoading(true);
    try {
      const response = await axios.get('http://localhost:8080/cart', { withCredentials: true });
      const cartItems = response.data?.data?.cart_items ?? [];
      setItems(Array.isArray(cartItems) ? cartItems : []);
    } catch (error) {
      console.error('Failed to load cart', error);
      setItems([]);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-slate-950 text-white">
      <div className="mx-auto max-w-5xl p-6">
        <div className="mb-6 flex items-center justify-between">
          <h1 className="text-3xl font-bold text-white">Your Cart</h1>
          <Link href="/dashboard" className="rounded-xl border border-white/10 bg-white/5 px-4 py-2 text-sm font-medium text-slate-100 transition hover:border-cyan-400/40 hover:bg-cyan-500/10">
            Continue Shopping
          </Link>
        </div>

        {loading ? (
          <div className="rounded-xl bg-slate-900/80 p-8 text-slate-300">Loading cart...</div>
        ) : items.length === 0 ? (
          <div className="rounded-xl bg-slate-900/80 p-12 text-center text-slate-300 shadow-sm">
            <p className="text-xl font-semibold text-white">Your cart is empty.</p>
            <Link href="/dashboard" className="mt-4 inline-block rounded-xl border border-white/10 bg-white/5 px-4 py-2 text-sm font-medium text-slate-100 transition hover:border-cyan-400/40 hover:bg-cyan-500/10">
              Browse Products
            </Link>
          </div>
        ) : (
          <div className="grid gap-6 lg:grid-cols-[2fr_1fr]">
            <div className="space-y-4">
              {items.map((item, index) => (
                <div key={`${item.product_id ?? index}`} className="flex items-center justify-between rounded-xl bg-slate-900/80 p-4 shadow-sm">
                  <div>
                    <h2 className="text-lg font-semibold text-white">{item.product_name ?? 'Product'}</h2>
                    <p className="text-sm text-slate-400">Qty: {item.qty ?? 1}</p>
                    <p className="text-sm text-slate-400">Unit: ₹{Number(item.price ?? 0).toFixed(2)}</p>
                    {Number(item.discount ?? 0) > 0 && (
                      <p className="text-xs text-rose-500">You saved: ₹{Number(item.discount).toFixed(2)}</p>
                    )}
                  </div>
                  <div className="text-right flex flex-col items-end gap-2">
                    <p className="font-semibold text-emerald-300">₹{Number(item.final_amount ?? item.total_amount ?? 0).toFixed(2)}</p>
                    <button
                      type="button"
                      onClick={() => handleRemove(item.product_id, item.qty)}
                      disabled={removing === String(item.product_id)}
                      className="mt-2 rounded-lg bg-rose-600/80 px-3 py-1 text-sm font-medium text-white hover:bg-rose-600/95 disabled:opacity-60"
                    >
                      {removing === String(item.product_id) ? 'Removing...' : 'Remove'}
                    </button>
                  </div>
                </div>
              ))}

            </div>

            <aside className="rounded-xl bg-slate-900/80 p-6 shadow-sm">
              <h3 className="text-xl font-bold text-white">Summary</h3>
              <div className="mt-4 flex items-center justify-between text-slate-300">
                <span>Subtotal</span>
                <span>₹{subtotal.toFixed(2)}</span>
              </div>
              <button type="button" className="mt-6 w-full rounded-lg bg-gradient-to-r from-cyan-500 to-violet-500 px-4 py-3 font-medium text-white hover:brightness-105">
                Checkout
              </button>
            </aside>
          </div>
        )}
      </div>
    </div>
  );
}
