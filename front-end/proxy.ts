import { NextResponse } from 'next/server';
import type { NextRequest } from 'next/server';

export function proxy(request: NextRequest) {
    // 1. Extract the HttpOnly cookie set by your Go backend
    const token = request.cookies.get('jwt_token')?.value;
    const { pathname } = request.nextUrl;

    const isLoginPage = pathname === '/admin/login';
    const isAdminRoute = pathname.startsWith('/admin');

    // Case 1: Unauthenticated user trying to access ANY admin route (except /admin/login)
    if (isAdminRoute && !isLoginPage && !token) {
        const loginUrl = new URL('/admin/login', request.url);
        return NextResponse.redirect(loginUrl);
    }

    // Case 2: Already authenticated user trying to access /admin/login
    if (isLoginPage && token) {
        const dashboardUrl = new URL('/admin', request.url);
        return NextResponse.redirect(dashboardUrl);
    }

    return NextResponse.next();
}

// Ensure matcher catches /admin, /admin/login, and all sub-routes
export const config = {
    matcher: ['/admin', '/admin/:path*'],
};