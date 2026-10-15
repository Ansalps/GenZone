'use client';

import axios from 'axios';
import Link from 'next/link';
import { useState } from 'react';
import { toast } from 'sonner';

const API = 'http://localhost:8080';

export default function InvoicePage() {
  const [orderId, setOrderId] = useState('');
  const [loading, setLoading] = useState(false);

  const handlePreview = async () => {
    const trimmedOrderId = orderId.trim();

    if (!trimmedOrderId) {
      toast.error('Enter an order ID to generate the invoice');
      return;
    }

    try {
      setLoading(true);
      const response = await axios.get(`${API}/admin/invoice/${trimmedOrderId}`, {
        withCredentials: true,
        responseType: 'blob',
      });

      const blob = new Blob([response.data], { type: 'application/pdf' });
      const url = window.URL.createObjectURL(blob);
      const newWindow = window.open(url, '_blank', 'noopener,noreferrer');

      if (!newWindow) {
        toast.error('Popup was blocked. Please allow popups to view the invoice');
      }
    } catch (error: any) {
      console.error('Failed to generate invoice:', error);
      const message = error?.response?.data?.message || 'Unable to generate invoice for this order';
      toast.error(message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-slate-950 px-4 py-8 text-white sm:px-6 lg:px-8">
      <div className="mx-auto max-w-3xl">
        <header className="mb-6 flex items-center justify-between rounded-3xl border border-white/10 bg-white/5 p-5 backdrop-blur-xl">
          <div className="flex items-center gap-4">
            <Link href="/admin" className="rounded-full border border-white/10 bg-slate-900/60 px-3 py-2 text-sm font-medium text-slate-200 transition hover:bg-slate-900">
              ← Back
            </Link>
            <div>
              <p className="text-sm uppercase tracking-[0.2em] text-fuchsia-400">Finance</p>
              <h1 className="text-2xl font-bold">Invoices</h1>
            </div>
          </div>
        </header>

        <main className="rounded-3xl border border-white/10 bg-slate-900/60 p-6 shadow-2xl shadow-slate-950/30">
          <div className="mb-4 text-sm text-slate-300">
            Fetch the generated invoice PDF for any completed order using the backend invoice API.
          </div>

          <div className="flex flex-col gap-3 sm:flex-row">
            <input
              type="number"
              value={orderId}
              onChange={(e) => setOrderId(e.target.value)}
              placeholder="Enter order ID"
              className="flex-1 rounded-xl border border-white/10 bg-slate-950 px-4 py-3 text-white outline-none"
            />
            <button
              type="button"
              onClick={handlePreview}
              disabled={loading}
              className="rounded-xl bg-gradient-to-r from-fuchsia-500 to-pink-500 px-5 py-3 text-sm font-semibold text-white shadow-lg shadow-fuchsia-500/20 transition hover:brightness-105 disabled:cursor-not-allowed disabled:opacity-60"
            >
              {loading ? 'Generating...' : 'Generate Invoice'}
            </button>
          </div>
        </main>
      </div>
    </div>
  );
}
