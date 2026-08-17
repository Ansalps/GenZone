import Link from "next/link";

import { Product } from "@/types/product";

interface ProductCardProps {
    product: Product;
    onDelete: (id: number) => void;
}

export default function ProductCard({
    product,
    onDelete,
}: ProductCardProps) {
    const discountPercent = product.discount_percentage || 0;
    const hasOffer = discountPercent > 0;

    const discountAmount = hasOffer
        ? product.price * (discountPercent / 100)
        : 0;

    const finalPrice = product.price - discountAmount;

    return (
        <div className="bg-white rounded-xl shadow-sm border overflow-hidden hover:shadow-md transition">

            {/* Product Image */}
            <div className="relative w-full h-56 bg-gray-100">
                {product.product_image_url ? (
                    <img
                        src={product.product_image_url}
                        alt={product.product_name}
                        className="w-full h-full object-cover"
                    />
                ) : (
                    <div className="w-full h-full flex items-center justify-center text-gray-400">
                        No Image
                    </div>
                )}

                {/* Popular badge */}
                {product.popular && (
                    <span className="absolute top-3 left-3 bg-yellow-500 text-white text-xs font-semibold px-2 py-1 rounded">
                        Popular
                    </span>
                )}

                {/* Offer badge */}
                {hasOffer && (
                    <span className="absolute top-3 right-3 bg-red-500 text-white text-xs font-semibold px-2 py-1 rounded">
                        {discountPercent}% OFF
                    </span>
                )}
            </div>

            {/* Product information */}
            <div className="p-5">

                {/* Category */}
                <p className="text-xs text-gray-500 uppercase tracking-wide">
                    {product.category_name}
                </p>

                {/* Name */}
                <h2 className="text-lg font-semibold mt-1 truncate">
                    {product.product_name}
                </h2>

                {/* Description */}
                <p
                    className="text-sm text-gray-500 mt-2 line-clamp-2"
                    title={product.product_description}
                >
                    {product.product_description}
                </p>

                {/* Price */}
                <div className="mt-4">
                    {hasOffer ? (
                        <div className="flex items-center gap-2">
                            <span className="text-xl font-bold text-green-700">
                                ₹{finalPrice.toFixed(2)}
                            </span>

                            <span className="text-sm text-gray-400 line-through">
                                ₹{product.price.toFixed(2)}
                            </span>
                        </div>
                    ) : (
                        <span className="text-xl font-bold">
                            ₹{product.price.toFixed(2)}
                        </span>
                    )}
                </div>

                {/* Discount */}
                {hasOffer && (
                    <p className="text-xs text-green-600 mt-1">
                        You save ₹{discountAmount.toFixed(2)}
                    </p>
                )}

                {/* Product details */}
                <div className="grid grid-cols-2 gap-3 mt-4 text-sm">

                    <div>
                        <span className="text-gray-500">
                            Stock
                        </span>

                        <p
                            className={
                                product.stock === 0
                                    ? "font-semibold text-red-600"
                                    : "font-semibold"
                            }
                        >
                            {product.stock}
                        </p>
                    </div>

                    <div>
                        <span className="text-gray-500">
                            Size
                        </span>

                        <p className="font-semibold">
                            {product.size}
                        </p>
                    </div>

                </div>

                {/* Actions */}
                <div className="flex gap-2 mt-5">

                    <Link
                        href={`/admin/products/edit/${product.id}`}
                        className="flex-1 text-center bg-blue-500 hover:bg-blue-600 text-white px-3 py-2 rounded text-sm font-medium transition"
                    >
                        Edit
                    </Link>

                    <button
                        onClick={() => onDelete(product.id)}
                        className="flex-1 bg-red-500 hover:bg-red-600 text-white px-3 py-2 rounded text-sm font-medium transition"
                    >
                        Delete
                    </button>

                </div>

            </div>
        </div>
    );
}