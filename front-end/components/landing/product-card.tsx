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

export default function ProductCard({
    product,
}: ProductCardProps) {

    const discountPercentage =
        Number(product.discount_percentage) || 0;

    const hasDiscount = discountPercentage > 0;

    const discountAmount = hasDiscount
        ? product.price * (discountPercentage / 100)
        : 0;

    const finalPrice =
        product.price - discountAmount;

    return (
        <div className="bg-white rounded-xl shadow-md overflow-hidden hover:shadow-xl transition duration-300">

            {/* Image */}
            <div className="relative h-64 overflow-hidden">

                <img
                    src={product.product_image_url}
                    alt={product.product_name}
                    className="w-full h-full object-cover hover:scale-105 transition duration-300"
                />

                {/* Popular */}
                {product.popular && (
                    <span className="absolute top-3 left-3 bg-yellow-400 text-yellow-900 text-xs font-bold px-3 py-1 rounded-full">
                        Popular
                    </span>
                )}

                {/* Discount */}
                {hasDiscount && (
                    <span className="absolute top-3 right-3 bg-red-500 text-white text-xs font-bold px-3 py-1 rounded-full">
                        {discountPercentage}% OFF
                    </span>
                )}
            </div>

            {/* Content */}
            <div className="p-4">

                <p className="text-xs text-gray-500 mb-1">
                    {product.category_name}
                </p>

                <h3 className="text-lg font-semibold text-gray-800">
                    {product.product_name}
                </h3>

                <p className="text-sm text-gray-500 mt-1 line-clamp-2">
                    {product.product_description}
                </p>

                {/* Price */}
                <div className="mt-4">

                    {hasDiscount ? (
                        <div className="flex items-center gap-2">

                            <span className="text-xl font-bold text-green-700">
                                ₹{finalPrice.toFixed(2)}
                            </span>

                            <span className="text-sm text-gray-400 line-through">
                                ₹{product.price.toFixed(2)}
                            </span>

                        </div>
                    ) : (
                        <span className="text-xl font-bold text-green-700">
                            ₹{product.price.toFixed(2)}
                        </span>
                    )}

                </div>

                {/* Size */}
                <div className="mt-3 text-sm text-gray-500">
                    Size:{" "}
                    <span className="font-medium text-gray-700">
                        {product.size}
                    </span>
                </div>

            </div>
        </div>
    );
}