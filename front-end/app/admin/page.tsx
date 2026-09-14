'use client';

import axios from 'axios';
import Link from 'next/link';

const dashboardCards = [
    { title: 'Users', href: '/admin/users', emoji: '👥', accent: 'from-cyan-500 to-sky-500' },
    { title: 'Categories', href: '/admin/categories', emoji: '🗂️', accent: 'from-violet-500 to-purple-500' },
    { title: 'Products', href: '/admin/products', emoji: '📦', accent: 'from-emerald-500 to-teal-500' },
    { title: 'Orders', href: '/admin/orders', emoji: '🧾', accent: 'from-amber-500 to-orange-500' },
    { title: 'Coupons', href: '/admin/coupons', emoji: '🎟️', accent: 'from-pink-500 to-rose-500' },
    { title: 'Offers', href: '/admin/offers', emoji: '🏷️', accent: 'from-indigo-500 to-blue-500' },
    { title: 'Sales Report', href: '/admin/salesreport', emoji: '📊', accent: 'from-green-500 to-lime-500' },
    { title: 'Best Selling', href: '/admin/bestselling', emoji: '⭐', accent: 'from-yellow-500 to-amber-500' },
    { title: 'Invoices', href: '/admin/invoice', emoji: '🧾', accent: 'from-fuchsia-500 to-pink-500' },
];

const stats = [
    { label: 'Total Users', value: '1,284', detail: '+12.5% this month' },
    { label: 'Products', value: '342', detail: '+24 new this week' },
    { label: 'Orders', value: '968', detail: '89 pending' },
    { label: 'Revenue', value: '$24.8k', detail: '+18.4% vs last month' },
];

export default function AdminDashboard() {
    async function handleLogout() {
        try {
            const response = await axios.post(
                'http://localhost:8080/admin/logout',
                {},
                { withCredentials: true }
            );

            console.log(response.data);
            window.location.href = '/admin/login';
        } catch (error) {
            console.error('Logout failed:', error);
        }
    }

    return (
        <div className="min-h-screen bg-slate-950 px-4 py-8 text-white sm:px-6 lg:px-8">
            <div className="mx-auto max-w-7xl">
                <header className="mb-8 flex flex-col gap-4 rounded-3xl border border-white/10 bg-white/5 p-5 shadow-2xl shadow-slate-950/30 backdrop-blur-xl sm:flex-row sm:items-center sm:justify-between">
                    <div>
                        <p className="text-xs font-semibold uppercase tracking-[0.3em] text-cyan-300">
                            Control Center
                        </p>
                        <h1 className="mt-2 text-3xl font-bold">Admin dashboard</h1>
                    </div>

                    <button
                        type="button"
                        onClick={handleLogout}
                        className="rounded-xl border border-red-500/40 bg-red-500/10 px-4 py-2 text-sm font-medium text-red-200 transition hover:bg-red-500/20"
                    >
                        Log out
                    </button>
                </header>

                <section className="mb-8 grid gap-4 md:grid-cols-2 xl:grid-cols-4">
                    {stats.map((stat) => (
                        <div
                            key={stat.label}
                            className="rounded-2xl border border-white/10 bg-slate-900/80 p-5 shadow-lg shadow-slate-950/25"
                        >
                            <p className="text-sm text-slate-400">{stat.label}</p>
                            <p className="mt-3 text-3xl font-bold text-white">{stat.value}</p>
                            <p className="mt-2 text-xs text-emerald-300">{stat.detail}</p>
                        </div>
                    ))}
                </section>

                <section className="rounded-3xl border border-white/10 bg-slate-900/60 p-6 shadow-2xl shadow-slate-950/30">
                    <div className="mb-6 flex items-center justify-between gap-3">
                        <h2 className="text-xl font-semibold text-white">Quick actions</h2>
                        <span className="rounded-full border border-cyan-400/30 bg-cyan-500/10 px-3 py-1 text-xs font-medium text-cyan-200">
                            Management
                        </span>
                    </div>

                    <div className="grid gap-4 sm:grid-cols-2 xl:grid-cols-3">
                        {dashboardCards.map((card) => (
                            <Link
                                key={card.title}
                                href={card.href}
                                className="group rounded-2xl border border-white/10 bg-gradient-to-br from-slate-800 to-slate-900 p-4 transition hover:-translate-y-1 hover:border-cyan-400/40"
                            >
                                <div className={`mb-4 flex h-12 w-12 items-center justify-center rounded-xl bg-gradient-to-br ${card.accent} text-2xl`}>
                                    {card.emoji}
                                </div>
                                <div className="flex items-center justify-between gap-3">
                                    <h3 className="text-lg font-semibold text-white">{card.title}</h3>
                                    <span className="text-xl text-slate-300 transition group-hover:translate-x-1 group-hover:text-cyan-300">
                                        →
                                    </span>
                                </div>
                            </Link>
                        ))}
                    </div>
                </section>
            </div>
        </div>
    );
}