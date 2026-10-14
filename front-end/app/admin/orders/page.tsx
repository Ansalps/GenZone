'use client';

import axios from 'axios';
import Link from 'next/link';
import { useEffect, useState } from 'react';
import { toast } from 'sonner';

const API = 'http://localhost:8080';

type UserRecord = {
  id: number;
  first_name?: string;
  last_name?: string;
  email?: string;
  phone?: string;
  status?: string;
  profile_picture?: string | null;
};

type OrderItem = {
  id: number;
  order_id?: number | string;
  product_name?: string;
  size?: string;
  quantity?: number;
  price?: number;
  total_amount?: number;
  order_status?: string;
  delivery_status?: string;
  payment_method?: string;
  coupon_discount?: number;
  offer_discount?: number;
  total_discount?: number;
  paid_amount?: number;
};

type OrderRecord = {
  id: number;
  user_id?: number | string;
  address_id?: number;
  total_amount?: number;
  offer_discount?: number;
  coupon_discount?: number;
  offer_applied?: number;
  total_discount_amount?: number;
  final_amount?: number;
  payment_method?: string;
  payment_status?: string;
  order_status?: string;
  delivery_status?: string;
  coupon_code?: string;
  created_at?: string;
  updated_at?: string;
  address?: {
    id?: number;
    phone?: string;
    country?: string;
    state?: string;
    city?: string;
    district?: string;
    street_name?: string;
    pin_code?: string;
    default?: boolean;
  };
};

const formatMoney = (amount: number = 0) =>
  new Intl.NumberFormat('en-IN', {
    style: 'currency',
    currency: 'INR',
    maximumFractionDigits: 2,
  }).format(amount);

const formatDate = (value?: string) => {
  if (!value) return '—';

  try {
    return new Date(value).toLocaleString('en-IN', {
      dateStyle: 'medium',
      timeStyle: 'short',
    });
  } catch {
    return value;
  }
};

const getDisplayName = (user?: UserRecord) => {
  if (!user) return 'Unknown User';

  const first = user.first_name?.trim();
  const last = user.last_name?.trim();

  if (first || last) {
    return `${first || ''} ${last || ''}`.trim();
  }

  return user.email || 'Unknown User';
};

const getInitials = (user?: UserRecord) => {
  const name = getDisplayName(user);
  return name
    .split(' ')
    .filter(Boolean)
    .slice(0, 2)
    .map((part) => part[0]?.toUpperCase() || '')
    .join('') || 'U';
};

const getOrderStatusClasses = (status?: string) => {
  const normalized = (status || 'pending').toLowerCase();

  if (normalized === 'delivered') {
    return 'border-emerald-400/30 bg-emerald-500/10 text-emerald-200';
  }

  if (normalized === 'cancelled' || normalized === 'failed') {
    return 'border-red-400/30 bg-red-500/10 text-red-200';
  }

  if (normalized === 'shipped') {
    return 'border-amber-400/30 bg-amber-500/10 text-amber-200';
  }

  return 'border-sky-400/30 bg-sky-500/10 text-sky-200';
};

export default function AdminOrdersPage() {
  const [orders, setOrders] = useState<OrderRecord[]>([]);
  const [users, setUsers] = useState<UserRecord[]>([]);
  const [expandedIds, setExpandedIds] = useState<number[]>([]);
  const [itemsByOrder, setItemsByOrder] = useState<Record<number, OrderItem[]>>({});
  const [loading, setLoading] = useState(false);
  const [loadingItemsId, setLoadingItemsId] = useState<number | null>(null);

  const loadUsers = async () => {
    try {
      const response = await axios.get(`${API}/admin/listusers`, { withCredentials: true });
      const list = response.data?.data?.users ?? [];
      setUsers(Array.isArray(list) ? list : []);
    } catch (error) {
      console.error('Failed to load users for order listing:', error);
      setUsers([]);
    }
  };

  const loadOrders = async () => {
    try {
      const response = await axios.get(`${API}/admin/orderlist`, { withCredentials: true });
      const list = response.data?.Order ?? response.data?.orders ?? [];
      setOrders(Array.isArray(list) ? list : []);
    } catch (error) {
      console.error('Failed to load orders:', error);
      setOrders([]);
      toast.error('Failed to load orders');
    }
  };

  useEffect(() => {
    const init = async () => {
      setLoading(true);
      await Promise.all([loadUsers(), loadOrders()]);
      setLoading(false);
    };

    init();
  }, []);

  const loadOrderItems = async (orderId: number) => {
    try {
      setLoadingItemsId(orderId);
      const response = await axios.get(`${API}/admin/orderlist/items/${orderId}`, { withCredentials: true });
      const items = response.data?.order_items ?? [];

      setItemsByOrder((prev) => ({
        ...prev,
        [orderId]: Array.isArray(items) ? items : [],
      }));

      setExpandedIds((prev) => (prev.includes(orderId) ? prev : [...prev, orderId]));
    } catch (error) {
      console.error(`Failed to load order items for ${orderId}:`, error);
      toast.error('Unable to load order items');
    } finally {
      setLoadingItemsId(null);
    }
  };

  const toggleOrder = async (order: OrderRecord) => {
    const isExpanded = expandedIds.includes(order.id);

    if (isExpanded) {
      setExpandedIds((prev) => prev.filter((itemId) => itemId !== order.id));
      return;
    }

    if (!itemsByOrder[order.id]) {
      await loadOrderItems(order.id);
      return;
    }

    setExpandedIds((prev) => (prev.includes(order.id) ? prev : [...prev, order.id]));
  };

  const totalRevenue = orders.reduce((sum, order) => sum + Number(order.final_amount ?? order.total_amount ?? 0), 0);

  return (
    <div className="min-h-screen bg-slate-950 px-4 py-8 text-white sm:px-6 lg:px-8">
      <div className="mx-auto max-w-7xl">
        <header className="mb-6 flex flex-col gap-4 rounded-3xl border border-white/10 bg-white/5 p-5 shadow-2xl shadow-slate-950/30 backdrop-blur-xl sm:flex-row sm:items-center sm:justify-between">
          <div className="flex items-center gap-4">
            <Link href="/admin" className="rounded-full border border-white/10 bg-slate-900/60 px-3 py-2 text-sm font-medium text-slate-200 transition hover:bg-slate-900">
              ← Back
            </Link>
            <div>
              <h1 className="text-2xl font-bold">Orders Management</h1>
              <p className="text-sm text-slate-400">Review customer orders and order details</p>
            </div>
          </div>
        </header>

        <section className="mb-6 grid gap-4 md:grid-cols-3">
          <div className="rounded-2xl border border-white/10 bg-slate-900/80 p-5 shadow-lg shadow-slate-950/25">
            <p className="text-sm text-slate-400">Total Orders</p>
            <p className="mt-3 text-3xl font-bold text-white">{orders.length}</p>
          </div>
          <div className="rounded-2xl border border-sky-400/30 bg-sky-500/10 p-5 shadow-lg shadow-slate-950/25">
            <p className="text-sm text-sky-200">Revenue</p>
            <p className="mt-3 text-3xl font-bold text-white">{formatMoney(totalRevenue)}</p>
          </div>
          <div className="rounded-2xl border border-emerald-400/30 bg-emerald-500/10 p-5 shadow-lg shadow-slate-950/25">
            <p className="text-sm text-emerald-200">Delivered</p>
            <p className="mt-3 text-3xl font-bold text-white">
              {orders.filter((order) => (order.order_status || '').toLowerCase() === 'delivered').length}
            </p>
          </div>
        </section>

        <main className="space-y-4">
          {loading ? (
            <div className="rounded-2xl border border-white/10 bg-slate-900/60 p-8 text-center text-slate-300">
              Loading orders...
            </div>
          ) : orders.length === 0 ? (
            <div className="rounded-2xl border border-white/10 bg-slate-900/60 p-12 text-center text-slate-300">
              No orders found.
            </div>
          ) : (
            orders.map((order) => {
              const customer = users.find((user) => String(user.id) === String(order.user_id));
              const isExpanded = expandedIds.includes(order.id);
              const items = itemsByOrder[order.id] ?? [];

              return (
                <article
                  key={order.id}
                  className="overflow-hidden rounded-3xl border border-white/10 bg-slate-900/70 shadow-2xl shadow-slate-950/30"
                >
                  <div className="p-5">
                    <div className="flex flex-col gap-4 xl:flex-row xl:items-center xl:justify-between">
                      <div className="flex items-center gap-4">
                        <div className="flex h-12 w-12 items-center justify-center overflow-hidden rounded-full bg-gradient-to-br from-cyan-500 to-sky-500 font-semibold text-white">
                          {customer?.profile_picture ? (
                            <img src={customer.profile_picture} alt={getDisplayName(customer)} className="h-full w-full object-cover" />
                          ) : (
                            <span>{getInitials(customer)}</span>
                          )}
                        </div>

                        <div>
                          <p className="text-lg font-semibold text-white">{getDisplayName(customer)}</p>
                          <p className="text-sm text-slate-400">{customer?.email || 'No email'} · {customer?.phone || 'No phone'}</p>
                        </div>
                      </div>

                      <div className="flex flex-wrap items-center gap-3 xl:justify-end">
                        <span className={`inline-flex rounded-full border px-2.5 py-1 text-xs font-medium ${getOrderStatusClasses(order.order_status)}`}>
                          {order.order_status || 'pending'}
                        </span>
                        <span className="rounded-full border border-white/10 bg-white/5 px-2.5 py-1 text-xs text-slate-300">
                          Order #{order.id}
                        </span>
                      </div>
                    </div>

                    <div className="mt-5 grid gap-4 md:grid-cols-2 xl:grid-cols-4">
                      <div className="rounded-2xl border border-white/10 bg-slate-950/50 p-3">
                        <p className="text-xs uppercase tracking-[0.2em] text-slate-500">Placed</p>
                        <p className="mt-2 text-sm text-slate-200">{formatDate(order.created_at)}</p>
                      </div>

                      <div className="rounded-2xl border border-white/10 bg-slate-950/50 p-3">
                        <p className="text-xs uppercase tracking-[0.2em] text-slate-500">Payment</p>
                        <p className="mt-2 text-sm text-slate-200">{order.payment_method || '—'}</p>
                        <p className="mt-1 text-xs text-slate-400">Status: {order.payment_status || '—'}</p>
                      </div>

                      <div className="rounded-2xl border border-white/10 bg-slate-950/50 p-3">
                        <p className="text-xs uppercase tracking-[0.2em] text-slate-500">Amount</p>
                        <p className="mt-2 text-sm text-slate-200">{formatMoney(Number(order.final_amount ?? order.total_amount ?? 0))}</p>
                        <p className="mt-1 text-xs text-slate-400">Discount: {formatMoney(Number(order.total_discount_amount ?? order.offer_discount ?? order.coupon_discount ?? 0))}</p>
                      </div>

                      <div className="rounded-2xl border border-white/10 bg-slate-950/50 p-3">
                        <p className="text-xs uppercase tracking-[0.2em] text-slate-500">Address</p>
                        <p className="mt-2 text-sm text-slate-200">
                          {order.address?.street_name || '—'}
                          {order.address?.city ? `, ${order.address.city}` : ''}
                        </p>
                        <p className="mt-1 text-xs text-slate-400">
                          {order.address?.state || '—'}
                          {order.address?.country ? `, ${order.address.country}` : ''}
                          {order.address?.pin_code ? ` · ${order.address.pin_code}` : ''}
                        </p>
                      </div>
                    </div>

                    <div className="mt-5 grid gap-4 rounded-2xl border border-white/10 bg-slate-950/40 p-4 md:grid-cols-2 xl:grid-cols-4">
                      <div>
                        <p className="text-[10px] uppercase tracking-[0.2em] text-slate-500">Order status</p>
                        <p className="mt-2 text-sm font-medium text-slate-100">{order.order_status || '—'}</p>
                      </div>
                      <div>
                        <p className="text-[10px] uppercase tracking-[0.2em] text-slate-500">Payment status</p>
                        <p className="mt-2 text-sm font-medium text-slate-100">{order.payment_status || '—'}</p>
                      </div>
                      <div>
                        <p className="text-[10px] uppercase tracking-[0.2em] text-slate-500">Delivery status</p>
                        <p className="mt-2 text-sm font-medium text-slate-100">{order.delivery_status || '—'}</p>
                      </div>
                      <div>
                        <p className="text-[10px] uppercase tracking-[0.2em] text-slate-500">Coupon</p>
                        <p className="mt-2 text-sm font-medium text-slate-100">{order.coupon_code || 'No coupon'}</p>
                      </div>
                    </div>

                    <div className="mt-5 grid gap-4 rounded-2xl border border-white/10 bg-slate-950/40 p-4 md:grid-cols-2 xl:grid-cols-5">
                      <div>
                        <p className="text-[10px] uppercase tracking-[0.2em] text-slate-500">Subtotal</p>
                        <p className="mt-2 text-sm text-slate-100">{formatMoney(Number(order.total_amount ?? 0))}</p>
                      </div>
                      <div>
                        <p className="text-[10px] uppercase tracking-[0.2em] text-slate-500">Offer discount</p>
                        <p className="mt-2 text-sm text-slate-100">{formatMoney(Number(order.offer_discount ?? order.offer_applied ?? 0))}</p>
                      </div>
                      <div>
                        <p className="text-[10px] uppercase tracking-[0.2em] text-slate-500">Coupon discount</p>
                        <p className="mt-2 text-sm text-slate-100">{formatMoney(Number(order.coupon_discount ?? 0))}</p>
                      </div>
                      <div>
                        <p className="text-[10px] uppercase tracking-[0.2em] text-slate-500">Total discount</p>
                        <p className="mt-2 text-sm text-slate-100">{formatMoney(Number(order.total_discount_amount ?? 0))}</p>
                      </div>
                      <div>
                        <p className="text-[10px] uppercase tracking-[0.2em] text-slate-500">Final amount</p>
                        <p className="mt-2 text-sm font-semibold text-emerald-300">{formatMoney(Number(order.final_amount ?? order.total_amount ?? 0))}</p>
                      </div>
                    </div>

                    <div className="mt-5 flex flex-wrap items-center justify-between gap-3 border-t border-white/10 pt-4">
                      <div className="text-sm text-slate-400">
                        {order.address?.phone ? `Contact: ${order.address.phone}` : 'No address phone'}
                        {order.address?.default ? ' · Default address' : ''}
                      </div>

                      <button
                        type="button"
                        onClick={() => toggleOrder(order)}
                        disabled={loadingItemsId === order.id}
                        className="rounded-xl border border-cyan-400/30 bg-cyan-500/10 px-4 py-2 text-sm font-medium text-cyan-100 transition hover:bg-cyan-500/20 disabled:cursor-not-allowed disabled:opacity-60"
                      >
                        {loadingItemsId === order.id
                          ? 'Loading...'
                          : isExpanded
                            ? 'Hide items'
                            : 'View items'}
                      </button>
                    </div>
                  </div>

                  {isExpanded && (
                    <div className="border-t border-white/10 bg-slate-950/50 p-5">
                      {items.length === 0 ? (
                        <div className="rounded-xl border border-dashed border-white/10 p-6 text-center text-slate-400">
                          No items found for this order.
                        </div>
                      ) : (
                        <div className="space-y-3">
                          {items.map((item) => (
                            <div key={item.id} className="rounded-2xl border border-white/10 bg-slate-900/80 p-4">
                              <div className="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
                                <div>
                                  <p className="text-base font-semibold text-white">{item.product_name || 'Product'}</p>
                                  <p className="text-sm text-slate-400">
                                    Qty: {item.quantity ?? 0} · Size: {item.size || 'N/A'}
                                  </p>
                                </div>

                                <div className="flex flex-wrap items-center gap-3 text-sm text-slate-300">
                                  <span>Price: {formatMoney(Number(item.price ?? 0))}</span>
                                  <span>Total: {formatMoney(Number(item.total_amount ?? 0))}</span>
                                </div>
                              </div>

                              <div className="mt-3 flex flex-wrap gap-2 text-xs text-slate-300">
                                <span className="rounded-full border border-white/10 bg-white/5 px-2 py-1">
                                  Status: {item.order_status || 'pending'}
                                </span>
                                <span className="rounded-full border border-white/10 bg-white/5 px-2 py-1">
                                  Delivery: {item.delivery_status || '—'}
                                </span>
                                <span className="rounded-full border border-white/10 bg-white/5 px-2 py-1">
                                  Payment: {item.payment_method || '—'}
                                </span>
                              </div>
                            </div>
                          ))}
                        </div>
                      )}
                    </div>
                  )}
                </article>
              );
            })
          )}
        </main>
      </div>
    </div>
  );
}
