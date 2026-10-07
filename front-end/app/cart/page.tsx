'use client';

import axios from 'axios';
import Link from 'next/link';
import { useEffect, useState } from 'react';

interface CartItem {
  id?: number | string;
  cart_id?: number | string;
  user_id?: number | string;
  product_id?: number | string;
  product_name?: string;
  category_name?: string;
  product_description?: string;
  product_image_url?: string;
  popular?: boolean;
  size?: string;
  stock?: number;
  discount_percentage?: number;
  original_price?: number;
  offer_unit_price?: number;
  original_total_price?: number;
  offer_total_price?: number;
  total_amount?: number;
  qty?: number;
  price?: number;
  discount?: number;
  final_amount?: number;
  final_price?: number;
}

interface CartSummary {
  items_count: number;
  total_quantity: number;
  original_total: number;
  discount_amount: number;
  final_total: number;
}

export default function CartPage() {
  const [items, setItems] = useState<CartItem[]>([]);
  const [summary, setSummary] = useState<CartSummary>({
    items_count: 0,
    total_quantity: 0,
    original_total: 0,
    discount_amount: 0,
    final_total: 0,
  });
  const [loading, setLoading] = useState(true);
  const [removing, setRemoving] = useState<string | null>(null);
  const [quantityErrors, setQuantityErrors] = useState<Record<string, string>>({});

  const loadCart = async () => {
    setLoading(true);
    try {
      const [cartResponse, summaryResponse] = await Promise.all([
        axios.get('http://localhost:8080/cart', { withCredentials: true }),
        axios.get('http://localhost:8080/cart/summary', { withCredentials: true }),
      ]);

      const cartItems = cartResponse.data?.data?.cart_items ?? [];
      const summaryData = summaryResponse.data?.data ?? {};

      setItems(Array.isArray(cartItems) ? cartItems : []);
      setSummary({
        items_count: Number(summaryData.items_count ?? 0),
        total_quantity: Number(summaryData.total_quantity ?? 0),
        original_total: Number(summaryData.original_total ?? 0),
        discount_amount: Number(summaryData.discount_amount ?? 0),
        final_total: Number(summaryData.final_total ?? 0),
      });
    } catch (error) {
      console.error('Failed to load cart', error);
      setItems([]);
      setSummary({ items_count: 0, total_quantity: 0, original_total: 0, discount_amount: 0, final_total: 0 });
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    let isMounted = true;

    const fetchCart = async () => {
      setLoading(true);
      try {
        const [cartResponse, summaryResponse] = await Promise.all([
          axios.get('http://localhost:8080/cart', { withCredentials: true }),
          axios.get('http://localhost:8080/cart/summary', { withCredentials: true }),
        ]);

        const cartItems = cartResponse.data?.data?.cart_items ?? [];
        const summaryData = summaryResponse.data?.data ?? {};

        if (isMounted) {
          setItems(Array.isArray(cartItems) ? cartItems : []);
          setSummary({
            items_count: Number(summaryData.items_count ?? 0),
            total_quantity: Number(summaryData.total_quantity ?? 0),
            original_total: Number(summaryData.original_total ?? 0),
            discount_amount: Number(summaryData.discount_amount ?? 0),
            final_total: Number(summaryData.final_total ?? 0),
          });
        }
      } catch (error) {
        console.error('Failed to load cart', error);
        if (isMounted) {
          setItems([]);
          setSummary({ items_count: 0, total_quantity: 0, original_total: 0, discount_amount: 0, final_total: 0 });
        }
      } finally {
        if (isMounted) {
          setLoading(false);
        }
      }
    };

    void fetchCart();

    return () => {
      isMounted = false;
    };
  }, []);

  const subtotal = summary.final_total || items.reduce((sum, item) => sum + Number(item.final_amount ?? item.total_amount ?? 0), 0);

  const handleRemove = async (productId?: number | string, qty?: number, size?: string) => {
    if (!productId) return;
    const should = window.confirm('Remove this item from cart?');
    if (!should) return;
    try {
      setRemoving(String(productId));
      await axios.delete('http://localhost:8080/cart', {
        data: { product_id: productId, quantity: qty ?? 0, size: size ?? '' },
        withCredentials: true,
      });
      await loadCart();
    } catch (err) {
      console.error('Failed to remove item', err);
    } finally {
      setRemoving(null);
    }
  };

  const updateQuantity = async (item: CartItem, nextQty: number) => {
    const productId = item.product_id;
    const size = item.size ?? '';

    if (!productId) return;

    const stock = Number(item.stock ?? 0);
    if (Number.isNaN(stock) || stock <= 0) {
      setQuantityErrors((prev) => ({ ...prev, [String(productId) + '-' + size]: 'This product is currently out of stock.' }));
      return;
    }

    if (nextQty < 1) {
      setQuantityErrors((prev) => ({ ...prev, [String(productId) + '-' + size]: 'Quantity must be at least 1.' }));
      return;
    }

    if (nextQty > 5) {
      setQuantityErrors((prev) => ({ ...prev, [String(productId) + '-' + size]: 'Maximum quantity allowed for this product is 5.' }));
      return;
    }

    if (nextQty > stock) {
      setQuantityErrors((prev) => ({ ...prev, [String(productId) + '-' + size]: `Only ${stock} item(s) available in stock for this variant. Quantity cannot exceed this limit.` }));
      return;
    }

    try {
      setQuantityErrors((prev) => ({ ...prev, [String(productId) + '-' + size]: '' }));
      await axios.put(
        'http://localhost:8080/cart/update-quantity',
        { product_id: productId, size, quantity: nextQty },
        { withCredentials: true },
      );
      await loadCart();
    } catch (err: any) {
      const message = err?.response?.data?.message || 'Unable to update quantity right now.';
      setQuantityErrors((prev) => ({ ...prev, [String(productId) + '-' + size]: message }));
      console.error('Failed to update quantity', err);
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
              {items.map((item, index) => {
                const itemKey = item.id != null
                  ? `${item.id}-${item.size ?? 'unknown'}`
                  : `${item.product_id ?? 'product'}-${item.size ?? 'size'}-${index}`;

                const productImage = item.product_image_url || 'https://images.unsplash.com/photo-1521572267360-ee0c2909d518?auto=format&fit=crop&w=900&q=80';
                const offerUnitPrice = Number(item.offer_unit_price ?? item.price ?? 0);
                const originalTotal = Number(item.original_total_price ?? item.total_amount ?? 0);
                const offerTotal = Number(item.offer_total_price ?? item.final_price ?? item.final_amount ?? 0);
                const itemQty = Number(item.qty ?? 1);
                const originalPrice = Number(item.original_price ?? item.price ?? 0);
                const hasOffer = Number(item.discount_percentage ?? 0) > 0;

                return (
                  <div key={itemKey} className="overflow-hidden rounded-2xl border border-white/10 bg-slate-900/80 shadow-lg transition hover:border-cyan-500/30">
                    <div className="flex flex-col sm:flex-row">
                      <div className="h-40 w-full sm:w-40 flex-shrink-0 overflow-hidden bg-slate-800">
                        <img src={productImage} alt={item.product_name ?? 'Product'} className="h-full w-full object-cover" />
                      </div>

                      <div className="flex flex-1 flex-col justify-between p-4">
                        <div className="flex items-start justify-between gap-3">
                          <div>
                            <div className="mb-2 flex flex-wrap items-center gap-2">
                              {item.category_name && (
                                <span className="rounded-full border border-cyan-400/30 bg-cyan-500/10 px-2 py-1 text-[10px] uppercase tracking-[0.2em] text-cyan-200">
                                  {item.category_name}
                                </span>
                              )}
                              {item.popular && (
                                <span className="rounded-full border border-amber-400/30 bg-amber-500/10 px-2 py-1 text-[10px] uppercase tracking-[0.2em] text-amber-200">
                                  Popular
                                </span>
                              )}
                            </div>
                            <h2 className="text-xl font-semibold text-white">{item.product_name ?? 'Product'}</h2>
                            {item.product_description && (
                              <p className="mt-1 line-clamp-2 text-sm text-slate-400">{item.product_description}</p>
                            )}
                          </div>

                          <p className="text-lg font-bold text-emerald-300">₹{offerTotal.toFixed(2)}</p>
                        </div>

                        <div className="mt-4 flex flex-wrap items-center gap-2 text-xs text-slate-300">
                          {item.size && (
                            <span className="rounded-full border border-white/10 bg-white/5 px-2 py-1">Size: {item.size}</span>
                          )}
                          <span className="rounded-full border border-white/10 bg-white/5 px-2 py-1">Qty: {itemQty}</span>
                          <span className="rounded-full border border-amber-500/30 bg-amber-500/10 px-2 py-1 text-amber-200">Original total: ₹{originalTotal.toFixed(2)}</span>
                          <span className="rounded-full border border-emerald-500/30 bg-emerald-500/10 px-2 py-1 text-emerald-200">Offer total: ₹{offerTotal.toFixed(2)}</span>
                          {typeof item.stock === 'number' && (
                            <span className="rounded-full border border-white/10 bg-white/5 px-2 py-1">Stock: {item.stock}</span>
                          )}
                        </div>

                        <div className="mt-3 space-y-1 text-xs text-slate-300">
                          <p className="text-slate-400 line-through">Original unit: ₹{originalPrice.toFixed(2)}</p>
                          {hasOffer ? (
                            <>
                              <p className="text-amber-300">Offer: {Number(item.discount_percentage ?? 0)}%</p>
                              <p className="text-emerald-300">Offer unit: ₹{offerUnitPrice.toFixed(2)}</p>
                              <p className="text-rose-400">Savings: ₹{(originalTotal - offerTotal).toFixed(2)}</p>
                            </>
                          ) : (
                            <p className="text-emerald-300">Offer unit: ₹{offerUnitPrice.toFixed(2)}</p>
                          )}
                        </div>

                        {!hasOffer && Number(item.discount ?? 0) > 0 && (
                          <p className="mt-3 text-xs text-rose-400">You saved: ₹{Number(item.discount).toFixed(2)}</p>
                        )}

                        <div className="mt-4 flex items-center justify-between gap-3">
                          <div className="text-sm text-slate-400">
                            {hasOffer ? 'Discount applied' : 'No offer active'}
                          </div>

                          <div className="flex items-center gap-2 rounded-lg border border-white/10 bg-slate-950/60 px-2 py-1">
                            <button
                              type="button"
                              onClick={() => updateQuantity(item, Number(item.qty ?? 1) - 1)}
                              className="h-8 w-8 rounded-md bg-slate-700 text-lg font-semibold text-white hover:bg-slate-600"
                              aria-label="Decrease quantity"
                            >
                              −
                            </button>
                            <span className="min-w-8 text-center text-sm font-semibold text-white">{itemQty}</span>
                            <button
                              type="button"
                              onClick={() => updateQuantity(item, Number(item.qty ?? 1) + 1)}
                              className="h-8 w-8 rounded-md bg-slate-700 text-lg font-semibold text-white hover:bg-slate-600"
                              aria-label="Increase quantity"
                            >
                              +
                            </button>
                          </div>

                          <button
                            type="button"
                            onClick={() => handleRemove(item.product_id, item.qty, item.size)}
                            disabled={removing === String(item.product_id)}
                            className="rounded-lg bg-rose-600/80 px-3 py-2 text-sm font-medium text-white hover:bg-rose-600/95 disabled:cursor-not-allowed disabled:opacity-60"
                          >
                            {removing === String(item.product_id) ? 'Removing...' : 'Remove'}
                          </button>
                        </div>

                        {quantityErrors[String(item.product_id ?? 'product') + '-' + (item.size ?? '')] && (
                          <p className="mt-3 text-xs text-rose-300">
                            {quantityErrors[String(item.product_id ?? 'product') + '-' + (item.size ?? '')]}
                          </p>
                        )}
                      </div>
                    </div>
                  </div>
                );
              })}
            </div>

            <aside className="rounded-xl bg-slate-900/80 p-6 shadow-sm">
              <h3 className="text-xl font-bold text-white">Summary</h3>

              <div className="mt-4 space-y-3 text-sm text-slate-300">
                <div className="flex items-center justify-between">
                  <span>Items</span>
                  <span>{summary.items_count}</span>
                </div>
                <div className="flex items-center justify-between">
                  <span>Total quantity</span>
                  <span>{summary.total_quantity}</span>
                </div>
                <div className="flex items-center justify-between">
                  <span>Original total</span>
                  <span>₹{summary.original_total.toFixed(2)}</span>
                </div>
                <div className="flex items-center justify-between text-rose-300">
                  <span>Discount</span>
                  <span>-₹{summary.discount_amount.toFixed(2)}</span>
                </div>
                <div className="flex items-center justify-between border-t border-white/10 pt-3 text-base font-semibold text-white">
                  <span>Final total</span>
                  <span>₹{subtotal.toFixed(2)}</span>
                </div>
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
