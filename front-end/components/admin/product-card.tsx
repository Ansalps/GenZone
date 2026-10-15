import Link from "next/link";

import { Product } from "@/types/product";

interface ProductCardProps {
    product: Product;
    onRequestDelete: (id: number, name?: string) => void;
}

export default function ProductCard({
    product,
    onRequestDelete,
}: ProductCardProps) {
    const hasDeletedTimestamp = Boolean(product.deleted_at && product.deleted_at !== '0001-01-01T00:00:00Z' && product.deleted_at !== '0001-01-01T00:00:00.000Z');
    const isDiscontinued = ((product.product_status || '').toLowerCase() === 'discontinued') || hasDeletedTimestamp;
    const discountPercent = product.discount_percentage || 0;
    const hasOffer = discountPercent > 0;

    const discountAmount = hasOffer
        ? product.price * (discountPercent / 100)
        : 0;

    const finalPrice = product.price - discountAmount;
    const inventory = product.inventory && product.inventory.length > 0 ? product.inventory : [];
    const sizeOrder = ["Small", "Medium", "Large"];
    const shortSizeMap: Record<string, string> = {
        Small: 'S',
        Medium: 'M',
        Large: 'L',
    };
    const inventoryRows = sizeOrder.map((size) => {
        const match = inventory.find((item) => item.size === size);
        return {
            size,
            label: shortSizeMap[size] || size.slice(0, 1).toUpperCase(),
            stock: match?.stock ?? 0,
        };
    });

    // Offer timing calculation
    let daysLeft: number | null = null
    let daysUntilStart: number | null = null
    let offerStatus: 'upcoming' | 'active' | 'expired' | 'unknown' = 'unknown'

    try {
        const today = new Date()
        if (product.start_date) {
            const s = new Date(product.start_date)
            daysUntilStart = Math.floor((Date.UTC(s.getFullYear(), s.getMonth(), s.getDate()) - Date.UTC(today.getFullYear(), today.getMonth(), today.getDate())) / (1000 * 60 * 60 * 24))
        }
        if (product.end_date) {
            const e = new Date(product.end_date)
            daysLeft = Math.floor((Date.UTC(e.getFullYear(), e.getMonth(), e.getDate()) - Date.UTC(today.getFullYear(), today.getMonth(), today.getDate())) / (1000 * 60 * 60 * 24))
        }

        if (product.start_date && daysUntilStart !== null && daysUntilStart > 0) {
            offerStatus = 'upcoming'
        } else {
            // start is today or past, or not provided
            if (product.end_date) {
                if (daysLeft !== null && daysLeft >= 0) {
                    offerStatus = 'active'
                } else if (daysLeft !== null && daysLeft < 0) {
                    offerStatus = 'expired'
                }
            } else if (hasOffer) {
                offerStatus = 'active'
            }
        }
    } catch (e) {
        offerStatus = 'unknown'
        daysLeft = null
        daysUntilStart = null
    }

    const showDiscount = hasOffer && offerStatus === 'active' && !isDiscontinued

    return (
        <div className={`rounded-2xl border overflow-hidden shadow-lg transition ${isDiscontinued ? 'border-slate-700 bg-slate-800/90 opacity-80 grayscale-[0.2]' : 'border-white/10 bg-slate-900/80 hover:shadow-2xl'}`}>

            {/* Product Image */}
            <div className="relative w-full h-52 bg-gray-800">
                {product.product_image_url ? (
                    <img
                        src={product.product_image_url}
                        alt={product.product_name}
                        className={`w-full h-full object-cover ${isDiscontinued ? 'grayscale' : ''}`}
                    />
                ) : (
                    <div className="w-full h-full flex items-center justify-center text-gray-400">
                        No Image
                    </div>
                )}

                {/* Popular badge */}
                {product.popular && (
                    <span className="absolute top-3 left-3 bg-yellow-400 text-black text-xs font-semibold px-2 py-1 rounded">
                        Popular
                    </span>
                )}

                {isDiscontinued && (
                    <span className="absolute top-3 right-3 bg-red-500/90 text-white text-xs font-semibold px-2 py-1 rounded">
                        Discontinued
                    </span>
                )}

                {/* Offer badge */}
                {showDiscount && (
                    <span className={`absolute top-3 right-3 text-xs font-semibold px-2 py-1 rounded ${offerStatus==='active' ? 'bg-emerald-500 text-white' : offerStatus==='upcoming' ? 'bg-blue-500 text-white' : offerStatus==='expired' ? 'bg-gray-500 text-white' : 'bg-red-600 text-white'}`}>
                        {discountPercent}% OFF
                    </span>
                )}
            </div>

            {/* Product information */}
            <div className="p-4">

                {/* Category */}
                <p className="text-xs text-slate-400 uppercase tracking-wide">
                    {product.category_name}
                </p>

                {/* Name */}
                <h2 className={`text-lg font-semibold mt-1 truncate ${isDiscontinued ? 'text-slate-300' : 'text-white'}`}>
                    {product.product_name}
                </h2>

                {/* Description */}
                <p className={`text-sm mt-2 line-clamp-2 ${isDiscontinued ? 'text-slate-400' : 'text-slate-300'}`} title={product.product_description}>
                    {product.product_description}
                </p>

                {/* Price */}
                <div className="mt-4">
                    {hasOffer && !isDiscontinued ? (
                        <div className="flex items-center gap-2">
                            <span className="text-xl font-bold text-emerald-300">
                                ₹{finalPrice.toFixed(2)}
                            </span>

                            <span className="text-sm text-slate-400 line-through">
                                ₹{product.price.toFixed(2)}
                            </span>
                        </div>
                    ) : (
                        <span className={`text-xl font-bold ${isDiscontinued ? 'text-slate-400' : 'text-white'}`}>₹{product.price.toFixed(2)}</span>
                    )}
                </div>

                {/* Discount */}
                {showDiscount && (
                    <p className="text-xs text-emerald-300 mt-1">You save ₹{discountAmount.toFixed(2)}</p>
                )}

                {showDiscount && offerStatus === 'active' && daysLeft !== null && (
                    <p className="text-xs mt-1">
                        <span className="inline-block rounded-full bg-emerald-600/10 text-emerald-300 px-2 py-1">Offer ends in {daysLeft} day{daysLeft===1?'':'s'}</span>
                    </p>
                )}

                {hasOffer && offerStatus === 'expired' && (
                    <p className="text-xs mt-1">
                        <span className="inline-block rounded-full bg-gray-600/10 text-gray-300 px-2 py-1">Offer expired</span>
                    </p>
                )}

                {/* Product details */}
                <div className="mt-4 rounded-xl border border-white/10 bg-slate-950/40 p-3 text-sm text-slate-300">
                    <div className="mb-2 flex items-center justify-between text-[11px] font-semibold uppercase tracking-[0.2em] text-slate-400">
                        <span>Sizes</span>
                        <span>Stock</span>
                    </div>

                    <div className="space-y-2">
                        <div className="flex items-center justify-between gap-2">
                            <div className="flex gap-2">
                                {inventoryRows.map((row) => (
                                    <span key={row.size} className="inline-flex min-w-7 justify-center rounded-lg border border-white/10 bg-slate-900 px-2 py-1 text-xs font-semibold text-slate-100">
                                        {row.label}
                                    </span>
                                ))}
                            </div>
                            <div className="flex gap-2">
                                {inventoryRows.map((row) => (
                                    <span key={`${row.size}-stock`} className={`inline-flex min-w-10 justify-center rounded-lg px-2 py-1 text-xs font-semibold ${row.stock === 0 ? 'bg-red-500/10 text-red-300' : 'bg-emerald-500/10 text-emerald-300'}`}>
                                        {row.stock}
                                    </span>
                                ))}
                            </div>
                        </div>
                    </div>
                </div>

                {/* Actions */}
                <div className="flex gap-2 mt-5">
                    <Link
                        href={`/admin/products/edit/${product.id}`}
                        className={`flex-1 text-center rounded-xl ${isDiscontinued ? 'bg-slate-700 text-slate-300 cursor-not-allowed opacity-70' : 'bg-gradient-to-r from-emerald-500 to-teal-500 text-white hover:opacity-95'} text-sm px-3 py-2 font-medium transition`}
                        onClick={(event) => {
                            if (isDiscontinued) {
                                event.preventDefault();
                            }
                        }}
                    >
                        Edit
                    </Link>

                    <Link
                        href={`/admin/products/update-stock/${product.id}`}
                        className={`flex-1 text-center rounded-xl ${isDiscontinued ? 'bg-slate-700 text-slate-300 cursor-not-allowed opacity-70' : 'bg-gradient-to-r from-violet-500 to-purple-500 text-white hover:opacity-95'} text-sm px-3 py-2 font-medium transition`}
                        onClick={(event) => {
                            if (isDiscontinued) {
                                event.preventDefault();
                            }
                        }}
                    >
                        Update Stock
                    </Link>

                    <button
                        onClick={() => onRequestDelete(product.id, product.product_name)}
                        className={`flex-1 rounded-xl ${isDiscontinued ? 'bg-slate-700 text-slate-300 cursor-not-allowed opacity-70' : 'bg-red-600 text-white hover:opacity-95'} text-sm px-3 py-2 font-medium transition`}
                        disabled={isDiscontinued}
                    >
                        {isDiscontinued ? 'Discontinued' : 'Delete'}
                    </button>
                </div>

            </div>
        </div>
    );
}