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
  product_variant_id?: number;
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
  final_amount?: number;
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
  payment_status?: string;
  order_status?: string;
  delivery_status?: string;
  coupon_code?: string;
  address?: Address;
  items?: OrderItem[];
};

export default function OrdersPage() {
  const [orders, setOrders] = useState<Order[]>([]);
  const [orderItems, setOrderItems] = useState<Record<number, OrderItem[]>>({});
  const [expandedOrders, setExpandedOrders] = useState<Record<number, boolean>>({});
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [sortOrder, setSortOrder] = useState<'DSC' | 'ASC'>('DSC');
  const [cancelingOrderId, setCancelingOrderId] = useState<number | null>(null);
  const [cancelingItemId, setCancelingItemId] = useState<number | null>(null);

  const fetchOrderItems = async (orderId: number) => {
    try {
      const response = await axios.get(`${API}/order-items/${orderId}`, { withCredentials: true });
      const items = Array.isArray(response.data?.['order items']) ? response.data['order items'] : [];
      setOrderItems((prev) => ({ ...prev, [orderId]: items }));
      return items;
    } catch (err) {
      console.error(`Failed to load items for order ${orderId}`, err);
      return [];
    }
  };

  const toggleOrderItems = async (orderId: number) => {
    const isExpanded = !!expandedOrders[orderId];
    if (!isExpanded && !(orderItems[orderId]?.length)) {
      await fetchOrderItems(orderId);
    }

    setExpandedOrders((prev) => ({
      ...prev,
      [orderId]: !prev[orderId],
    }));
  };

  const refreshOrders = async (selectedOrder: 'DSC' | 'ASC' = sortOrder) => {
    try {
      const response = await axios.get(`${API}/orders`, {
        withCredentials: true,
        params: { list_order: selectedOrder },
      });
      const orderList = Array.isArray(response.data?.data?.Order) ? response.data.data.Order : [];
      setOrders(orderList);

      for (const order of orderList) {
        await fetchOrderItems(order.id);
      }
    } catch (err: any) {
      console.error('Failed to refresh orders', err);
      setError(err.response?.data?.message || 'Unable to refresh orders right now.');
      setOrders([]);
    }
  };

  const canCancelByDeliveryStatus = (status?: string) => (status || 'pending').toLowerCase() === 'pending';
  const canReturnByDeliveryStatus = (status?: string) => (status || 'pending').toLowerCase() === 'delivered';

  const returnOrder = async (orderId: number) => {
    if (!orderId) return;

    const confirmed = window.confirm('Are you sure you want to return this entire order?');
    if (!confirmed) return;

    try {
      const response = await axios.put(`${API}/orders/${orderId}/return`, {}, { withCredentials: true });
      const message = response.data?.message || 'Order return processed successfully.';
      setError('');
      await refreshOrders();
      alert(message);
    } catch (err: any) {
      console.error('Failed to return order', err);
      alert(err.response?.data?.message || 'Unable to return this order right now.');
    }
  };

  const cancelOrder = async (order: Order) => {
    const orderId = order.id;
    if (!orderId) return;
    if (!canCancelByDeliveryStatus(order.delivery_status)) {
      alert('Order cancellation is only available while delivery status is pending.');
      return;
    }

    const confirmed = window.confirm('Are you sure you want to cancel this order?');
    if (!confirmed) return;

    setCancelingOrderId(orderId);
    try {
      const response = await axios.put(`${API}/orders/${orderId}/cancel`, {}, { withCredentials: true });
      const message = response.data?.message || 'Order cancelled successfully.';
      setError('');
      await refreshOrders();
      alert(message);
    } catch (err: any) {
      console.error('Failed to cancel order', err);
      alert(err.response?.data?.message || 'Unable to cancel this order right now.');
    } finally {
      setCancelingOrderId(null);
    }
  };

  const returnSingleOrderItem = async (item: OrderItem) => {
    const itemId = item.id;
    if (!itemId) return;
    if (!canReturnByDeliveryStatus(item.delivery_status)) {
      alert('Item return is only available when delivery status is delivered.');
      return;
    }

    const maxQuantity = Number(item.quantity ?? 1);
    const enteredQuantity = window.prompt(
      `How many quantity would you like to return? Enter a number between 1 and ${maxQuantity}.`,
      String(maxQuantity),
    );

    if (enteredQuantity === null) return;

    const quantity = Number(enteredQuantity);
    if (!Number.isInteger(quantity) || quantity < 1 || quantity > maxQuantity) {
      alert(`Please enter a valid quantity between 1 and ${maxQuantity}.`);
      return;
    }

    const confirmed = window.confirm(`Return ${quantity} quantity from this item?`);
    if (!confirmed) return;

    setCancelingItemId(itemId);
    try {
      const response = await axios.put(
        `${API}/orders/item/${itemId}/return`,
        { quantity },
        { withCredentials: true },
      );
      const message = response.data?.message || 'Item return processed successfully.';
      setError('');
      await refreshOrders();
      alert(message);
    } catch (err: any) {
      console.error('Failed to return order item', err);
      alert(err.response?.data?.message || 'Unable to return this item right now.');
    } finally {
      setCancelingItemId(null);
    }
  };

  const cancelSingleOrderItem = async (item: OrderItem) => {
    const itemId = item.id;
    if (!itemId) return;
    if (!canCancelByDeliveryStatus(item.delivery_status)) {
      alert('Item cancellation is only available while delivery status is pending.');
      return;
    }

    const maxQuantity = Number(item.quantity ?? 1);
    const enteredQuantity = window.prompt(
      `How many quantity would you like to cancel? Enter a number between 1 and ${maxQuantity}.`,
      String(maxQuantity),
    );

    if (enteredQuantity === null) return;

    const quantity = Number(enteredQuantity);
    if (!Number.isInteger(quantity) || quantity < 1 || quantity > maxQuantity) {
      alert(`Please enter a valid quantity between 1 and ${maxQuantity}.`);
      return;
    }

    const confirmed = window.confirm(`Cancel ${quantity} quantity from this item?`);
    if (!confirmed) return;

    setCancelingItemId(itemId);
    try {
      const response = await axios.put(
        `${API}/orders/item/${itemId}/cancel`,
        { quantity },
        { withCredentials: true },
      );
      const message = response.data?.message || 'Item cancelled successfully.';
      setError('');
      await refreshOrders();
      alert(message);
    } catch (err: any) {
      console.error('Failed to cancel order item', err);
      alert(err.response?.data?.message || 'Unable to cancel this item right now.');
    } finally {
      setCancelingItemId(null);
    }
  };

  const retryPayment = (order: Order) => {
    if (!order.address?.id) {
      alert('Address is not available for this order.');
      return;
    }

    const params = new URLSearchParams({
      address_id: String(order.address.id),
      coupon_code: order.coupon_code ?? '',
      retry_order_id: String(order.id),
      retry_amount: String(Number(order.final_amount ?? order.total_amount ?? 0)),
    });

    window.location.href = `/payment-option?${params.toString()}`;
  };

  useEffect(() => {
    const fetchOrders = async () => {
      try {
        setLoading(true);
        const response = await axios.get(`${API}/orders`, {
          withCredentials: true,
          params: { list_order: sortOrder },
        });
        const orderList = Array.isArray(response.data?.data?.Order) ? response.data.data.Order : [];
        setOrders(orderList);

        for (const order of orderList) {
          await fetchOrderItems(order.id);
        }
      } catch (err: any) {
        console.error('Failed to load orders', err);
        setError(err.response?.data?.message || 'Unable to load orders right now.');
        setOrders([]);
      } finally {
        setLoading(false);
      }
    };

    fetchOrders();
  }, [sortOrder]);

  return (
    <div className="min-h-screen bg-slate-950 px-6 py-10 text-white">
      <div className="mx-auto max-w-5xl">
        <div className="mb-6 flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
          <div>
            <p className="text-xs font-semibold uppercase tracking-[0.28em] text-cyan-300">Your orders</p>
            <h1 className="mt-2 text-3xl font-bold text-white">Order history</h1>
          </div>
          <div className="flex items-center gap-3">
            <label className="flex items-center gap-2 rounded-xl border border-white/10 bg-slate-900/80 px-3 py-2 text-sm text-slate-200">
              <span>Sort</span>
              <select
                value={sortOrder}
                onChange={(e) => setSortOrder(e.target.value as 'DSC' | 'ASC')}
                className="rounded-lg border border-white/10 bg-slate-950 px-2 py-1.5 text-sm text-white outline-none"
              >
                <option value="DSC">Newest</option>
                <option value="ASC">Oldest</option>
              </select>
            </label>
            <Link
              href="/dashboard"
              className="rounded-xl border border-white/10 bg-white/5 px-4 py-2 text-sm font-medium text-slate-100 transition hover:border-cyan-400/40 hover:bg-cyan-500/10"
            >
              Back to dashboard
            </Link>
          </div>
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
            {orders.map((order) => {
              const paymentExpired =
                (order.payment_status?.toLowerCase() === 'expired' || order.order_status?.toLowerCase() === 'cancelled') ||
                (!!order.created_at &&
                  Date.now() - new Date(order.created_at).getTime() > 60000 &&
                  (order.payment_status?.toLowerCase() === 'pending' || order.payment_status?.toLowerCase() === 'failed'));

              return (
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
                        Order: {order.order_status || 'pending'}
                      </span>
                      <span className="rounded-full border border-amber-500/40 bg-amber-500/10 px-2.5 py-1 text-amber-200">
                        Delivery: {order.delivery_status || 'pending'}
                      </span>
                      <span className="rounded-full border border-emerald-500/40 bg-emerald-500/10 px-2.5 py-1 text-emerald-200">
                        Payment: {order.payment_status || 'pending'}
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

                  <div className="mt-4 flex flex-wrap items-center justify-end gap-3">
                    {(order.payment_method?.toLowerCase() === 'razorpay' || order.payment_method?.toLowerCase() === 'razorpay') &&
                      (order.payment_status?.toLowerCase() === 'pending' || order.payment_status?.toLowerCase() === 'failed' || order.order_status?.toLowerCase() === 'pending') && (
                        <button
                          type="button"
                          onClick={() => !paymentExpired && retryPayment(order)}
                          disabled={paymentExpired}
                          className={`rounded-lg border px-3 py-2 text-sm font-medium transition ${paymentExpired ? 'cursor-not-allowed border-red-500/40 bg-red-500/10 text-red-200' : 'border-amber-500/40 bg-amber-500/10 text-amber-200 hover:border-amber-300 hover:bg-amber-500/20'}`}
                        >
                          {paymentExpired ? 'Payment Expired' : 'Retry Payment'}
                        </button>
                      )}

                    <button
                      type="button"
                      onClick={() => toggleOrderItems(order.id)}
                      className="rounded-lg border border-cyan-500/40 bg-cyan-500/10 px-3 py-2 text-sm font-medium text-cyan-200 transition hover:border-cyan-300 hover:bg-cyan-500/20"
                    >
                      {expandedOrders[order.id] ? 'Minimize order items' : 'View order items'}
                    </button>

                    {canReturnByDeliveryStatus(order.delivery_status) && order.order_status !== 'cancelled' && (
                      <button
                        type="button"
                        onClick={() => returnOrder(order.id)}
                        className="rounded-lg border border-purple-500/40 bg-purple-500/10 px-3 py-2 text-sm font-medium text-purple-200 transition hover:border-purple-300 hover:bg-purple-500/20"
                      >
                        Return Order
                      </button>
                    )}

                    {canCancelByDeliveryStatus(order.delivery_status) && order.order_status !== 'cancelled' && (
                      <button
                        type="button"
                        onClick={() => cancelOrder(order)}
                        disabled={cancelingOrderId === order.id}
                        className="rounded-lg border border-red-500/40 bg-red-500/10 px-3 py-2 text-sm font-medium text-red-200 transition hover:border-red-300 hover:bg-red-500/20 disabled:cursor-not-allowed disabled:opacity-60"
                      >
                        {cancelingOrderId === order.id ? 'Cancelling...' : 'Cancel Order'}
                      </button>
                    )}
                  </div>

                  {(expandedOrders[order.id] ?? false) && (orderItems[order.id]?.length ?? 0) > 0 && (
                    <div className="mt-4 rounded-xl border border-white/10 bg-slate-950/60 p-4">
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
                                <div className="mt-1 flex flex-wrap items-center gap-2">
                                  <span className="text-xs text-cyan-300">Qty: {Number(item.quantity ?? 0)}</span>
                                  {item.size && <span className="rounded-full border border-white/10 bg-white/5 px-2 py-0.5 text-[11px] font-medium text-slate-200">Size: {item.size}</span>}
                                </div>
                              </div>

                              <div className="flex flex-wrap items-center gap-3 text-sm text-slate-300">
                                <span className="rounded-full border border-violet-500/40 bg-violet-500/10 px-2 py-1 text-violet-200">
                                  Order: {item.order_status || 'pending'}
                                </span>
                                <span className="rounded-full border border-amber-500/40 bg-amber-500/10 px-2 py-1 text-amber-200">
                                  Delivery: {item.delivery_status || 'pending'}
                                </span>
                                <span>Price: ₹{Number(item.price ?? 0).toFixed(2)}</span>
                                <span className="text-emerald-300">Total: ₹{Number(item.total_amount ?? (item.price ?? 0) * (item.quantity ?? 1)).toFixed(2)}</span>
                                <span className="text-cyan-300">Final: ₹{Number(item.final_amount ?? item.paid_amount ?? item.price ?? 0).toFixed(2)}</span>
                              </div>
                            </div>

                            <div className="mt-3 flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
                              <div className="grid flex-1 gap-2 text-xs text-slate-300 sm:grid-cols-3">
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

                              {canReturnByDeliveryStatus(item.delivery_status) && (item.order_status || '').toLowerCase() !== 'cancelled' && (
                                <button
                                  type="button"
                                  onClick={() => returnSingleOrderItem(item)}
                                  disabled={cancelingItemId === item.id}
                                  className="rounded-lg border border-purple-500/40 bg-purple-500/10 px-3 py-2 text-xs font-medium text-purple-200 transition hover:border-purple-300 hover:bg-purple-500/20 disabled:cursor-not-allowed disabled:opacity-60"
                                >
                                  {cancelingItemId === item.id ? 'Processing...' : 'Return'}
                                </button>
                              )}

                              {canCancelByDeliveryStatus(item.delivery_status) && (item.order_status || '').toLowerCase() !== 'cancelled' && (
                                <button
                                  type="button"
                                  onClick={() => cancelSingleOrderItem(item)}
                                  disabled={cancelingItemId === item.id}
                                  className="rounded-lg border border-red-500/40 bg-red-500/10 px-3 py-2 text-xs font-medium text-red-200 transition hover:border-red-300 hover:bg-red-500/20 disabled:cursor-not-allowed disabled:opacity-60"
                                >
                                  {cancelingItemId === item.id ? 'Cancelling...' : 'Cancel Item'}
                                </button>
                              )}
                            </div>
                          </div>
                        ))}
                      </div>
                    </div>
                  )}
                </div>
              );
            })}
          </div>
        )}
      </div>
    </div>
  );
}
