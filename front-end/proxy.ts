import { NextResponse } from 'next/server';
import type { NextRequest } from 'next/server';

const protectedUserRoutes = [
  '/dashboard',
  '/cart',
  '/checkout',
  '/orders',
  '/profile',
  '/wallet',
  '/wishlist',
  '/payment-option',
];

export function proxy(request: NextRequest) {
    const { pathname } = request.nextUrl;
    const isAdminRoute = pathname.startsWith('/admin') && pathname !== '/admin/login';
    const isUserProtectedRoute = protectedUserRoutes.some((route) => pathname === route || pathname.startsWith(`${route}/`));

    const adminToken = request.cookies.get('jwt_admin_token')?.value || request.cookies.get('jwt_token')?.value;
    const userToken = request.cookies.get('jwt_user_token')?.value || request.cookies.get('jwt_token')?.value;

    if (isAdminRoute) {
        if (!adminToken) {
            return NextResponse.redirect(new URL('/admin/login', request.url));
        }
        return NextResponse.next();
    }

    if (pathname === '/admin/login' && adminToken) {
        return NextResponse.redirect(new URL('/admin', request.url));
    }

    if (isUserProtectedRoute && !userToken) {
        return NextResponse.redirect(new URL('/login', request.url));
    }

    if ((pathname === '/login' || pathname === '/signup' || pathname === '/verify-otp') && userToken) {
        return NextResponse.redirect(new URL('/dashboard', request.url));
    }

    return NextResponse.next();
}

export const config = {
    matcher: ['/admin/:path*', '/dashboard/:path*', '/cart/:path*', '/checkout/:path*', '/orders/:path*', '/profile/:path*', '/wallet/:path*', '/wishlist/:path*', '/payment-option/:path*', '/login', '/signup', '/verify-otp'],
};