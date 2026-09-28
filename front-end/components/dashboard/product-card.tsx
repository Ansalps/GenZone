"use client"

import { useState } from "react"
import { Product, ProductInventoryItem } from "@/types/product"
import axios from "axios"

interface Props {
    product: Product
    onAdded: (success: boolean) => void
}

export default function DashboardProductCard({ product, onAdded }: Props) {
    const inventory: ProductInventoryItem[] = product.inventory ?? []

    const sizes = inventory.length > 0
        ? inventory.map((i) => i.size)
        : (product.size ? product.size.split(",").map(s => s.trim()).filter(Boolean) : [])

    const defaultSize = sizes[0] ?? ""
    const [selectedSize, setSelectedSize] = useState<string>(defaultSize)
    const [qty, setQty] = useState<number>(1)
    const [loading, setLoading] = useState(false)
    const discountPercentage = Number(product.discount_percentage) || 0
    const hasDiscount = discountPercentage > 0
    const discountAmount = hasDiscount ? product.price * (discountPercentage / 100) : 0
    const finalPrice = product.price - discountAmount
    const showDiscountPlaceholder = false

    const getStockForSize = (size: string) => {
        if (inventory.length > 0) {
            const match = inventory.find((it) => it.size === size)
            return match?.stock ?? 0
        }
        return product.stock ?? 0
    }

    const maxStock = getStockForSize(selectedSize)

    const handleAdd = async () => {
        if (qty <= 0) return onAdded(false)
        if (qty > maxStock) return onAdded(false)

        setLoading(true)
        try {
            await axios.post("http://localhost:8080/cart", {
                product_id: String(product.id),
                quantity: qty,
                size: selectedSize,
            }, { withCredentials: true })

            onAdded(true)
            alert("Added to cart")
        } catch (e) {
            console.error("Add to cart failed", e)
            onAdded(false)
            alert("Failed to add to cart")
        } finally {
            setLoading(false)
        }
    }

    // Offer timing
    let offerText: string | null = null
    let offerStatus: 'upcoming' | 'active' | 'expired' | 'unknown' = 'unknown'
    let daysUntilStart: number | null = null
    let daysLeft: number | null = null

    try {
        const today = new Date()
        if (product.start_date) {
            const s = new Date(product.start_date)
            // use integer day difference (start - today)
            daysUntilStart = Math.floor((Date.UTC(s.getFullYear(), s.getMonth(), s.getDate()) - Date.UTC(today.getFullYear(), today.getMonth(), today.getDate())) / (1000*60*60*24))
        }
        if (product.end_date) {
            const e = new Date(product.end_date)
            daysLeft = Math.floor((Date.UTC(e.getFullYear(), e.getMonth(), e.getDate()) - Date.UTC(today.getFullYear(), today.getMonth(), today.getDate())) / (1000*60*60*24))
        }

        // Determine status:
        // - upcoming: start exists and daysUntilStart > 0
        // - active: start is missing or daysUntilStart <= 0, and end is missing or daysLeft >= 0
        // - expired: end exists and daysLeft < 0
        if (product.start_date && daysUntilStart !== null && daysUntilStart > 0) {
            offerStatus = 'upcoming'
            offerText = `Starts in ${daysUntilStart} day${daysUntilStart===1?'':'s'}`
        } else {
            // either no start_date or start_date is today/past
            if (product.end_date) {
                if (daysLeft !== null && daysLeft >= 0) {
                    offerStatus = 'active'
                    offerText = `Ends in ${daysLeft} day${daysLeft===1?'':'s'}`
                } else if (daysLeft !== null && daysLeft < 0) {
                    offerStatus = 'expired'
                    offerText = `Offer expired`
                }
            } else if (hasDiscount) {
                // no end date, discount present -> active
                offerStatus = 'active'
                offerText = null
            }
        }
    } catch (err) {
        offerStatus = 'unknown'
        offerText = null
        daysLeft = null
        daysUntilStart = null
    }

    const showDiscount = hasDiscount && offerStatus === 'active'

    return (
        <div className="rounded-2xl border border-white/10 bg-slate-900/80 shadow-lg overflow-hidden hover:shadow-2xl transition">

            <div className="relative w-full h-52 bg-gray-800">
                {product.product_image_url ? (
                    <img src={product.product_image_url} alt={product.product_name} className="w-full h-full object-cover" />
                ) : (
                    <div className="w-full h-full flex items-center justify-center text-gray-400">No Image</div>
                )}

                {product.popular && (
                    <span className="absolute top-3 left-3 bg-yellow-400 text-black text-xs font-semibold px-2 py-1 rounded">Popular</span>
                )}

                {showDiscount && (
                    <span className={`absolute top-3 right-3 text-xs font-semibold px-2 py-1 rounded ${offerStatus==='active' ? 'bg-emerald-500 text-white' : offerStatus==='upcoming' ? 'bg-blue-500 text-white' : offerStatus==='expired' ? 'bg-gray-500 text-white' : 'bg-red-600 text-white'}`}>{discountPercentage}% OFF</span>
                )}
            </div>

            <div className="p-4">
                <p className="text-xs text-slate-400 uppercase tracking-wide">{product.category_name}</p>
                <h4 className="mt-2 text-lg font-semibold truncate text-white">{product.product_name}</h4>
                <p className="mt-2 line-clamp-3 text-sm text-slate-300" title={product.product_description}>{product.product_description}</p>

                <div className="mt-4">
                    {showDiscount ? (
                        <div className="flex items-center gap-2">
                            <span className="text-xl font-bold text-emerald-300">₹{finalPrice.toFixed(2)}</span>
                            <span className="text-sm text-slate-400 line-through">₹{product.price.toFixed(2)}</span>
                        </div>
                    ) : (
                        <span className="text-xl font-bold text-white">₹{product.price.toFixed(2)}</span>
                    )}
                </div>

                {showDiscount && (
                    <p className="text-xs text-emerald-300 mt-1">You save ₹{(product.price - finalPrice).toFixed(2)}</p>
                )}

                {showDiscount && offerStatus === 'active' && daysLeft !== null && (
                    <p className="text-xs mt-1"><span className="inline-block rounded-full bg-emerald-600/10 text-emerald-300 px-2 py-1">Offer ends in {daysLeft} day{daysLeft===1?'':'s'}</span></p>
                )}

                {hasDiscount && offerStatus === 'expired' && (
                    <p className="text-xs mt-1"><span className="inline-block rounded-full bg-gray-600/10 text-gray-300 px-2 py-1">Offer expired</span></p>
                )}

                <div className="mt-4 rounded-xl border border-white/10 bg-slate-950/40 p-3 text-sm text-slate-300">
                    {sizes.length > 0 && (
                        <div className="mb-3">
                            <label className="text-[11px] font-semibold uppercase tracking-[0.12em] text-slate-400">Size</label>
                            <div className="mt-2">
                                <select value={selectedSize} onChange={(e) => setSelectedSize(e.target.value)} className="w-full rounded-md border border-white/10 bg-slate-900/70 px-2 py-1 text-white">
                                    {sizes.map((s) => (
                                        <option key={s} value={s}>{s} ({getStockForSize(s)} in stock)</option>
                                    ))}
                                </select>
                            </div>
                        </div>
                    )}

                    <div className="mb-3">
                        <label className="text-[11px] font-semibold uppercase tracking-[0.12em] text-slate-400">Quantity</label>
                        <div className="mt-2">
                            <input type="number" min={1} max={maxStock} value={qty} onChange={(e) => setQty(Math.max(1, Math.min(Number(e.target.value) || 1, maxStock)))} className="w-full rounded-md border border-white/10 bg-slate-900/70 px-2 py-1 text-white" />
                            <p className="text-xs text-slate-400 mt-1">Max: {maxStock}</p>
                        </div>
                    </div>

                    <div>
                        <button onClick={handleAdd} disabled={loading || qty < 1 || qty > maxStock} className="w-full rounded-xl bg-gradient-to-r from-emerald-500 to-teal-500 text-white text-sm px-3 py-2 font-medium transition hover:opacity-95 disabled:opacity-60">
                            {loading ? 'Adding...' : 'Add to cart'}
                        </button>
                    </div>
                </div>
            </div>
        </div>
    )
}
