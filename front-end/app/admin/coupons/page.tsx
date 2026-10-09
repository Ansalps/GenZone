'use client';

import axios from 'axios';
import Link from 'next/link';
import { useEffect, useState } from 'react';
import { toast } from 'sonner';
import ConfirmModal from '@/components/admin/confirm-modal';

const API = 'http://localhost:8080';

type CouponRecord = {
    id: number;
    code: string;
    discount: number;
    min_purchase: number;
    start_at?: string;
    end_at?: string;
    is_active: boolean;
    created_at?: string;
    updated_at?: string;
};

const normalizeDate = (value?: string) => {
    if (!value) return '';
    const trimmed = value.trim();
    if (!trimmed) return '';
    if (/^\d{4}-\d{2}-\d{2}$/.test(trimmed)) return trimmed;
    const date = new Date(trimmed);
    if (Number.isNaN(date.getTime())) return '';
    return date.toISOString().slice(0, 10);
};

export default function CouponsPage() {
    const [coupons, setCoupons] = useState<CouponRecord[]>([]);
    const [loading, setLoading] = useState(false);
    const [submitting, setSubmitting] = useState(false);
    const [editingId, setEditingId] = useState<number | null>(null);
    const [code, setCode] = useState('');
    const [discount, setDiscount] = useState(0);
    const [minPurchase, setMinPurchase] = useState(0);
    const [startAt, setStartAt] = useState('');
    const [endAt, setEndAt] = useState('');
    const [isActive, setIsActive] = useState(false);
    const [confirmOpen, setConfirmOpen] = useState(false);
    const [confirmId, setConfirmId] = useState<number | null>(null);
    const [confirmCode, setConfirmCode] = useState<string | undefined>(undefined);

    const fetchCoupons = async () => {
        try {
            const response = await axios.get(`${API}/admin/coupon`, { withCredentials: true });
            const list = response.data?.data?.coupons ?? response.data?.data ?? [];
            setCoupons(Array.isArray(list) ? list : []);
        } catch (error) {
            console.error('Failed to load coupons:', error);
            setCoupons([]);
        }
    };

    const isSubmitDisabled = Boolean(submitting || loading);

    useEffect(() => {
        const loadData = async () => {
            setLoading(true);
            await fetchCoupons();
            setLoading(false);
        };

        loadData();
    }, []);

    const resetForm = () => {
        setEditingId(null);
        setCode('');
        setDiscount(0);
        setMinPurchase(0);
        setStartAt('');
        setEndAt('');
        setIsActive(false);
    };

    const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
        event.preventDefault();

        const trimmedCode = code.trim();
        if (!trimmedCode) {
            toast.error('Please enter a coupon code');
            return;
        }

        if (Number(discount) < 0 || Number(minPurchase) < 0) {
            toast.error('Discount and minimum purchase cannot be negative');
            return;
        }

        if (!startAt || !endAt) {
            toast.error('Please provide both start and end dates');
            return;
        }

        if (new Date(endAt) < new Date(startAt)) {
            toast.error('End date must be on or after the start date');
            return;
        }

        setSubmitting(true);

        try {
            const payload = {
                code: trimmedCode,
                discount: Number(discount),
                min_purchase: Number(minPurchase),
                start_at: startAt,
                end_at: endAt,
                is_active: isActive,
            };

            if (editingId) {
                await axios.put(`${API}/admin/coupon/${editingId}`, payload, { withCredentials: true });
                toast.success('Coupon updated successfully');
            } else {
                await axios.post(`${API}/admin/coupon`, payload, { withCredentials: true });
                toast.success('Coupon created successfully');
            }

            resetForm();
            await fetchCoupons();
        } catch (error: any) {
            console.error('Coupon save failed:', error);
            const message = error?.response?.data?.message || 'Failed to save coupon';
            toast.error(message);
        } finally {
            setSubmitting(false);
        }
    };

    const handleEdit = (coupon: CouponRecord) => {
        setEditingId(coupon.id);
        setCode(coupon.code);
        setDiscount(Number(coupon.discount || 0));
        setMinPurchase(Number(coupon.min_purchase || 0));
        setStartAt(normalizeDate(coupon.start_at));
        setEndAt(normalizeDate(coupon.end_at));
        setIsActive(Boolean(coupon.is_active));
    };

    const handleToggleStatus = async (id: number, active: boolean) => {
        try {
            const url = active ? `${API}/admin/coupon/${id}/inactivate` : `${API}/admin/coupon/${id}/activate`;
            await axios.put(url, {}, { withCredentials: true });
            await fetchCoupons();
            toast.success(active ? 'Coupon deactivated successfully' : 'Coupon activated successfully');
        } catch (error: any) {
            console.error('Coupon status update failed:', error);
            const message = error?.response?.data?.message || 'Failed to update coupon status';
            toast.error(message);
        }
    };

    const requestDelete = (id: number, code?: string) => {
        setConfirmId(id);
        setConfirmCode(code);
        setConfirmOpen(true);
    };

    const handleDelete = async () => {
        if (!confirmId) return;

        try {
            await axios.delete(`${API}/admin/coupon/${confirmId}`, { withCredentials: true });
            setCoupons((prev) => prev.filter((coupon) => coupon.id !== confirmId));
            toast.success('Coupon deleted successfully');
        } catch (error) {
            console.error('Failed to delete coupon:', error);
            toast.error('Failed to delete coupon');
        } finally {
            setConfirmOpen(false);
            setConfirmId(null);
            setConfirmCode(undefined);
        }
    };

    return (
        <div className="min-h-screen bg-slate-950 px-4 py-8 text-white sm:px-6 lg:px-8">
            <div className="mx-auto max-w-7xl">
                <header className="mb-6 flex flex-col gap-4 rounded-3xl border border-white/10 bg-white/5 p-5 shadow-2xl shadow-slate-950/30 backdrop-blur-xl sm:flex-row sm:items-center sm:justify-between">
                    <div className="flex items-center gap-4">
                        <Link href="/admin" className="rounded-full border border-white/10 bg-slate-900/60 px-3 py-2 text-sm font-medium text-slate-200 hover:bg-slate-900 transition">← Back</Link>
                        <div>
                            <h1 className="text-2xl font-bold">Coupon Management</h1>
                            <p className="text-sm text-slate-400">Create, update, and activate product coupons</p>
                        </div>
                    </div>
                </header>

                <main className="grid gap-6 lg:grid-cols-[420px_minmax(0,1fr)]">
                    <section className="rounded-3xl border border-white/10 bg-slate-900/60 p-6 shadow-2xl shadow-slate-950/30">
                        <div className="mb-6 flex items-center justify-between">
                            <h2 className="text-xl font-semibold">{editingId ? 'Edit Coupon' : 'Create Coupon'}</h2>
                            {editingId && (
                                <button type="button" onClick={resetForm} className="text-sm text-slate-300 hover:text-white">
                                    Cancel
                                </button>
                            )}
                        </div>

                        <form onSubmit={handleSubmit} className="space-y-5">
                            <div className="space-y-2">
                                <label htmlFor="code" className="text-sm font-medium text-slate-200">Coupon Code</label>
                                <input
                                    id="code"
                                    value={code}
                                    onChange={(e) => setCode(e.target.value.toUpperCase())}
                                    className="w-full rounded-xl border border-white/10 bg-slate-950 px-3 py-3 text-white outline-none"
                                    placeholder="SAVE10"
                                />
                            </div>

                            <div className="grid gap-4 sm:grid-cols-2">
                                <div className="space-y-2">
                                    <label htmlFor="discount" className="text-sm font-medium text-slate-200">Discount</label>
                                    <input
                                        id="discount"
                                        type="number"
                                        min={0}
                                        step="0.01"
                                        value={discount || ''}
                                        onChange={(e) => setDiscount(Number(e.target.value))}
                                        className="w-full rounded-xl border border-white/10 bg-slate-950 px-3 py-3 text-white outline-none"
                                        placeholder="10"
                                    />
                                </div>

                                <div className="space-y-2">
                                    <label htmlFor="minPurchase" className="text-sm font-medium text-slate-200">Min Purchase</label>
                                    <input
                                        id="minPurchase"
                                        type="number"
                                        min={0}
                                        step="0.01"
                                        value={minPurchase || ''}
                                        onChange={(e) => setMinPurchase(Number(e.target.value))}
                                        className="w-full rounded-xl border border-white/10 bg-slate-950 px-3 py-3 text-white outline-none"
                                        placeholder="500"
                                    />
                                </div>
                            </div>

                            <div className="grid gap-4 sm:grid-cols-2">
                                <div className="space-y-2">
                                    <label htmlFor="startAt" className="text-sm font-medium text-slate-200">Start At</label>
                                    <input
                                        id="startAt"
                                        type="date"
                                        value={startAt}
                                        onChange={(e) => setStartAt(e.target.value)}
                                        className="w-full rounded-xl border border-white/10 bg-slate-950 px-3 py-3 text-white outline-none"
                                    />
                                </div>

                                <div className="space-y-2">
                                    <label htmlFor="endAt" className="text-sm font-medium text-slate-200">End At</label>
                                    <input
                                        id="endAt"
                                        type="date"
                                        value={endAt}
                                        onChange={(e) => setEndAt(e.target.value)}
                                        className="w-full rounded-xl border border-white/10 bg-slate-950 px-3 py-3 text-white outline-none"
                                    />
                                </div>
                            </div>

                            <label className="flex items-center gap-3 rounded-xl border border-white/10 bg-slate-950/60 px-3 py-3 text-sm text-slate-200">
                                <input
                                    type="checkbox"
                                    checked={isActive}
                                    onChange={(e) => setIsActive(e.target.checked)}
                                    className="h-4 w-4 rounded border-white/10 bg-slate-900/50"
                                />
                                Active immediately
                            </label>

                            <button
                                type="submit"
                                disabled={isSubmitDisabled}
                                className="w-full rounded-xl bg-gradient-to-r from-pink-500 to-rose-500 px-4 py-3 font-semibold text-white disabled:cursor-not-allowed disabled:opacity-60"
                            >
                                {submitting ? (editingId ? 'Updating Coupon...' : 'Creating Coupon...') : (editingId ? 'Update Coupon' : 'Create Coupon')}
                            </button>
                        </form>
                    </section>

                    <section className="rounded-3xl border border-white/10 bg-slate-900/60 p-6 shadow-2xl shadow-slate-950/30">
                        <div className="mb-5 flex items-center justify-between">
                            <h2 className="text-xl font-semibold">Current Coupons</h2>
                            <span className="rounded-full border border-pink-500/30 bg-pink-500/10 px-3 py-1 text-xs font-medium text-pink-200">
                                {coupons.length} total
                            </span>
                        </div>

                        {loading ? (
                            <div className="rounded-xl border border-white/10 bg-slate-950/40 p-10 text-center text-slate-300">
                                Loading coupons...
                            </div>
                        ) : coupons.length === 0 ? (
                            <div className="rounded-xl border border-dashed border-white/10 bg-slate-950/30 p-10 text-center text-slate-300">
                                No coupons available yet.
                            </div>
                        ) : (
                            <div className="grid gap-4 xl:grid-cols-2">
                                {coupons.map((coupon) => (
                                    <div key={coupon.id} className="overflow-hidden rounded-2xl border border-white/10 bg-slate-950/50">
                                        <div className="border-b border-white/10 p-4">
                                            <div className="flex items-center justify-between gap-3">
                                                <div>
                                                    <p className="text-[10px] font-semibold uppercase tracking-[0.2em] text-slate-400">Coupon</p>
                                                    <h3 className="mt-1 text-lg font-semibold text-white">{coupon.code}</h3>
                                                </div>
                                                <span className={`rounded-full px-2.5 py-1 text-xs font-medium ${coupon.is_active ? 'bg-emerald-500/10 text-emerald-300 border border-emerald-500/20' : 'bg-slate-700 text-slate-300 border border-white/10'}`}>
                                                    {coupon.is_active ? 'Active' : 'Inactive'}
                                                </span>
                                            </div>
                                        </div>

                                        <div className="space-y-3 p-4 text-sm text-slate-300">
                                            <div className="flex items-center justify-between gap-3 rounded-xl border border-pink-500/20 bg-pink-500/5 px-3 py-2">
                                                <span>Discount</span>
                                                <span className="text-lg font-bold text-pink-300">₹{Number(coupon.discount || 0).toFixed(2)}</span>
                                            </div>

                                            <div className="grid gap-2 rounded-xl border border-white/10 bg-slate-900/60 p-3 text-xs text-slate-300">
                                                <div className="flex items-center justify-between gap-3">
                                                    <span>Min Purchase</span>
                                                    <span className="font-medium text-white">₹{Number(coupon.min_purchase || 0).toFixed(2)}</span>
                                                </div>
                                                <div className="flex items-center justify-between gap-3">
                                                    <span>Start</span>
                                                    <span className="font-medium text-white">{normalizeDate(coupon.start_at) || '—'}</span>
                                                </div>
                                                <div className="flex items-center justify-between gap-3">
                                                    <span>End</span>
                                                    <span className="font-medium text-white">{normalizeDate(coupon.end_at) || '—'}</span>
                                                </div>
                                            </div>

                                            <div className="flex gap-2 pt-2">
                                                <button
                                                    type="button"
                                                    onClick={() => handleEdit(coupon)}
                                                    className="flex-1 rounded-xl border border-cyan-500/30 bg-cyan-500/10 px-3 py-2 font-medium text-cyan-200 hover:bg-cyan-500/15"
                                                >
                                                    Edit
                                                </button>
                                                <button
                                                    type="button"
                                                    onClick={() => handleToggleStatus(coupon.id, coupon.is_active)}
                                                    className={`flex-1 rounded-xl border px-3 py-2 font-medium ${coupon.is_active ? 'border-amber-500/30 bg-amber-500/10 text-amber-200 hover:bg-amber-500/15' : 'border-emerald-500/30 bg-emerald-500/10 text-emerald-200 hover:bg-emerald-500/15'}`}
                                                >
                                                    {coupon.is_active ? 'Deactivate' : 'Activate'}
                                                </button>
                                                <button
                                                    type="button"
                                                    onClick={() => requestDelete(coupon.id, coupon.code)}
                                                    className="flex-1 rounded-xl border border-red-500/30 bg-red-500/10 px-3 py-2 font-medium text-red-200 hover:bg-red-500/15"
                                                >
                                                    Delete
                                                </button>
                                            </div>
                                        </div>
                                    </div>
                                ))}
                            </div>
                        )}
                    </section>
                </main>
            </div>

            <ConfirmModal
                open={confirmOpen}
                title={`Delete coupon${confirmCode ? `: ${confirmCode}` : ''}`}
                description="This removes the coupon from the admin list and prevents it from being used."
                confirmLabel="Delete"
                cancelLabel="Cancel"
                loading={false}
                onCancel={() => setConfirmOpen(false)}
                onConfirm={handleDelete}
            />
        </div>
    );
}
