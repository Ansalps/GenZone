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

interface AddressItem {
  id?: number;
  country?: string;
  state?: string;
  city?: string;
  district?: string;
  street_name?: string;
  pin_code?: string;
  phone?: string;
  default?: boolean;
  Default?: boolean;
}

interface AddressFormState {
  country: string;
  state: string;
  city: string;
  street_name: string;
  pin_code: string;
  phone: string;
}

const API = 'http://localhost:8080';

const createEmptyAddressForm = (): AddressFormState => ({
  country: '',
  state: '',
  city: '',
  street_name: '',
  pin_code: '',
  phone: '',
});

export default function ProfilePage() {
  const [profile, setProfile] = useState<UserProfile>({});
  const [addresses, setAddresses] = useState<AddressItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [addressesLoading, setAddressesLoading] = useState(true);
  const [form, setForm] = useState<AddressFormState>(createEmptyAddressForm());
  const [editingAddressId, setEditingAddressId] = useState<number | null>(null);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState('');

  const fetchProfile = async () => {
    try {
      const response = await axios.get(`${API}/profile`, { withCredentials: true });
      const user = response.data?.data?.User ?? response.data?.user ?? {};
      setProfile(user);
    } catch (err) {
      console.error('Failed to load profile', err);
      setProfile({});
    }
  };

  const fetchAddresses = async () => {
    try {
      const response = await axios.get(`${API}/profile/useraddress`, { withCredentials: true });
      const list = response.data?.data?.address ?? response.data?.address ?? [];
      setAddresses(Array.isArray(list) ? list : []);
    } catch (err) {
      console.error('Failed to load addresses', err);
      setAddresses([]);
    } finally {
      setAddressesLoading(false);
    }
  };

  useEffect(() => {
    const loadData = async () => {
      setLoading(true);
      try {
        await Promise.all([fetchProfile(), fetchAddresses()]);
      } finally {
        setLoading(false);
      }
    };

    void loadData();
  }, []);

  const fullName = [profile.first_name, profile.last_name].filter(Boolean).join(' ') || 'User';

  const handleFormChange = (
    event: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>
  ) => {
    const { name, value, type } = event.target;

    setForm((previous) => ({
      ...previous,
      [name]: type === 'checkbox' ? (event.target as HTMLInputElement).checked : value,
    }));
  };

  const resetForm = () => {
    setForm(createEmptyAddressForm());
    setEditingAddressId(null);
    setError('');
  };

  const submitAddress = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setError('');

    const payload = {
      ...form,
      default: false,
    };

    try {
      setSubmitting(true);

      if (editingAddressId) {
        await axios.put(`${API}/profile/useraddress/${editingAddressId}`, payload, {
          withCredentials: true,
        });
      } else {
        await axios.post(`${API}/profile/useraddress`, payload, {
          withCredentials: true,
        });
      }

      resetForm();
      await fetchAddresses();
    } catch (err) {
      console.error('Failed to save address', err);
      setError(
        axios.isAxiosError(err)
          ? err.response?.data?.message || err.response?.data?.error || 'Unable to save address.'
          : 'Unable to save address.'
      );
    } finally {
      setSubmitting(false);
    }
  };

  const handleEdit = (address: AddressItem) => {
    setEditingAddressId(Number(address.id));
    setForm({
      country: address.country ?? '',
      state: address.state ?? '',
      city: address.city ?? address.district ?? '',
      street_name: address.street_name ?? '',
      pin_code: address.pin_code ?? '',
      phone: address.phone ?? '',
    });
    setError('');
  };

  const handleDelete = async (addressId: number) => {
    if (!window.confirm('Delete this address?')) {
      return;
    }

    try {
      await axios.delete(`${API}/profile/useraddress/${addressId}`, { withCredentials: true });
      if (editingAddressId === addressId) {
        resetForm();
      }
      await fetchAddresses();
    } catch (err) {
      console.error('Failed to delete address', err);
      setError(
        axios.isAxiosError(err)
          ? err.response?.data?.message || err.response?.data?.error || 'Unable to delete address.'
          : 'Unable to delete address.'
      );
    }
  };

  const setAddressAsDefault = async (addressId: number) => {
    try {
      await axios.put(`${API}/profile/useraddress/default/${addressId}`, {}, { withCredentials: true });
      await fetchAddresses();
    } catch (err) {
      console.error('Failed to set default address', err);
      setError(
        axios.isAxiosError(err)
          ? err.response?.data?.message || err.response?.data?.error || 'Unable to set default address.'
          : 'Unable to set default address.'
      );
    }
  };

  return (
    <div className="min-h-screen bg-slate-950 px-6 py-10 text-white">
      <div className="mx-auto max-w-5xl">
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
          <>
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

            <div className="mt-8 rounded-2xl border border-white/10 bg-slate-900/80 p-6">
              <div className="mb-6 flex items-center justify-between gap-4">
                <div>
                  <h2 className="text-2xl font-bold text-white">Addresses</h2>
                  <p className="text-sm text-slate-400">Manage all your saved shipping addresses.</p>
                </div>
                {!editingAddressId && (
                  <button
                    type="button"
                    onClick={() => setForm(createEmptyAddressForm())}
                    className="rounded-xl bg-cyan-500 px-4 py-2 text-sm font-semibold text-slate-950 transition hover:bg-cyan-400"
                  >
                    Add Address
                  </button>
                )}
              </div>

              <div className="grid gap-6 lg:grid-cols-[1.3fr_0.7fr]">
                <div className="space-y-4">
                  {addressesLoading ? (
                    <div className="rounded-xl border border-white/10 bg-slate-950/60 p-5 text-slate-300">
                      Loading addresses...
                    </div>
                  ) : addresses.length === 0 ? (
                    <div className="rounded-xl border border-dashed border-slate-700 bg-slate-950/60 p-5 text-slate-300">
                      No addresses saved yet.
                    </div>
                  ) : (
                    addresses.map((address) => (
                      <div
                        key={address.id}
                        className="rounded-xl border border-white/10 bg-slate-950/60 p-4"
                      >
                        <div className="flex items-start justify-between gap-4">
                          <div>
                            <div className="mb-2 flex items-center gap-2">
                              <p className="text-base font-semibold text-white">
                                {address.street_name || 'Address'}
                              </p>
                              {Boolean(address.default ?? address.Default) && (
                                <span className="rounded-full border border-emerald-400/30 bg-emerald-500/10 px-2 py-0.5 text-[10px] font-semibold uppercase tracking-[0.15em] text-emerald-300">
                                  Default
                                </span>
                              )}
                            </div>
                            <p className="text-sm text-slate-300">
                              {address.street_name}, {address.city ?? address.district}, {address.state}, {address.country}
                            </p>
                            <p className="mt-1 text-sm text-slate-300">Pin: {address.pin_code}</p>
                            <p className="mt-1 text-sm text-slate-300">Phone: {address.phone}</p>
                          </div>

                          <div className="flex flex-col gap-2">
                            {!Boolean(address.default ?? address.Default) && (
                              <button
                                type="button"
                                onClick={() => setAddressAsDefault(Number(address.id))}
                                className="rounded-lg border border-emerald-400/30 bg-emerald-500/10 px-3 py-1.5 text-xs font-medium text-emerald-200 transition hover:bg-emerald-500/20"
                              >
                                Set as default
                              </button>
                            )}
                            <button
                              type="button"
                              onClick={() => handleEdit(address)}
                              className="rounded-lg border border-cyan-400/30 bg-cyan-500/10 px-3 py-1.5 text-xs font-medium text-cyan-200 transition hover:bg-cyan-500/20"
                            >
                              Edit
                            </button>
                            <button
                              type="button"
                              onClick={() => handleDelete(Number(address.id))}
                              className="rounded-lg border border-rose-400/30 bg-rose-500/10 px-3 py-1.5 text-xs font-medium text-rose-200 transition hover:bg-rose-500/20"
                            >
                              Delete
                            </button>
                          </div>
                        </div>
                      </div>
                    ))
                  )}
                </div>

                <form onSubmit={submitAddress} className="rounded-xl border border-white/10 bg-slate-950/60 p-4">
                  <h3 className="mb-4 text-xl font-semibold text-white">
                    {editingAddressId ? 'Edit Address' : 'Add New Address'}
                  </h3>

                  {error && (
                    <div className="mb-4 rounded-lg border border-rose-500/40 bg-rose-500/10 px-3 py-2 text-sm text-rose-200">
                      {error}
                    </div>
                  )}

                  <div className="space-y-4">
                    <div>
                      <label className="mb-1 block text-sm text-slate-300">Country</label>
                      <input
                        name="country"
                        value={form.country}
                        onChange={handleFormChange}
                        className="w-full rounded-lg border border-slate-700 bg-slate-900 px-3 py-2 text-white outline-none transition focus:border-cyan-400"
                        placeholder="India"
                        required
                      />
                    </div>

                    <div>
                      <label className="mb-1 block text-sm text-slate-300">State</label>
                      <input
                        name="state"
                        value={form.state}
                        onChange={handleFormChange}
                        className="w-full rounded-lg border border-slate-700 bg-slate-900 px-3 py-2 text-white outline-none transition focus:border-cyan-400"
                        placeholder="Kerala"
                        required
                      />
                    </div>

                    <div>
                      <label className="mb-1 block text-sm text-slate-300">Town/City</label>
                      <input
                        name="city"
                        value={form.city}
                        onChange={handleFormChange}
                        className="w-full rounded-lg border border-slate-700 bg-slate-900 px-3 py-2 text-white outline-none transition focus:border-cyan-400"
                        placeholder="Kochi"
                        required
                      />
                    </div>

                    <div>
                      <label className="mb-1 block text-sm text-slate-300">Street Name</label>
                      <input
                        name="street_name"
                        value={form.street_name}
                        onChange={handleFormChange}
                        className="w-full rounded-lg border border-slate-700 bg-slate-900 px-3 py-2 text-white outline-none transition focus:border-cyan-400"
                        placeholder="MG Road"
                        required
                      />
                    </div>

                    <div>
                      <label className="mb-1 block text-sm text-slate-300">Pin Code</label>
                      <input
                        name="pin_code"
                        value={form.pin_code}
                        onChange={handleFormChange}
                        className="w-full rounded-lg border border-slate-700 bg-slate-900 px-3 py-2 text-white outline-none transition focus:border-cyan-400"
                        placeholder="682001"
                        required
                      />
                    </div>

                    <div>
                      <label className="mb-1 block text-sm text-slate-300">Phone</label>
                      <input
                        name="phone"
                        value={form.phone}
                        onChange={handleFormChange}
                        className="w-full rounded-lg border border-slate-700 bg-slate-900 px-3 py-2 text-white outline-none transition focus:border-cyan-400"
                        placeholder="9876543210"
                        required
                      />
                    </div>

                  </div>

                  <div className="mt-5 flex gap-3">
                    <button
                      type="submit"
                      disabled={submitting}
                      className="flex-1 rounded-lg bg-cyan-500 px-4 py-2.5 font-semibold text-slate-950 transition hover:bg-cyan-400 disabled:cursor-not-allowed disabled:opacity-60"
                    >
                      {submitting ? 'Saving...' : editingAddressId ? 'Update Address' : 'Save Address'}
                    </button>
                    <button
                      type="button"
                      onClick={resetForm}
                      className="rounded-lg border border-white/10 bg-white/5 px-4 py-2.5 font-medium text-slate-200 transition hover:bg-white/10"
                    >
                      Cancel
                    </button>
                  </div>
                </form>
              </div>
            </div>
          </>
        )}
      </div>
    </div>
  );
}
