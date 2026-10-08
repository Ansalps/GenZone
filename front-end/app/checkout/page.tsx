'use client';

import axios from 'axios';
import Link from 'next/link';
import { useCallback, useEffect, useState } from 'react';

const API = 'http://localhost:8080';

interface CartItem {
  product_id?: number | string;
  product_name?: string;
  product_description?: string;
  product_image_url?: string;
  category_name?: string;
  qty?: number;
  price?: number;
  total_amount?: number;
  final_amount?: number;
  offer_unit_price?: number;
  original_total_price?: number;
  offer_total_price?: number;
  discount?: number;
}

interface AddressItem {
  id?: number;
  country?: string;
  state?: string;
  city?: string;
  district?: string;
  street_name?: string;
  pin_code?: string;
  phone?: string;
  default?: boolean;
  Default?: boolean;
}

interface CartSummary {
  items_count?: number;
  total_quantity?: number;
  original_total?: number;
  discount_amount?: number;
  final_total?: number;
}

interface CheckoutResponse {
  cart_items?: CartItem[];
  total_amount?: number;
  offer_applied?: number;
  coupon_discount?: number;
  coupon_applied_amount?: number;
  address?: AddressItem[];
}

export default function CheckoutPage() {
  const [checkout, setCheckout] = useState<CheckoutResponse>({});
  const [cartSummary, setCartSummary] = useState<CartSummary>({});
  const [cartItems, setCartItems] = useState<CartItem[]>([]);
  const [addresses, setAddresses] = useState<AddressItem[]>([]);
  const [selectedAddressId, setSelectedAddressId] = useState<number | null>(null);
  const [couponCode, setCouponCode] = useState('');
  const [placingOrder, setPlacingOrder] = useState(false);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  const loadCheckout = useCallback(async (value = couponCode) => {
    try {
      setLoading(true);
      setError('');

      const [cartResponse, cartSummaryResponse, checkoutResponse] = await Promise.all([
        axios.get(`${API}/cart`, { withCredentials: true }),
        axios.get(`${API}/cart/summary`, { withCredentials: true }),
        axios.get(`${API}/checkout`, {
          withCredentials: true,
          params: value ? { coupon_code: value } : {},
        }),
      ]);

      const cartPayload = Array.isArray(cartResponse.data?.data?.cart_items)
        ? cartResponse.data.data.cart_items
        : [];
      const nextCartItems = Array.isArray(cartPayload) ? cartPayload : [];
      const nextCartSummary = (cartSummaryResponse.data?.data ?? {}) as CartSummary;

      setCartItems(nextCartItems);
      setCartSummary(nextCartSummary);

      const result = checkoutResponse.data?.result ?? checkoutResponse.data?.data ?? checkoutResponse.data ?? {};
      const nextCheckout = (result as CheckoutResponse) ?? {};
      const nextAddresses = Array.isArray((result as { address?: AddressItem[] })?.address)
        ? (result as { address?: AddressItem[] }).address ?? []
        : Array.isArray(nextCheckout.address)
          ? nextCheckout.address
          : [];

      setCheckout(nextCheckout);
      setAddresses(nextAddresses);

      const defaultAddress = nextAddresses.find((address) => Boolean(address.default ?? address.Default));
      setSelectedAddressId(defaultAddress?.id ?? nextAddresses[0]?.id ?? null);
    } catch (err: unknown) {
      const message = (err as { response?: { data?: { message?: string } } })?.response?.data?.message ?? 'Unable to load checkout summary.';
      console.error('Failed to load checkout', err);
      setError(message);
      setCheckout({});
      setCartSummary({});
      setCartItems([]);
      setAddresses([]);
      setSelectedAddressId(null);
    } finally {
      setLoading(false);
    }
  }, [couponCode]);

  useEffect(() => {
    let active = true;

    const init = async () => {
      await loadCheckout();
      if (!active) return;
    };

    void init();

    return () => {
      active = false;
    };
  }, [loadCheckout]);

  const subtotal = Number(cartSummary.original_total ?? 0);
  const discount = Number(checkout.coupon_discount ?? 0);
  const offerDiscount = Number(cartSummary.discount_amount ?? 0);
  const finalTotal = Math.max(subtotal - offerDiscount - discount, 0);

  const handleApplyCoupon = async () => {
    await loadCheckout(couponCode.trim());
  };

  const handlePlaceOrder = async () => {
    if (!selectedAddressId) {
      setError('Please choose a delivery address before placing the order.');
      return;
    }

    try {
      setPlacingOrder(true);
      setError('');

      await axios.post(
        `${API}/checkout/order`,
        {
          address_id: selectedAddressId,
          coupon_code: couponCode.trim(),
        },
        { withCredentials: true }
      );

      window.alert('Order placed successfully!');
      window.location.href = '/dashboard';
    } catch (err: unknown) {
      const message = (err as { response?: { data?: { message?: string } } })?.response?.data?.message ?? 'Unable to place the order right now.';
      console.error('Failed to place order', err);
      setError(message);
    } finally {
      setPlacingOrder(false);
    }
  };

  return (
    <div className="min-h-screen bg-slate-950 text-white">
      <div className="mx-auto max-w-6xl px-6 py-10">
        <div className="mb-6 flex items-center justify-between gap-4">
          <div>
            <p className="text-xs font-semibold uppercase tracking-[0.25em] text-cyan-300">Checkout</p>
            <h1 className="mt-2 text-3xl font-bold">Complete your order</h1>
          </div>
          <Link href="/cart" className="rounded-xl border border-white/10 bg-white/5 px-4 py-2 text-sm font-medium text-slate-100 transition hover:border-cyan-400/40 hover:bg-cyan-500/10">
            Back to cart
          </Link>
        </div>

        {error && (
          <div className="mb-6 rounded-xl border border-rose-500/40 bg-rose-500/10 px-4 py-3 text-sm text-rose-200">
            {error}
          </div>
        )}

        {loading ? (
          <div className="rounded-2xl border border-white/10 bg-slate-900/80 p-10 text-center text-slate-300">
            Loading checkout details...
          </div>
        ) : (
          <div className="grid gap-6 lg:grid-cols-[1.5fr_0.9fr]">
            <div className="space-y-6">
              <section className="rounded-2xl border border-white/10 bg-slate-900/80 p-5 shadow-sm">
                <div className="mb-4 flex items-center justify-between">
                  <h2 className="text-xl font-semibold">Delivery address</h2>
                  <Link href="/profile" className="text-sm text-cyan-300 hover:text-cyan-200">
                    Manage addresses
                  </Link>
                </div>

                {addresses.length === 0 ? (
                  <div className="rounded-xl border border-dashed border-slate-700 bg-slate-950/60 p-5 text-slate-300">
                    No saved addresses yet. Add one in your profile before checkout.
                  </div>
                ) : (
                  <div className="space-y-3">
                    {addresses.map((address) => {
                      const isSelected = Number(address.id) === Number(selectedAddressId);
                      return (
                        <button
                          key={address.id ?? `${address.street_name}-${address.phone}`}
                          type="button"
                          onClick={() => setSelectedAddressId(Number(address.id))}
                          className={`w-full rounded-xl border p-4 text-left transition ${
                            isSelected
                              ? 'border-cyan-400/60 bg-cyan-500/10'
                              : 'border-white/10 bg-slate-950/60 hover:border-cyan-400/40'
                          }`}
                        >
                          <div className="flex items-start justify-between gap-4">
                            <div>
                              <div className="mb-2 flex items-center gap-2">
                                <p className="font-semibold text-white">{address.street_name || 'Address'}</p>
                                {Boolean(address.default ?? address.Default) && (
                                  <span className="rounded-full border border-emerald-500/30 bg-emerald-500/10 px-2 py-0.5 text-[10px] uppercase tracking-[0.2em] text-emerald-200">
                                    Default
                                  </span>
                                )}
                              </div>
                              <p className="text-sm text-slate-300">
                                {address.street_name}, {address.city ?? address.district}, {address.state}, {address.country}
                              </p>
                              <p className="mt-1 text-sm text-slate-300">Pin: {address.pin_code}</p>
                              <p className="mt-1 text-sm text-slate-300">Phone: {address.phone}</p>
                            </div>
                            {isSelected && (
                              <span className="rounded-full bg-cyan-500 px-2 py-1 text-[10px] font-semibold uppercase tracking-[0.2em] text-slate-950">
                                Selected
                              </span>
                            )}
                          </div>
                        </button>
                      );
                    })}
                  </div>
                )}
              </section>

              <section className="rounded-2xl border border-white/10 bg-slate-900/80 p-5 shadow-sm">
                <h2 className="mb-4 text-xl font-semibold">Checkout Items</h2>

                {cartItems.length === 0 ? (
                  <div className="rounded-xl border border-dashed border-slate-700 bg-slate-950/60 p-5 text-slate-300">
                    Your cart is empty.
                  </div>
                ) : (
                  <div className="space-y-4">
                    {cartItems.map((item, index) => (
                      <div key={`${item.product_id ?? 'product'}-${index}`} className="flex items-center gap-4 rounded-xl border border-white/10 bg-slate-950/60 p-3">
                        <div className="h-16 w-16 overflow-hidden rounded-lg border border-white/10 bg-slate-800">
                          <img
                            src={item.product_image_url || 'https://images.unsplash.com/photo-1521572267360-ee0c2909d518?auto=format&fit=crop&w=900&q=80'}
                            alt={item.product_name || 'Product'}
                            className="h-full w-full object-cover"
                          />
                        </div>
                        <div className="flex-1">
                          <p className="font-medium text-white">{item.product_name || 'Product'}</p>
                          <p className="mt-1 text-sm text-slate-400">Qty: {item.qty ?? 1}</p>
                        </div>
                        <div className="text-right">
                          <p className="font-semibold text-emerald-300">₹{Number(item.final_amount ?? item.total_amount ?? 0).toFixed(2)}</p>
                        </div>
                      </div>
                    ))}
                  </div>
                )}
              </section>
            </div>

            <aside className="rounded-2xl border border-white/10 bg-slate-900/80 p-5 shadow-sm">
              <h2 className="text-xl font-semibold">Summary</h2>

              <div className="mt-4 space-y-3 text-sm text-slate-300">
                <div className="flex items-center justify-between">
                  <span>Original total</span>
                  <span>₹{subtotal.toFixed(2)}</span>
                </div>
                <div className="flex items-center justify-between text-amber-300">
                  <span>Offer Discount</span>
                  <span>-₹{offerDiscount.toFixed(2)}</span>
                </div>
                <div className="flex items-center justify-between text-rose-300">
                  <span>Coupon Discount</span>
                  <span>-₹{discount.toFixed(2)}</span>
                </div>
                <div className="flex items-center justify-between border-t border-white/10 pt-3 text-base font-semibold text-white">
                  <span>Final total</span>
                  <span>₹{finalTotal.toFixed(2)}</span>
                </div>
              </div>

              <div className="mt-5 space-y-3">
                <label className="block text-sm text-slate-300">Coupon code</label>
                <div className="flex gap-2">
                  <input
                    type="text"
                    value={couponCode}
                    onChange={(e) => setCouponCode(e.target.value)}
                    placeholder="Enter coupon code"
                    className="w-full rounded-xl border border-white/10 bg-slate-950 px-3 py-2 text-white placeholder:text-slate-500 outline-none focus:border-cyan-400"
                  />
                  <button
                    type="button"
                    onClick={handleApplyCoupon}
                    className="rounded-xl border border-cyan-400/40 bg-cyan-500/10 px-3 py-2 text-sm font-medium text-cyan-200 hover:bg-cyan-500/20"
                  >
                    Apply
                  </button>
                </div>
              </div>

              <button
                type="button"
                onClick={handlePlaceOrder}
                disabled={placingOrder || cartItems.length === 0 || !selectedAddressId}
                className="mt-6 w-full rounded-xl bg-gradient-to-r from-cyan-500 to-violet-500 px-4 py-3 font-semibold text-white shadow-lg shadow-cyan-500/20 transition hover:brightness-110 disabled:cursor-not-allowed disabled:opacity-60"
              >
                {placingOrder ? 'Placing order...' : 'Place order'}
              </button>
            </aside>
          </div>
        )}
      </div>
    </div>
  );
}
