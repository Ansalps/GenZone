'use client';

import axios from 'axios';
import Link from 'next/link';
import { useEffect, useState } from 'react';

interface CartItem {
  user_id?: number;
  product_id?: number;
  product_name?: string;
  total_amount?: number;
  qty?: number;
  price?: number;
  discount?: number;
  final_amount?: number;
}

export default function CartPage() {
  const [items, setItems] = useState<CartItem[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const loadCart = async () => {
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

    loadCart();
  }, []);

  const subtotal = items.reduce((sum, item) => sum + Number(item.final_amount ?? item.total_amount ?? 0), 0);

  return (
    <div className="min-h-screen bg-slate-100 p-6">
      <div className="mx-auto max-w-5xl">
        <div className="mb-6 flex items-center justify-between">
          <h1 className="text-3xl font-bold text-slate-800">Your Cart</h1>
          <Link href="/dashboard" className="rounded-lg bg-green-700 px-4 py-2 text-white hover:bg-green-800">
            Continue Shopping
          </Link>
        </div>

        {loading ? (
          <div className="rounded-xl bg-white p-8 text-slate-600">Loading cart...</div>
        ) : items.length === 0 ? (
          <div className="rounded-xl bg-white p-12 text-center text-slate-600 shadow-sm">
            <p className="text-xl font-semibold">Your cart is empty.</p>
            <Link href="/dashboard" className="mt-4 inline-block rounded-lg bg-green-700 px-4 py-2 text-white">
              Browse Products
            </Link>
          </div>
        ) : (
          <div className="grid gap-6 lg:grid-cols-[2fr_1fr]">
            <div className="space-y-4">
              {items.map((item, index) => (
                <div key={`${item.product_id ?? index}`} className="flex items-center justify-between rounded-xl bg-white p-4 shadow-sm">
                  <div>
                    <h2 className="text-lg font-semibold text-slate-800">{item.product_name ?? 'Product'}</h2>
                    <p className="text-sm text-slate-500">Qty: {item.qty ?? 1}</p>
                  </div>
                  <div className="text-right">
                    <p className="font-semibold text-green-700">₹{Number(item.final_amount ?? item.total_amount ?? 0).toFixed(2)}</p>
                    {Number(item.discount ?? 0) > 0 && (
                      <p className="text-xs text-slate-500">Discount: ₹{Number(item.discount).toFixed(2)}</p>
                    )}
                  </div>
                </div>
              ))}
            </div>

            <aside className="rounded-xl bg-white p-6 shadow-sm">
              <h3 className="text-xl font-bold text-slate-800">Summary</h3>
              <div className="mt-4 flex items-center justify-between text-slate-600">
                <span>Subtotal</span>
                <span>₹{subtotal.toFixed(2)}</span>
              </div>
              <button type="button" className="mt-6 w-full rounded-lg bg-black px-4 py-3 font-medium text-white hover:bg-slate-800">
                Checkout
              </button>
            </aside>
          </div>
        )}
      </div>
    </div>
  );
}
