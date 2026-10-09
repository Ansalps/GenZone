'use client';

import axios from 'axios';
import Link from 'next/link';
import { useRouter, useSearchParams } from 'next/navigation';
import { useEffect, useState } from 'react';

const API = 'http://localhost:8080';

export default function PaymentOptionPage() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const addressId = Number(searchParams.get('address_id') ?? '0');
  const couponCode = searchParams.get('coupon_code') ?? '';
  const [selectedMethod, setSelectedMethod] = useState<'cod' | 'internet-banking'>('cod');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState('');
  const [cartTotal, setCartTotal] = useState(0);
  const [couponDiscount, setCouponDiscount] = useState(0);
  const [loadingAmount, setLoadingAmount] = useState(true);

  const isCod = selectedMethod === 'cod';
  const amountDue = Math.max(cartTotal - couponDiscount, 0);

  useEffect(() => {
    const loadAmount = async () => {
      try {
        setLoadingAmount(true);
        const summaryResponse = await axios.get(`${API}/cart/summary`, { withCredentials: true });
        const summary = summaryResponse.data?.data ?? {};
        const baseTotal = Number(summary.final_total ?? summary.original_total ?? 0);

        let couponValue = 0;
        if (couponCode.trim()) {
          const couponResponse = await axios.get(`${API}/checkout/coupon`, {
            withCredentials: true,
            params: { coupon_code: couponCode.trim() },
          });
          couponValue = Number(couponResponse.data?.data?.coupon_discount ?? 0);
        }

        setCartTotal(baseTotal);
        setCouponDiscount(couponValue);
      } catch (err) {
        console.error('Failed to load payment total', err);
        setCartTotal(0);
        setCouponDiscount(0);
      } finally {
        setLoadingAmount(false);
      }
    };

    void loadAmount();
  }, [couponCode]);

  const handlePlaceOrder = async () => {
    if (!addressId) {
      setError('Please select a delivery address from checkout before placing the order.');
      return;
    }

    setIsSubmitting(true);
    setError('');

    try {
      const payload = {
        address_id: addressId,
        coupon_code: couponCode,
        payment_method: isCod ? 'COD' : 'RazorPay',
      };

      await axios.post(`${API}/checkout/order`, payload, { withCredentials: true });
      router.push('/orders');
    } catch (err: any) {
      setError(err.response?.data?.message || 'Unable to place order right now. Please try again.');
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="min-h-screen bg-slate-950 px-6 py-10 text-white">
      <div className="mx-auto max-w-2xl rounded-3xl border border-white/10 bg-slate-900/80 p-6 shadow-2xl shadow-cyan-950/30 sm:p-8">
        <div className="mb-8">
          <p className="text-xs font-semibold uppercase tracking-[0.32em] text-cyan-300">Payment</p>
          <h1 className="mt-3 text-3xl font-bold text-white">Choose payment option</h1>
        </div>

        <div className="mb-6 rounded-2xl border border-white/10 bg-slate-950/60 p-4 shadow-inner shadow-cyan-950/20">
          <div className="flex items-center justify-between text-sm text-slate-300">
            <span>Amount to be paid</span>
            <span className="text-xl font-bold text-white">
              {loadingAmount ? 'Loading...' : `₹${amountDue.toFixed(2)}`}
            </span>
          </div>
          <div className="mt-2 flex items-center justify-between text-xs text-slate-400">
            <span>Cart total</span>
            <span>₹{cartTotal.toFixed(2)}</span>
          </div>
          <div className="mt-1 flex items-center justify-between text-xs text-emerald-300">
            <span>Coupon discount</span>
            <span>-₹{couponDiscount.toFixed(2)}</span>
          </div>
        </div>

        <div className="space-y-4">
          <button
            type="button"
            onClick={() => setSelectedMethod('cod')}
            className={`w-full rounded-2xl border p-5 text-left transition ${
              isCod
                ? 'border-cyan-400/60 bg-cyan-500/10 shadow-lg shadow-cyan-500/10'
                : 'border-white/10 bg-slate-950/60 hover:border-cyan-400/40'
            }`}
          >
            <div className="flex items-center justify-between gap-4">
              <div>
                <p className="text-lg font-semibold text-white">Cash On Delivery</p>
                <p className="mt-1 text-sm text-slate-400">Pay when your order is delivered.</p>
              </div>
              <span
                className={`flex h-5 w-5 items-center justify-center rounded-full border ${
                  isCod ? 'border-cyan-400 bg-cyan-400' : 'border-slate-500 bg-transparent'
                }`}
              >
                {isCod && <span className="h-2.5 w-2.5 rounded-full bg-slate-950" />}
              </span>
            </div>
          </button>

          <button
            type="button"
            onClick={() => setSelectedMethod('internet-banking')}
            className={`w-full rounded-2xl border p-5 text-left transition ${
              !isCod
                ? 'border-violet-400/60 bg-violet-500/10 shadow-lg shadow-violet-500/10'
                : 'border-white/10 bg-slate-950/60 hover:border-violet-400/40'
            }`}
          >
            <div className="flex items-center justify-between gap-4">
              <div>
                <p className="text-lg font-semibold text-white">Internet Banking</p>
                <p className="mt-1 text-sm text-slate-400">Complete payment through your bank online.</p>
              </div>
              <span
                className={`flex h-5 w-5 items-center justify-center rounded-full border ${
                  !isCod ? 'border-violet-400 bg-violet-400' : 'border-slate-500 bg-transparent'
                }`}
              >
                {!isCod && <span className="h-2.5 w-2.5 rounded-full bg-slate-950" />}
              </span>
            </div>
          </button>
        </div>

        {error && (
          <p className="mt-4 rounded-xl border border-red-500/40 bg-red-500/10 px-3 py-2 text-sm text-red-200">
            {error}
          </p>
        )}

        <div className="mt-8 flex flex-col gap-3 sm:flex-row sm:justify-between">
          <Link
            href="/checkout"
            className="inline-flex items-center justify-center rounded-xl border border-white/10 bg-white/5 px-4 py-3 text-sm font-medium text-slate-100 hover:border-cyan-400/40 hover:bg-cyan-500/10"
          >
            Back to checkout
          </Link>

          <button
            type="button"
            onClick={handlePlaceOrder}
            disabled={isSubmitting}
            className={`inline-flex items-center justify-center rounded-xl px-5 py-3 text-sm font-semibold text-white shadow-lg transition disabled:cursor-not-allowed disabled:opacity-60 ${
              isCod
                ? 'bg-gradient-to-r from-cyan-500 to-sky-500 hover:brightness-110'
                : 'bg-gradient-to-r from-violet-500 to-purple-500 hover:brightness-110'
            }`}
          >
            {isSubmitting ? 'Processing...' : isCod ? 'Place Order' : 'Buy Now'}
          </button>
        </div>
      </div>
    </div>
  );
}
