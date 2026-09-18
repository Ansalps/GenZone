interface Product {
    id: number;
    category_id: number;
    category_name: string;
    product_name: string;
    product_description: string;
    product_image_url: string;
    price: number;
    stock: number;
    popular: boolean;
    size: string;
    discount_percentage?: number;
}

interface ProductCardProps {
    product: Product;
}

export default function ProductCard({ product }: ProductCardProps) {
    const discountPercentage = Number(product.discount_percentage) || 0
    const hasDiscount = discountPercentage > 0
    const discountAmount = hasDiscount ? product.price * (discountPercentage / 100) : 0
    const finalPrice = product.price - discountAmount
    const formattedSizes = product.size
        ? product.size
              .split(',')
              .map((item) => item.trim())
              .filter(Boolean)
              .join(', ')
        : 'N/A'

    return (
        <div className="group overflow-hidden rounded-[1.5rem] border border-white/10 bg-slate-900/80 shadow-lg shadow-slate-950/20 transition duration-300 hover:-translate-y-1 hover:border-cyan-400/40 hover:shadow-cyan-950/20">
            <div className="relative h-64 overflow-hidden">
                <img
                    src={product.product_image_url}
                    alt={product.product_name}
                    className="h-full w-full object-cover transition duration-500 group-hover:scale-105"
                />

                {product.popular && (
                    <span className="absolute left-3 top-3 rounded-full bg-amber-400 px-3 py-1 text-[10px] font-bold uppercase tracking-[0.2em] text-slate-950">
                        Popular
                    </span>
                )}

                {hasDiscount && (
                    <span className="absolute right-3 top-3 rounded-full bg-red-500 px-3 py-1 text-[10px] font-bold uppercase tracking-[0.2em] text-white">
                        {discountPercentage}% off
                    </span>
                )}
            </div>

            <div className="space-y-4 p-5">
                <div>
                    <p className="text-[10px] font-semibold uppercase tracking-[0.25em] text-slate-400">{product.category_name}</p>
                    <h3 className="mt-2 text-xl font-semibold text-white">{product.product_name}</h3>
                </div>

                <p className="line-clamp-2 text-sm text-slate-300">{product.product_description}</p>

                <div className="flex items-center justify-between gap-3">
                    <div>
                        {hasDiscount ? (
                            <div className="flex items-end gap-2">
                                <span className="text-2xl font-bold text-emerald-300">₹{finalPrice.toFixed(2)}</span>
                                <span className="text-sm text-slate-500 line-through">₹{product.price.toFixed(2)}</span>
                            </div>
                        ) : (
                            <span className="text-2xl font-bold text-emerald-300">₹{product.price.toFixed(2)}</span>
                        )}
                    </div>

                    <span className={`rounded-full px-2.5 py-1 text-xs font-medium ${product.stock > 0 ? 'bg-emerald-500/10 text-emerald-200' : 'bg-red-500/10 text-red-200'}`}>
                        {product.stock > 0 ? `${product.stock} in stock` : 'Out of stock'}
                    </span>
                </div>

                <div className="rounded-xl border border-white/10 bg-slate-950/60 px-3 py-2 text-sm text-slate-300">
                    <span className="text-slate-400">Sizes:</span> {formattedSizes}
                </div>
            </div>
        </div>
    )
}