'use client';

import axios from 'axios';
import Link from 'next/link';
import { useEffect, useState } from 'react';
import { toast } from 'sonner';

const API = 'http://localhost:8080';

type UserRecord = {
    id: number;
    first_name: string;
    last_name: string;
    email: string;
    phone: string;
    status?: string;
    profile_picture?: string | null;
};

const getDisplayName = (user: UserRecord) => {
    const first = user.first_name?.trim();
    const last = user.last_name?.trim();

    if (first || last) {
        return `${first || ''} ${last || ''}`.trim();
    }

    return 'Unknown User';
};

const getStatusStyles = (status?: string) => {
    const normalized = (status || 'Active').toLowerCase();

    if (normalized === 'blocked') {
        return 'border-red-400/30 bg-red-500/10 text-red-200';
    }

    return 'border-emerald-400/30 bg-emerald-500/10 text-emerald-200';
};

export default function UsersPage() {
    const [users, setUsers] = useState<UserRecord[]>([]);
    const [loading, setLoading] = useState(false);
    const [processingId, setProcessingId] = useState<number | null>(null);

    const fetchUsers = async () => {
        try {
            const response = await axios.get(`${API}/admin/listusers`, { withCredentials: true });
            const list = response.data?.data?.users ?? [];
            setUsers(Array.isArray(list) ? list : []);
        } catch (error) {
            console.error('Failed to load users:', error);
            setUsers([]);
            toast.error('Failed to load users');
        }
    };

    useEffect(() => {
        const loadUsers = async () => {
            setLoading(true);
            await fetchUsers();
            setLoading(false);
        };

        loadUsers();
    }, []);

    const handleStatusToggle = async (user: UserRecord) => {
        const isBlocked = (user.status || 'Active').toLowerCase() === 'blocked';
        const endpoint = isBlocked ? `${API}/admin/listusers/unblockuser` : `${API}/admin/listusers/blockuser`;

        try {
            setProcessingId(user.id);
            await axios.put(endpoint, { user_id: user.id }, { withCredentials: true });
            setUsers((prev) =>
                prev.map((item) =>
                    item.id === user.id
                        ? { ...item, status: isBlocked ? 'Active' : 'Blocked' }
                        : item
                )
            );
            toast.success(isBlocked ? 'User unblocked successfully' : 'User blocked successfully');
        } catch (error: any) {
            console.error('User status update failed:', error);
            const message = error?.response?.data?.message || 'Failed to update user status';
            toast.error(message);
        } finally {
            setProcessingId(null);
        }
    };

    const totalUsers = users.length;
    const activeUsers = users.filter((user) => (user.status || 'Active').toLowerCase() === 'active').length;
    const blockedUsers = users.filter((user) => (user.status || 'Active').toLowerCase() === 'blocked').length;

    return (
        <div className="min-h-screen bg-slate-950 px-4 py-8 text-white sm:px-6 lg:px-8">
            <div className="mx-auto max-w-7xl">
                <header className="mb-6 flex flex-col gap-4 rounded-3xl border border-white/10 bg-white/5 p-5 shadow-2xl shadow-slate-950/30 backdrop-blur-xl sm:flex-row sm:items-center sm:justify-between">
                    <div className="flex items-center gap-4">
                        <Link href="/admin" className="rounded-full border border-white/10 bg-slate-900/60 px-3 py-2 text-sm font-medium text-slate-200 hover:bg-slate-900 transition">← Back</Link>
                        <div>
                            <h1 className="text-2xl font-bold">Users Management</h1>
                            <p className="text-sm text-slate-400">Review customer accounts and manage access</p>
                        </div>
                    </div>
                </header>

                <section className="mb-6 grid gap-4 md:grid-cols-3">
                    <div className="rounded-2xl border border-white/10 bg-slate-900/80 p-5 shadow-lg shadow-slate-950/25">
                        <p className="text-sm text-slate-400">Total Users</p>
                        <p className="mt-3 text-3xl font-bold text-white">{totalUsers}</p>
                    </div>
                    <div className="rounded-2xl border border-emerald-400/30 bg-emerald-500/10 p-5 shadow-lg shadow-slate-950/25">
                        <p className="text-sm text-emerald-200">Active</p>
                        <p className="mt-3 text-3xl font-bold text-white">{activeUsers}</p>
                    </div>
                    <div className="rounded-2xl border border-red-400/30 bg-red-500/10 p-5 shadow-lg shadow-slate-950/25">
                        <p className="text-sm text-red-200">Blocked</p>
                        <p className="mt-3 text-3xl font-bold text-white">{blockedUsers}</p>
                    </div>
                </section>

                <main className="rounded-3xl border border-white/10 bg-slate-900/60 p-4 shadow-2xl shadow-slate-950/30 sm:p-6">
                    {loading ? (
                        <div className="rounded-xl border border-white/10 bg-slate-900/50 p-8 text-center text-slate-300">
                            Loading users...
                        </div>
                    ) : users.length === 0 ? (
                        <div className="rounded-xl border border-white/10 bg-slate-900/50 p-12 text-center text-slate-300">
                            No users found.
                        </div>
                    ) : (
                        <div className="overflow-x-auto">
                            <table className="min-w-full border-separate border-spacing-y-3 text-left text-sm">
                                <thead>
                                    <tr className="text-slate-400">
                                        <th className="px-4 py-3 font-medium">User</th>
                                        <th className="px-4 py-3 font-medium">Email</th>
                                        <th className="px-4 py-3 font-medium">Phone</th>
                                        <th className="px-4 py-3 font-medium">Status</th>
                                        <th className="px-4 py-3 font-medium text-right">Action</th>
                                    </tr>
                                </thead>
                                <tbody>
                                    {users.map((user) => {
                                        const isBlocked = (user.status || 'Active').toLowerCase() === 'blocked';
                                        return (
                                            <tr key={user.id} className="rounded-2xl bg-slate-900/80 text-slate-200 shadow-lg shadow-slate-950/15">
                                                <td className="rounded-l-2xl px-4 py-4">
                                                    <div className="flex items-center gap-3">
                                                        <div className="flex h-10 w-10 items-center justify-center overflow-hidden rounded-full bg-gradient-to-br from-cyan-500 to-sky-500 font-semibold text-white">
                                                            {user.profile_picture ? (
                                                                <img
                                                                    src={user.profile_picture}
                                                                    alt={getDisplayName(user)}
                                                                    className="h-full w-full object-cover"
                                                                />
                                                            ) : (
                                                                <span>{(getDisplayName(user).charAt(0) || 'U').toUpperCase()}</span>
                                                            )}
                                                        </div>
                                                        <div>
                                                            <p className="font-medium text-white">{getDisplayName(user)}</p>
                                                            <p className="text-xs text-slate-400">ID #{user.id}</p>
                                                        </div>
                                                    </div>
                                                </td>
                                                <td className="px-4 py-4 text-slate-300">{user.email || '—'}</td>
                                                <td className="px-4 py-4 text-slate-300">{user.phone || '—'}</td>
                                                <td className="px-4 py-4">
                                                    <span className={`inline-flex rounded-full border px-2.5 py-1 text-xs font-medium ${getStatusStyles(user.status)}`}>
                                                        {user.status || 'Active'}
                                                    </span>
                                                </td>
                                                <td className="rounded-r-2xl px-4 py-4 text-right">
                                                    <button
                                                        type="button"
                                                        onClick={() => handleStatusToggle(user)}
                                                        disabled={processingId === user.id}
                                                        className={`rounded-xl px-3 py-2 text-xs font-semibold transition ${
                                                            isBlocked
                                                                ? 'border border-emerald-400/30 bg-emerald-500/10 text-emerald-200 hover:bg-emerald-500/20'
                                                                : 'border border-red-400/30 bg-red-500/10 text-red-200 hover:bg-red-500/20'
                                                        } disabled:cursor-not-allowed disabled:opacity-60`}
                                                    >
                                                        {processingId === user.id ? 'Updating...' : isBlocked ? 'Unblock' : 'Block'}
                                                    </button>
                                                </td>
                                            </tr>
                                        );
                                    })}
                                </tbody>
                            </table>
                        </div>
                    )}
                </main>
            </div>
        </div>
    );
}
