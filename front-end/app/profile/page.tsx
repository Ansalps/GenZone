'use client';

import axios from 'axios';
import Link from 'next/link';
import { useEffect, useState } from 'react';

interface UserProfile {
  first_name?: string;
  last_name?: string;
  email?: string;
  phone?: string;
  status?: string;
}

const API = 'http://localhost:8080';

export default function ProfilePage() {
  const [profile, setProfile] = useState<UserProfile>({});
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchProfile = async () => {
      try {
        const response = await axios.get(`${API}/profile`, { withCredentials: true });
        const user = response.data?.data?.user ?? response.data?.user ?? {};
        setProfile(user);
      } catch (error) {
        console.error('Failed to load profile', error);
        setProfile({});
      } finally {
        setLoading(false);
      }
    };

    void fetchProfile();
  }, []);

  const fullName = [profile.first_name, profile.last_name].filter(Boolean).join(' ') || 'User';

  return (
    <div className="min-h-screen bg-slate-950 px-6 py-10 text-white">
      <div className="mx-auto max-w-4xl">
        <div className="mb-6 flex items-center justify-between">
          <h1 className="text-3xl font-bold">Profile</h1>
          <Link
            href="/dashboard"
            className="rounded-xl border border-white/10 bg-white/5 px-4 py-2 text-sm font-medium text-slate-100 transition hover:border-cyan-400/40 hover:bg-cyan-500/10"
          >
            Back to Dashboard
          </Link>
        </div>

        {loading ? (
          <div className="rounded-2xl border border-white/10 bg-slate-900/80 p-8 text-slate-300">
            Loading profile...
          </div>
        ) : (
          <div className="grid gap-6 md:grid-cols-[260px_1fr]">
            <div className="rounded-2xl border border-white/10 bg-slate-900/80 p-6 text-center shadow-sm">
              <div className="mx-auto mb-4 flex h-20 w-20 items-center justify-center rounded-full bg-gradient-to-br from-cyan-500 to-violet-500 text-2xl font-bold text-white">
                {fullName.charAt(0).toUpperCase()}
              </div>
              <h2 className="text-xl font-semibold">{fullName}</h2>
              <p className="mt-2 text-sm text-slate-400">{profile.status || 'Active'}</p>
            </div>

            <div className="rounded-2xl border border-white/10 bg-slate-900/80 p-6 shadow-sm">
              <div className="grid gap-5 md:grid-cols-2">
                <div>
                  <p className="text-xs uppercase tracking-[0.2em] text-slate-500">First Name</p>
                  <p className="mt-2 text-lg font-medium text-white">{profile.first_name || '—'}</p>
                </div>
                <div>
                  <p className="text-xs uppercase tracking-[0.2em] text-slate-500">Last Name</p>
                  <p className="mt-2 text-lg font-medium text-white">{profile.last_name || '—'}</p>
                </div>
                <div>
                  <p className="text-xs uppercase tracking-[0.2em] text-slate-500">Email</p>
                  <p className="mt-2 text-lg font-medium text-white">{profile.email || '—'}</p>
                </div>
                <div>
                  <p className="text-xs uppercase tracking-[0.2em] text-slate-500">Phone</p>
                  <p className="mt-2 text-lg font-medium text-white">{profile.phone || '—'}</p>
                </div>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
