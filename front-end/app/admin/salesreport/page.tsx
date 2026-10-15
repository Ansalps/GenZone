'use client';

import axios from 'axios';
import Link from 'next/link';
import { useEffect, useMemo, useState } from 'react';
import { toast } from 'sonner';

const API = 'http://localhost:8080';

type SalesReportRow = {
  order_id?: number | string;
  product_id?: number | string;
  product_name?: string;
  qty?: number;
  price?: number;
  order_status?: string;
  payment_method?: string;
  coupon_discount?: number;
  offer_discount?: number;
  total_discount?: number;
  paid_amount?: number;
  order_date?: string;
  delivered_date?: string;
};

const money = (value: number = 0) =>
  new Intl.NumberFormat('en-IN', {
    style: 'currency',
    currency: 'INR',
    maximumFractionDigits: 2,
  }).format(value);

export default function SalesReportPage() {
  const [rows, setRows] = useState<SalesReportRow[]>([]);
  const [filters, setFilters] = useState({ day: '', month: '', year: '' });
  const [loading, setLoading] = useState(false);
  const [generating, setGenerating] = useState(false);
  const [downloading, setDownloading] = useState<'pdf' | 'xlsx' | null>(null);
  const [summary, setSummary] = useState({
    totalQuantity: 0,
    totalPaidAmount: 0,
    totalDiscount: 0,
  });

  const activeParams = useMemo(() => {
    const params: Record<string, string> = {};
    if (filters.day) params.day = filters.day;
    if (filters.month) params.month = filters.month;
    if (filters.year) params.year = filters.year;
    return params;
  }, [filters]);

  const loadSalesReport = async (nextFilters = filters) => {
    setLoading(true);

    try {
      const params = {
        ...(nextFilters.day ? { day: nextFilters.day } : {}),
        ...(nextFilters.month ? { month: nextFilters.month } : {}),
        ...(nextFilters.year ? { year: nextFilters.year } : {}),
      };

      const response = await axios.get(`${API}/admin/salesreport`, {
        params,
        withCredentials: true,
      });

      const reportRows = Array.isArray(response.data?.sales_report) ? response.data.sales_report : [];
      setRows(reportRows);
      setSummary({
        totalQuantity: Number(response.data?.overall_sales_count ?? 0),
        totalPaidAmount: Number(response.data?.overall_order_amount ?? 0),
        totalDiscount: Number(response.data?.overall_discount ?? 0),
      });
    } catch (error) {
      console.error('Failed to load sales report:', error);
      setRows([]);
      toast.error('Unable to fetch the sales report');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadSalesReport();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const handleGenerate = async () => {
    setGenerating(true);
    try {
      await axios.post(`${API}/admin/salesreport`, {}, { withCredentials: true });
      toast.success('Sales report generated successfully');
      await loadSalesReport();
    } catch (error) {
      console.error('Failed to generate sales report:', error);
      toast.error('Unable to generate sales report');
    } finally {
      setGenerating(false);
    }
  };

  const handleDownload = async (format: 'pdf' | 'xlsx') => {
    try {
      setDownloading(format);
      const response = await axios.get(`${API}/admin/salesreportdownload`, {
        params: {
          ...activeParams,
          format,
        },
        withCredentials: true,
        responseType: 'blob',
      });

      const blob = new Blob([response.data], {
        type: format === 'pdf' ? 'application/pdf' : 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
      });

      const url = window.URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.href = url;
      link.download = `sales-report.${format === 'pdf' ? 'pdf' : 'xlsx'}`;
      link.click();
      window.URL.revokeObjectURL(url);
      toast.success('Report download started');
    } catch (error) {
      console.error('Failed to download report:', error);
      toast.error('Unable to download the sales report');
    } finally {
      setDownloading(null);
    }
  };

  const handleFilterSubmit = async () => {
    await loadSalesReport(filters);
  };

  const handleReset = async () => {
    const resetValues = { day: '', month: '', year: '' };
    setFilters(resetValues);
    await loadSalesReport(resetValues);
  };

  return (
    <div className="min-h-screen bg-slate-950 px-4 py-8 text-white sm:px-6 lg:px-8">
      <div className="mx-auto max-w-7xl">
        <header className="mb-6 flex flex-col gap-4 rounded-3xl border border-white/10 bg-white/5 p-5 backdrop-blur-xl sm:flex-row sm:items-center sm:justify-between">
          <div className="flex items-center gap-4">
            <Link href="/admin" className="rounded-full border border-white/10 bg-slate-900/60 px-3 py-2 text-sm font-medium text-slate-200 transition hover:bg-slate-900">
              ← Back
            </Link>
            <div>
              <p className="text-sm uppercase tracking-[0.2em] text-cyan-400">Admin</p>
              <h1 className="text-2xl font-bold">Sales Report</h1>
            </div>
          </div>

          <div className="flex flex-wrap items-center gap-3">
            <button
              type="button"
              onClick={handleGenerate}
              disabled={generating}
              className="rounded-xl bg-gradient-to-r from-emerald-500 to-teal-500 px-4 py-2 text-sm font-semibold text-white shadow-lg shadow-emerald-500/20 transition hover:brightness-105 disabled:cursor-not-allowed disabled:opacity-60"
            >
              {generating ? 'Generating...' : 'Generate Report'}
            </button>
            <button
              type="button"
              onClick={() => handleDownload('pdf')}
              disabled={downloading !== null}
              className="rounded-xl border border-white/10 bg-slate-900/70 px-4 py-2 text-sm font-medium text-slate-100 transition hover:bg-slate-800 disabled:cursor-not-allowed disabled:opacity-60"
            >
              {downloading === 'pdf' ? 'Preparing PDF...' : 'Download PDF'}
            </button>
            <button
              type="button"
              onClick={() => handleDownload('xlsx')}
              disabled={downloading !== null}
              className="rounded-xl border border-cyan-500/30 bg-cyan-500/10 px-4 py-2 text-sm font-medium text-cyan-100 transition hover:bg-cyan-500/20 disabled:cursor-not-allowed disabled:opacity-60"
            >
              {downloading === 'xlsx' ? 'Preparing Excel...' : 'Download Excel'}
            </button>
          </div>
        </header>

        <section className="mb-6 grid gap-4 rounded-3xl border border-white/10 bg-slate-900/60 p-5 shadow-xl shadow-slate-950/30 md:grid-cols-4">
          <label className="text-sm text-slate-300">
            Day
            <input
              type="date"
              value={filters.day}
              onChange={(e) => setFilters((prev) => ({ ...prev, day: e.target.value }))}
              className="mt-2 w-full rounded-xl border border-white/10 bg-slate-950 px-3 py-2 text-white outline-none"
            />
          </label>
          <label className="text-sm text-slate-300">
            Month
            <input
              type="number"
              min={1}
              max={12}
              value={filters.month}
              onChange={(e) => setFilters((prev) => ({ ...prev, month: e.target.value }))}
              placeholder="1-12"
              className="mt-2 w-full rounded-xl border border-white/10 bg-slate-950 px-3 py-2 text-white outline-none"
            />
          </label>
          <label className="text-sm text-slate-300">
            Year
            <input
              type="number"
              value={filters.year}
              onChange={(e) => setFilters((prev) => ({ ...prev, year: e.target.value }))}
              placeholder="2026"
              className="mt-2 w-full rounded-xl border border-white/10 bg-slate-950 px-3 py-2 text-white outline-none"
            />
          </label>
          <div className="flex items-end gap-2">
            <button
              type="button"
              onClick={handleFilterSubmit}
              className="flex-1 rounded-xl bg-white/10 px-4 py-2 text-sm font-medium text-white transition hover:bg-white/15"
            >
              Apply
            </button>
            <button
              type="button"
              onClick={handleReset}
              className="flex-1 rounded-xl border border-white/10 bg-slate-900 px-4 py-2 text-sm font-medium text-slate-200 transition hover:bg-slate-800"
            >
              Reset
            </button>
          </div>
        </section>

        <section className="mb-6 grid gap-4 md:grid-cols-3">
          <div className="rounded-2xl border border-white/10 bg-slate-900/70 p-5">
            <p className="text-sm text-slate-400">Total quantity</p>
            <h2 className="mt-3 text-3xl font-bold text-cyan-300">{summary.totalQuantity}</h2>
          </div>
          <div className="rounded-2xl border border-white/10 bg-slate-900/70 p-5">
            <p className="text-sm text-slate-400">Overall order amount</p>
            <h2 className="mt-3 text-3xl font-bold text-emerald-300">{money(summary.totalPaidAmount)}</h2>
          </div>
          <div className="rounded-2xl border border-white/10 bg-slate-900/70 p-5">
            <p className="text-sm text-slate-400">Discounts</p>
            <h2 className="mt-3 text-3xl font-bold text-pink-300">{money(summary.totalDiscount)}</h2>
          </div>
        </section>

        <main className="overflow-hidden rounded-3xl border border-white/10 bg-slate-900/60 shadow-2xl shadow-slate-950/30">
          {loading ? (
            <div className="p-8 text-center text-slate-300">Loading sales report...</div>
          ) : rows.length === 0 ? (
            <div className="p-8 text-center text-slate-300">No sales data found for the selected filter.</div>
          ) : (
            <div className="overflow-x-auto">
              <table className="min-w-full text-left text-sm text-slate-200">
                <thead className="bg-slate-800/80 text-xs uppercase tracking-[0.12em] text-slate-400">
                  <tr>
                    <th className="px-4 py-3">Order</th>
                    <th className="px-4 py-3">Product</th>
                    <th className="px-4 py-3">Qty</th>
                    <th className="px-4 py-3">Price</th>
                    <th className="px-4 py-3">Status</th>
                    <th className="px-4 py-3">Payment</th>
                    <th className="px-4 py-3">Discount</th>
                    <th className="px-4 py-3">Paid</th>
                    <th className="px-4 py-3">Order Date</th>
                  </tr>
                </thead>
                <tbody>
                  {rows.map((row, index) => (
                    <tr key={`${row.order_id ?? 'o'}-${row.product_id ?? 'p'}-${index}`} className="border-t border-white/5 hover:bg-slate-800/40">
                      <td className="px-4 py-3 font-medium text-slate-100">#{row.order_id ?? '—'}</td>
                      <td className="px-4 py-3">
                        <div className="font-medium text-slate-100">{row.product_name ?? '—'}</div>
                        <div className="text-xs text-slate-400">Product #{row.product_id ?? '—'}</div>
                      </td>
                      <td className="px-4 py-3">{row.qty ?? 0}</td>
                      <td className="px-4 py-3">{money(Number(row.price ?? 0))}</td>
                      <td className="px-4 py-3">
                        <span className="rounded-full border border-cyan-400/30 bg-cyan-500/10 px-2 py-1 text-xs text-cyan-100">
                          {row.order_status ?? '—'}
                        </span>
                      </td>
                      <td className="px-4 py-3">{row.payment_method ?? '—'}</td>
                      <td className="px-4 py-3">{money(Number(row.total_discount ?? 0))}</td>
                      <td className="px-4 py-3">{money(Number(row.paid_amount ?? 0))}</td>
                      <td className="px-4 py-3">{row.order_date ? new Date(row.order_date).toLocaleDateString('en-IN') : '—'}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </main>
      </div>
    </div>
  );
}
