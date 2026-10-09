'use client';

import axios from 'axios';
import Link from 'next/link';
import { useEffect, useState } from 'react';
import { toast } from 'sonner';
import ConfirmModal from '@/components/admin/confirm-modal';

const API = 'http://localhost:8080';

type ProductOption = {
    id: number;
    product_name: string;
    category_name?: string;
    price?: number;
};

type OfferRecord = {
    id: number;
    product_id: number;
    product_name: string;
    category_name?: string;
    description?: string;
    image_url?: string;
    price?: number;
    stock?: number;
    popular?: boolean;
    size?: string;
    discount_percentage: number;
    start_at?: string;
    end_at?: string;
};

const formatDateInput = (value?: string) => {
    if (!value) return '';
    const trimmed = value.trim();
    if (!trimmed) return '';
    if (/^\d{4}-\d{2}-\d{2}$/.test(trimmed)) return trimmed;
    const date = new Date(trimmed);
    if (Number.isNaN(date.getTime())) return '';
    return date.toISOString().slice(0, 10);
};

export default function OffersPage() {
    const [offers, setOffers] = useState<OfferRecord[]>([]);
    const [products, setProducts] = useState<ProductOption[]>([]);
    const [selectedProductId, setSelectedProductId] = useState('');
    const [discountPercentage, setDiscountPercentage] = useState(0);
    const [startAt, setStartAt] = useState('');
    const [endAt, setEndAt] = useState('');
    const [editingId, setEditingId] = useState<number | null>(null);
    const [loading, setLoading] = useState(true);
    const [submitting, setSubmitting] = useState(false);
    const [confirmOpen, setConfirmOpen] = useState(false);
    const [confirmId, setConfirmId] = useState<number | null>(null);
    const [confirmName, setConfirmName] = useState<string | undefined>(undefined);

    const fetchProducts = async () => {
        try {
            const response = await axios.get(`${API}/admin/product`, { withCredentials: true });
            const list = response.data?.data?.products ?? [];
            setProducts(Array.isArray(list) ? list : []);
        } catch (error) {
            console.error('Failed to load products:', error);
            setProducts([]);
        }
    };

    const fetchOffers = async () => {
        try {
            const response = await axios.get(`${API}/admin/offer`, { withCredentials: true });
            const list = response.data?.data?.offers ?? response.data?.data ?? [];
            setOffers(Array.isArray(list) ? list : []);
        } catch (error) {
            console.error('Failed to load offers:', error);
            setOffers([]);
        }
    };

    useEffect(() => {
        const loadData = async () => {
            setLoading(true);
            await Promise.all([fetchProducts(), fetchOffers()]);
            setLoading(false);
        };

        loadData();
    }, []);

    const resetForm = () => {
        setEditingId(null);
        setSelectedProductId('');
        setDiscountPercentage(0);
        setStartAt('');
        setEndAt('');
    };

    const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
        event.preventDefault();

        if (!selectedProductId) {
            toast.error('Please select a product');
            return;
        }

        if (Number(discountPercentage) < 1 || Number(discountPercentage) > 100) {
            toast.error('Discount percentage must be between 1 and 100');
            return;
        }

        if (!startAt || !endAt) {
            toast.error('Please provide both start and end dates');
            return;
        }

        if (new Date(endAt) < new Date(startAt)) {
            toast.error('End date must be the same as or later than the start date');
            return;
        }

        setSubmitting(true);

        try {
            const payload = {
                product_id: Number(selectedProductId),
                discount_percentage: Number(discountPercentage),
                start_at: startAt,
                end_at: endAt,
            };

            if (editingId) {
                await axios.put(`${API}/admin/offer/${editingId}`, payload, { withCredentials: true });
                toast.success('Offer updated successfully');
            } else {
                await axios.post(`${API}/admin/offer`, payload, { withCredentials: true });
                toast.success('Offer created successfully');
            }

            resetForm();
            await Promise.all([fetchProducts(), fetchOffers()]);
        } catch (error: any) {
            console.error('Offer save failed:', error);
            const message = error?.response?.data?.message || 'Failed to save offer';
            toast.error(message);
        } finally {
            setSubmitting(false);
        }
    };

    const handleEdit = (offer: OfferRecord) => {
        setEditingId(offer.id);
        setSelectedProductId(String(offer.product_id));
        setDiscountPercentage(Number(offer.discount_percentage || 0));
        setStartAt(formatDateInput(offer.start_at));
        setEndAt(formatDateInput(offer.end_at));
    };

    const requestDelete = (id: number, name?: string) => {
        setConfirmId(id);
        setConfirmName(name);
        setConfirmOpen(true);
    };

    const handleDelete = async () => {
        if (!confirmId) return;

        try {
            await axios.delete(`${API}/admin/offer/${confirmId}`, { withCredentials: true });
            setOffers((prev) => prev.filter((offer) => offer.id !== confirmId));
            toast.success('Offer deleted successfully');
        } catch (error) {
            console.error('Failed to delete offer:', error);
            toast.error('Failed to delete offer');
        } finally {
            setConfirmOpen(false);
            setConfirmId(null);
            setConfirmName(undefined);
        }
    };

    return (
        <div className="min-h-screen bg-slate-950 px-4 py-8 text-white sm:px-6 lg:px-8">
            <div className="mx-auto max-w-7xl">
                <header className="mb-6 flex flex-col gap-4 rounded-3xl border border-white/10 bg-white/5 p-5 shadow-2xl shadow-slate-950/30 backdrop-blur-xl sm:flex-row sm:items-center sm:justify-between">
                    <div className="flex items-center gap-4">
                        <Link href="/admin" className="rounded-full border border-white/10 bg-slate-900/60 px-3 py-2 text-sm font-medium text-slate-200 hover:bg-slate-900 transition">← Back</Link>
                        <div>
                            <h1 className="text-2xl font-bold">Offer Management</h1>
                            <p className="text-sm text-slate-400">Create and update product discount offers</p>
                        </div>
                    </div>
                </header>

                <main className="grid gap-6 lg:grid-cols-[420px_minmax(0,1fr)]">
                    <section className="rounded-3xl border border-white/10 bg-slate-900/60 p-6 shadow-2xl shadow-slate-950/30">
                        <div className="mb-6 flex items-center justify-between">
                            <h2 className="text-xl font-semibold">{editingId ? 'Edit Offer' : 'Create Offer'}</h2>
                            {editingId && (
                                <button
                                    type="button"
                                    onClick={resetForm}
                                    className="text-sm text-slate-300 hover:text-white"
                                >
                                    Cancel
                                </button>
                            )}
                        </div>

                        <form onSubmit={handleSubmit} className="space-y-5">
                            <div className="space-y-2">
                                <label htmlFor="product" className="text-sm font-medium text-slate-200">Product</label>
                                <select
                                    id="product"
                                    value={selectedProductId}
                                    onChange={(e) => setSelectedProductId(e.target.value)}
                                    className="w-full rounded-xl border border-white/10 bg-slate-950 px-3 py-3 text-white outline-none"
                                >
                                    <option value="">Select a product</option>
                                    {products.map((product) => (
                                        <option key={product.id} value={product.id}>
                                            {product.product_name}
                                            {product.category_name ? ` · ${product.category_name}` : ''}
                                        </option>
                                    ))}
                                </select>
                            </div>

                            <div className="space-y-2">
                                <label htmlFor="discount" className="text-sm font-medium text-slate-200">Discount Percentage (%)</label>
                                <input
                                    id="discount"
                                    type="number"
                                    min={1}
                                    max={100}
                                    value={discountPercentage || ''}
                                    onChange={(e) => setDiscountPercentage(Number(e.target.value))}
                                    className="w-full rounded-xl border border-white/10 bg-slate-950 px-3 py-3 text-white outline-none"
                                    placeholder="e.g. 15"
                                />
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

                            <button
                                type="submit"
                                disabled={submitting || loading}
                                className="w-full rounded-xl bg-gradient-to-r from-indigo-500 to-blue-500 px-4 py-3 font-semibold text-white disabled:cursor-not-allowed disabled:opacity-60"
                            >
                                {submitting ? (editingId ? 'Updating Offer...' : 'Creating Offer...') : (editingId ? 'Update Offer' : 'Create Offer')}
                            </button>
                        </form>
                    </section>

                    <section className="rounded-3xl border border-white/10 bg-slate-900/60 p-6 shadow-2xl shadow-slate-950/30">
                        <div className="mb-5 flex items-center justify-between">
                            <h2 className="text-xl font-semibold">Current Offers</h2>
                            <span className="rounded-full border border-indigo-500/30 bg-indigo-500/10 px-3 py-1 text-xs font-medium text-indigo-200">
                                {offers.length} active
                            </span>
                        </div>

                        {loading ? (
                            <div className="rounded-xl border border-white/10 bg-slate-950/40 p-10 text-center text-slate-300">
                                Loading offers...
                            </div>
                        ) : offers.length === 0 ? (
                            <div className="rounded-xl border border-dashed border-white/10 bg-slate-950/30 p-10 text-center text-slate-300">
                                No offers available yet.
                            </div>
                        ) : (
                            <div className="grid gap-4 xl:grid-cols-2">
                                {offers.map((offer) => (
                                    <div key={offer.id} className="overflow-hidden rounded-2xl border border-white/10 bg-slate-950/50">
                                        <div className="flex gap-3 border-b border-white/10 p-4">
                                            <div className="h-16 w-16 overflow-hidden rounded-xl bg-slate-800">
                                                {offer.image_url ? (
                                                    <img src={offer.image_url} alt={offer.product_name} className="h-full w-full object-cover" />
                                                ) : (
                                                    <div className="flex h-full items-center justify-center text-xs text-slate-400">IMG</div>
                                                )}
                                            </div>
                                            <div className="min-w-0 flex-1">
                                                <p className="text-[10px] font-semibold uppercase tracking-[0.2em] text-slate-400">{offer.category_name || 'Category'}</p>
                                                <h3 className="mt-1 truncate text-lg font-semibold text-white">{offer.product_name}</h3>
                                                <p className="mt-1 text-sm text-slate-300">₹{Number(offer.price ?? 0).toFixed(2)}</p>
                                            </div>
                                        </div>

                                        <div className="space-y-3 p-4 text-sm text-slate-300">
                                            <div className="flex items-center justify-between gap-3 rounded-xl border border-indigo-500/20 bg-indigo-500/5 px-3 py-2">
                                                <span className="text-slate-300">Discount</span>
                                                <span className="text-lg font-bold text-indigo-300">{Number(offer.discount_percentage || 0)}%</span>
                                            </div>

                                            <div className="grid gap-2 rounded-xl border border-white/10 bg-slate-900/60 p-3 text-xs text-slate-300">
                                                <div className="flex items-center justify-between gap-3">
                                                    <span>Start</span>
                                                    <span className="font-medium text-white">{formatDateInput(offer.start_at) || '—'}</span>
                                                </div>
                                                <div className="flex items-center justify-between gap-3">
                                                    <span>End</span>
                                                    <span className="font-medium text-white">{formatDateInput(offer.end_at) || '—'}</span>
                                                </div>
                                            </div>

                                            <div className="flex flex-wrap gap-2 text-xs text-slate-400">
                                                {offer.size && <span className="rounded-full border border-white/10 bg-white/5 px-2 py-1">Size: {offer.size}</span>}
                                                {typeof offer.stock === 'number' && (
                                                    <span className="rounded-full border border-white/10 bg-white/5 px-2 py-1">Stock: {offer.stock}</span>
                                                )}
                                                {offer.popular && <span className="rounded-full border border-amber-500/20 bg-amber-500/10 px-2 py-1 text-amber-200">Popular</span>}
                                            </div>

                                            <div className="flex gap-2 pt-2">
                                                <button
                                                    type="button"
                                                    onClick={() => handleEdit(offer)}
                                                    className="flex-1 rounded-xl border border-cyan-500/30 bg-cyan-500/10 px-3 py-2 font-medium text-cyan-200 hover:bg-cyan-500/15"
                                                >
                                                    Edit
                                                </button>
                                                <button
                                                    type="button"
                                                    onClick={() => requestDelete(offer.id, offer.product_name)}
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
                title={`Delete offer${confirmName ? `: ${confirmName}` : ''}`}
                description="This will remove the active discount from the selected product."
                confirmLabel="Delete"
                cancelLabel="Cancel"
                loading={false}
                onCancel={() => setConfirmOpen(false)}
                onConfirm={handleDelete}
            />
        </div>
    );
}
