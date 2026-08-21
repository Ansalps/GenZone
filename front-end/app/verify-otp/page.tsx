"use client"

import { useState, Suspense } from "react";
import { useSearchParams, useRouter } from "next/navigation";
import axios from "axios";
import { z } from "zod";
import { toast } from "sonner";

const otpSchema = z.object({
  otp: z
    .string()
    .min(1, "OTP is required")
    .max(20, "OTP is too long"),
});

function VerifyOtpForm() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const email = searchParams.get("email") || "";

  const [otp, setOtp] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [resending, setResending] = useState(false);

  async function handleSubmit(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault();
    setError(null);

    const validationResult = otpSchema.safeParse({ otp });

    if (!validationResult.success) {
      const errorMessage = validationResult.error.issues[0].message;
      setError(errorMessage);
      toast.error(errorMessage);
      return;
    }

    if (!email) {
      toast.error("Email parameter missing from URL");
      return;
    }

    try {
      setLoading(true);
      const response = await axios.post(
        `http://localhost:8080/public/verify-otp?email=${encodeURIComponent(email)}`,
        { otp: validationResult.data.otp }, // Fixed: Using validationResult instead of formData
        { withCredentials: true }
      );

      toast.success(response.data?.message || "OTP verified successfully");
      router.push("/login");
    } catch (err: any) {
      toast.error(err?.response?.data?.message || "Failed to verify OTP");
    } finally {
      setLoading(false);
    }
  }

  async function handleResendOtp() {
    if (!email) {
      toast.error("Email parameter missing from URL");
      return;
    }

    try {
      setResending(true);
      const response = await axios.post(
        `http://localhost:8080/public/resendotp/${encodeURIComponent(email)}`,
        {},
        { withCredentials: true }
      );

      toast.success(response.data?.message || "OTP resent successfully");
    } catch (err: any) {
      toast.error(err?.response?.data?.message || "Failed to resend OTP");
    } finally {
      setResending(false);
    }
  }

  return (
    <div className="flex flex-col justify-center items-center min-h-screen">
      <form onSubmit={handleSubmit} className="flex flex-col gap-4 w-full max-w-xs">
        <div className="flex flex-col">
          <label htmlFor="otp" className="mb-1 font-medium">OTP Code</label>
          <input
            id="otp"
            type="text"
            maxLength={20}
            value={otp}
            onChange={(e) => setOtp(e.target.value)}
            className="border-black border p-2 text-center tracking-widest text-lg uppercase"
            placeholder="HRDO4SO2"
          />
          {error && <span className="text-red-500 text-sm mt-1">{error}</span>}
        </div>

        <button
          type="submit"
          disabled={loading}
          className="bg-black text-white p-2 rounded disabled:opacity-50"
        >
          {loading ? "Verifying..." : "Submit OTP"}
        </button>

        <button
          type="button"
          disabled={resending}
          onClick={handleResendOtp}
          className="border border-black p-2 rounded disabled:opacity-50"
        >
          {resending ? "Resending..." : "Resend OTP"}
        </button>
      </form>
    </div>
  );
}

// Wrapped with Suspense for App Router searchParams support
export default function VerifyOtp() {
  return (
    <Suspense fallback={<div>Loading...</div>}>
      <VerifyOtpForm />
    </Suspense>
  );
}