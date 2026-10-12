'use client';

import axios from 'axios';
import Link from 'next/link';
import { useEffect, useState } from 'react';

const API = 'http://localhost:8080';

interface WalletTransaction {
  id: number;
  amount: number;
  transaction_type: string;
  description: string;
  created_at: string;
}

export default function WalletPage() {
  const [balance, setBalance] = useState(0);
  const [transactions, setTransactions] = useState<WalletTransaction[]>([]);
  const [loading, setLoading] = useState(true);
  const [addingMoney, setAddingMoney] = useState(false);

  const formatTransactionDate = (value?: string) => {
    if (!value) return 'Unknown date';

    const parsed = new Date(value);
    if (Number.isNaN(parsed.getTime())) return 'Unknown date';

    return parsed.toLocaleString();
  };

  const loadWallet = async () => {
    try {
      const [walletResponse, txResponse] = await Promise.all([
        axios.get(`${API}/wallet`, { withCredentials: true }),
        axios.get(`${API}/wallet/history`, { withCredentials: true }),
      ]);

      const walletData = walletResponse.data?.data ?? walletResponse.data ?? {};
      const txData = txResponse.data?.data ?? txResponse.data ?? [];

      setBalance(Number(walletData.balance ?? 0));
      setTransactions(Array.isArray(txData) ? txData : []);
    } catch (error) {
      console.error('Failed to load wallet', error);
      setBalance(0);
      setTransactions([]);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    void loadWallet();
  }, []);

  const handleAddMoney = async () => {
    const amountInput = window.prompt('Enter amount to add to wallet (INR)', '100');
    if (amountInput === null) return;

    const amount = Number(amountInput);
    if (!Number.isFinite(amount) || amount <= 0) {
      window.alert('Please enter a valid number greater than zero.');
      return;
    }

    try {
      setAddingMoney(true);
      const response = await axios.post(
        `${API}/wallet/add-money`,
        { amount },
        { withCredentials: true }
      );

      const nextBalance = Number(response.data?.data?.balance ?? balance ?? 0);
      setBalance(nextBalance);
      await loadWallet();
      window.alert(response.data?.message || 'Money added to wallet successfully');
    } catch (error: any) {
      console.error('Failed to add money', error);
      window.alert(error.response?.data?.message || 'Unable to add money to wallet');
    } finally {
      setAddingMoney(false);
    }
  };

  return (
    <div className="min-h-screen bg-slate-950 text-white">
      <header className="sticky top-0 z-20 border-b border-white/10 bg-slate-950/80 backdrop-blur-xl">
        <div className="mx-auto flex max-w-7xl items-center justify-between px-6 py-4">
          <div className="flex items-center gap-4">
            <Link href="/dashboard" className="text-lg font-semibold text-white hover:text-cyan-400">
              GenZone
            </Link>
          </div>

          <div className="flex items-center gap-3">
            <Link href="/dashboard" className="rounded-xl border border-white/10 bg-white/5 px-4 py-2 text-sm font-medium text-slate-100 hover:border-cyan-400/40 hover:bg-cyan-500/10">
              Back to Dashboard
            </Link>
          </div>
        </div>
      </header>

      <main className="mx-auto max-w-5xl px-6 py-10">
        <div className="mb-8">
          <p className="text-sm uppercase tracking-[0.25em] text-amber-300">Wallet</p>
          <h1 className="mt-2 text-3xl font-bold">My Wallet</h1>
        </div>

        <section className="mb-8 rounded-3xl border border-amber-400/30 bg-gradient-to-br from-amber-500/15 via-slate-900 to-slate-900 p-6 shadow-lg shadow-amber-950/10">
          <div className="flex items-center justify-between gap-4">
            <div>
              <p className="text-sm text-slate-300">Current balance</p>
              <p className="mt-3 text-4xl font-bold text-amber-300">₹{balance.toFixed(2)}</p>
            </div>
            <div className="flex flex-col items-end gap-3">
              <div className="rounded-2xl border border-amber-300/30 bg-amber-500/10 px-4 py-3 text-right">
                <p className="text-xs uppercase tracking-[0.2em] text-amber-200">Status</p>
                <p className="mt-2 text-sm font-semibold text-amber-100">Active</p>
              </div>

              <button
                type="button"
                onClick={handleAddMoney}
                disabled={addingMoney}
                className="rounded-xl bg-gradient-to-r from-amber-400 to-yellow-500 px-5 py-3 text-sm font-semibold text-slate-950 shadow-lg shadow-amber-500/20 transition hover:brightness-110 disabled:cursor-not-allowed disabled:opacity-60"
              >
                {addingMoney ? 'Adding...' : 'Add Money'}
              </button>
            </div>
          </div>
        </section>

        <section className="rounded-3xl border border-white/10 bg-slate-900/80 p-6">
          <div className="mb-5 flex items-center justify-between gap-4">
            <h2 className="text-xl font-semibold text-white">Transactions</h2>
          </div>

          {loading ? (
            <div className="rounded-xl border border-white/10 bg-white/5 p-5 text-slate-300">Loading wallet history...</div>
          ) : transactions.length === 0 ? (
            <div className="rounded-xl border border-dashed border-slate-600 bg-white/5 p-8 text-center text-slate-400">
              No wallet transactions yet.
            </div>
          ) : (
            <div className="space-y-3">
              {transactions.map((transaction) => {
                const isCredit = String(transaction.transaction_type).toLowerCase() === 'credit';
                const sign = isCredit ? '+' : '-';

                return (
                  <div key={transaction.id} className="flex items-center justify-between gap-3 rounded-2xl border border-white/10 bg-slate-950/60 p-4">
                    <div>
                      <p className="font-medium text-white">{transaction.description}</p>
                      <p className="mt-1 text-xs text-slate-400">
                        {formatTransactionDate(transaction.created_at)}
                      </p>
                    </div>

                    <div className="text-right">
                      <p className={`text-lg font-semibold ${isCredit ? 'text-emerald-300' : 'text-red-300'}`}>
                        {sign}₹{Math.abs(Number(transaction.amount)).toFixed(2)}
                      </p>
                      <p className="text-xs uppercase tracking-[0.2em] text-slate-400">
                        {transaction.transaction_type}
                      </p>
                    </div>
                  </div>
                );
              })}
            </div>
          )}
        </section>
      </main>
    </div>
  );
}
