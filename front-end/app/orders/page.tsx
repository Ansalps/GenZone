'use client';

import axios from 'axios';
import Link from 'next/link';
import { useEffect, useState } from 'react';

const API = 'http://localhost:8080';

type Address = {
  id?: number;
  city?: string;
  state?: string;
  country?: string;
  street_name?: string;
  pin_code?: string;
  phone?: string;
};

type OrderItem = {
  id?: number;
  product_id?: number;
  product_name?: string;
  price?: number;
  order_staus?: string;
  payment_method?: string;
  coupon_discount?: number;
  offer_discount?: number;
  total_discount?: number;
  paid_amount?: number;
};

type Order = {
  id: number;
  created_at?: string;
  total_amount?: number;
  offer_discount?: number;
  coupon_discount?: number;
  total_discount_amount?: number;
  final_amount?: number;
  payment_method?: string;
  order_status?: string;
  address?: Address;
  items?: OrderItem[];
};

export default function OrdersPage() {
  const [orders, setOrders] = useState<Order[]>([]);
  const [orderItems, setOrderItems] = useState<Record<number, OrderItem[]>>({});
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    const fetchOrderItems = async (orderId: number) => {
      try {
        const response = await axios.get(`${API}/order-items/${orderId}`, { withCredentials: true });
        const items = Array.isArray(response.data?.['order items']) ? response.data['order items'] : [];
        setOrderItems((prev) => ({ ...prev, [orderId]: items }));
      } catch (err) {
        console.error(`Failed to load items for order ${orderId}`, err);
      }
    };

    const fetchOrders = async () => {
      try {
        setLoading(true);
        const response = await axios.get(`${API}/orders`, { withCredentials: true });
        const orderList = Array.isArray(response.data?.data?.Order) ? response.data.data.Order : [];
        setOrders(orderList);

        orderList.forEach((order: Order) => {
          fetchOrderItems(order.id);
        });
      } catch (err: any) {
        console.error('Failed to load orders', err);
        setError(err.response?.data?.message || 'Unable to load orders right now.');
        setOrders([]);
      } finally {
        setLoading(false);
      }
    };

    fetchOrders();
  }, []);

  return (
    <div className="min-h-screen bg-slate-950 px-6 py-10 text-white">
      <div className="mx-auto max-w-5xl">
        <div className="mb-6 flex items-center justify-between gap-4">
          <div>
            <p className="text-xs font-semibold uppercase tracking-[0.28em] text-cyan-300">Your orders</p>
            <h1 className="mt-2 text-3xl font-bold text-white">Order history</h1>
          </div>
          <Link
            href="/dashboard"
            className="rounded-xl border border-white/10 bg-white/5 px-4 py-2 text-sm font-medium text-slate-100 transition hover:border-cyan-400/40 hover:bg-cyan-500/10"
          >
            Back to dashboard
          </Link>
        </div>

        {error && (
          <div className="mb-6 rounded-2xl border border-red-500/40 bg-red-500/10 px-4 py-3 text-sm text-red-200">
            {error}
          </div>
        )}

        {loading ? (
          <div className="rounded-2xl border border-white/10 bg-slate-900/80 p-6 text-slate-300">Loading orders...</div>
        ) : orders.length === 0 ? (
          <div className="rounded-2xl border border-dashed border-white/10 bg-slate-900/80 p-10 text-center text-slate-300">
            No orders yet. Start shopping from the dashboard.
          </div>
        ) : (
          <div className="space-y-4">
            {orders.map((order) => (
              <div key={order.id} className="rounded-2xl border border-white/10 bg-slate-900/80 p-5 shadow-lg shadow-cyan-900/10">
                <div className="mb-4 flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
                  <div>
                    <p className="text-sm text-slate-400">Order #{order.id}</p>
                    <p className="text-lg font-semibold text-white">{new Date(order.created_at ?? Date.now()).toLocaleDateString()}</p>
                  </div>
                  <div className="flex flex-wrap items-center gap-2 text-sm">
                    <span className="rounded-full border border-cyan-500/40 bg-cyan-500/10 px-2.5 py-1 text-cyan-200">
                      {order.payment_method || 'COD'}
                    </span>
                    <span className="rounded-full border border-violet-500/40 bg-violet-500/10 px-2.5 py-1 text-violet-200">
                      {order.order_status || 'pending'}
                    </span>
                  </div>
                </div>

                <div className="grid gap-3 text-sm text-slate-300 md:grid-cols-4">
                  <div>
                    <p className="text-xs uppercase tracking-[0.2em] text-slate-500">Total</p>
                    <p className="mt-1 text-base font-semibold text-white">₹{Number(order.total_amount ?? 0).toFixed(2)}</p>
                  </div>
                  <div>
                    <p className="text-xs uppercase tracking-[0.2em] text-slate-500">Offer disc.</p>
                    <p className="mt-1 text-base font-semibold text-amber-300">-₹{Number(order.offer_discount ?? 0).toFixed(2)}</p>
                  </div>
                  <div>
                    <p className="text-xs uppercase tracking-[0.2em] text-slate-500">Coupon disc.</p>
                    <p className="mt-1 text-base font-semibold text-purple-300">-₹{Number(order.coupon_discount ?? 0).toFixed(2)}</p>
                  </div>
                  <div>
                    <p className="text-xs uppercase tracking-[0.2em] text-slate-500">Total disc.</p>
                    <p className="mt-1 text-base font-semibold text-emerald-300">-₹{Number(order.total_discount_amount ?? 0).toFixed(2)}</p>
                  </div>
                </div>

                <div className="mt-3 grid gap-3 text-sm text-slate-300 md:grid-cols-2">
                  <div className="rounded-lg border border-white/10 bg-slate-950/50 p-3">
                    <p className="text-xs uppercase tracking-[0.2em] text-slate-500">Final amount</p>
                    <p className="mt-1 text-base font-semibold text-white">₹{Number(order.final_amount ?? order.total_amount ?? 0).toFixed(2)}</p>
                  </div>
                  <div className="rounded-lg border border-white/10 bg-slate-950/50 p-3">
                    <p className="text-xs uppercase tracking-[0.2em] text-slate-500">Delivery</p>
                    <p className="mt-1 text-base font-semibold text-white">
                      {order.address ? `${order.address.city || ''}${order.address.city && order.address.state ? ', ' : ''}${order.address.state || ''}` : 'Address unavailable'}
                    </p>
                  </div>
                </div>

                {(orderItems[order.id]?.length ?? 0) > 0 && (
                  <div className="mt-5 rounded-xl border border-white/10 bg-slate-950/60 p-4">
                    <div className="mb-3 flex items-center justify-between">
                      <p className="text-sm font-semibold uppercase tracking-[0.2em] text-cyan-300">Order items</p>
                      <span className="text-xs text-slate-400">{orderItems[order.id].length} item(s)</span>
                    </div>

                    <div className="space-y-3">
                      {orderItems[order.id].map((item) => (
                        <div key={item.id ?? `${order.id}-${item.product_id}`} className="rounded-lg border border-white/10 bg-slate-900/70 p-3">
                          <div className="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
                            <div>
                              <p className="font-medium text-white">{item.product_name || 'Product'}</p>
                              <p className="text-xs text-slate-400">{item.payment_method || 'Payment pending'}</p>
                            </div>

                            <div className="flex flex-wrap items-center gap-3 text-sm text-slate-300">
                              <span>₹{Number(item.price ?? 0).toFixed(2)}</span>
                              <span className="text-cyan-300">Paid: ₹{Number(item.paid_amount ?? item.price ?? 0).toFixed(2)}</span>
                            </div>
                          </div>

                          <div className="mt-3 grid gap-2 text-xs text-slate-300 sm:grid-cols-3">
                            <div className="rounded border border-amber-500/30 bg-amber-500/10 px-2 py-1.5">
                              <span className="block text-[10px] uppercase tracking-[0.2em] text-amber-300">Offer disc.</span>
                              <span className="font-semibold text-amber-200">-₹{Number(item.offer_discount ?? 0).toFixed(2)}</span>
                            </div>
                            <div className="rounded border border-purple-500/30 bg-purple-500/10 px-2 py-1.5">
                              <span className="block text-[10px] uppercase tracking-[0.2em] text-purple-300">Coupon disc.</span>
                              <span className="font-semibold text-purple-200">-₹{Number(item.coupon_discount ?? 0).toFixed(2)}</span>
                            </div>
                            <div className="rounded border border-emerald-500/30 bg-emerald-500/10 px-2 py-1.5">
                              <span className="block text-[10px] uppercase tracking-[0.2em] text-emerald-300">Total disc.</span>
                              <span className="font-semibold text-emerald-200">-₹{Number(item.total_discount ?? 0).toFixed(2)}</span>
                            </div>
                          </div>
                        </div>
                      ))}
                    </div>
                  </div>
                )}
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
