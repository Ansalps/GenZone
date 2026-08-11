interface Category {
    id: number;
    category_name: string;
}

// Add Validation Errors interface
export interface FormErrors {
    categoryName?: string;
    productName?: string;
    productDescription?: string;
    productImageUrl?: string;
    price?: string;
    stock?: string;
    size?: string;
    discountPercentage?: string;
}

interface ProductFormProps {
    formData: {
        categoryName: string;
        productName: string;
        productDescription: string;
        productImageUrl: string;
        price: number;
        stock: number;
        size: string;
        popular: boolean;
        discountPercentage: number;
        // discountAmount: number;
        // totalDiscountedAmount: number;
    };
    onChange: (
        e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>
    ) => void;
    isLoading: boolean;
    categories: Category[];
    errors?: FormErrors; // <-- Passed from parent component
}

export default function ProductForm({
    formData,
    onChange,
    isLoading,
    categories,
    errors = {}, // Default to empty object
}: ProductFormProps) {
    // Dynamically calculate live preview amounts on the frontend
    const price = Number(formData.price) || 0;
    const discountPercent = Number(formData.discountPercentage) || 0;

    const discountAmount = price > 0 && discountPercent > 0 
        ? ((price * discountPercent) / 100).toFixed(2) 
        : "0.00";

    const finalPrice = price > 0 && discountPercent > 0 
        ? (price - Number(discountAmount)).toFixed(2) 
        : price.toFixed(2);

    return (
        <>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-5">

                {/* Category */}
                <div className="flex flex-col">
                    <label htmlFor="categoryName" className="mb-1">
                        Category *
                    </label>

                    <select
                        id="categoryName"
                        name="categoryName"
                        value={formData.categoryName}
                        onChange={onChange}
                        className="border rounded p-2 bg-white"
                    >
                        <option value="">Select Category</option>

                        {categories.map((category) => (
                            <option
                                key={category.id}
                                value={category.category_name}
                            >
                                {category.category_name}
                            </option>
                        ))}
                    </select>
                    {errors.categoryName && (
                        <span className="text-red-500 text-xs mt-1">{errors.categoryName}</span>
                    )}
                </div>

                {/* Product Name */}
                <div className="flex flex-col">
                    <label htmlFor="productName" className="mb-1">
                        Product Name *
                    </label>

                    <input
                        type="text"
                        id="productName"
                        name="productName"
                        value={formData.productName}
                        onChange={onChange}
                        className="border rounded p-2"
                        required
                    />
                    {errors.productName && (
                        <span className="text-red-500 text-xs mt-1">{errors.productName}</span>
                    )}
                </div>

                {/* Description */}
                <div className="flex flex-col">
                    <label htmlFor="productDescription" className="mb-1">
                        Description *
                    </label>

                    <input
                        type="text"
                        id="productDescription"
                        name="productDescription"
                        value={formData.productDescription}
                        onChange={onChange}
                        className="border rounded p-2"
                        required
                    />
                    {errors.productDescription && (
                        <span className="text-red-500 text-xs mt-1">{errors.productDescription}</span>
                    )}
                </div>

                {/* Image URL */}
                <div className="flex flex-col">
                    <label htmlFor="productImageUrl" className="mb-1">
                        Image URL *
                    </label>

                    <input
                        type="text"
                        id="productImageUrl"
                        name="productImageUrl"
                        value={formData.productImageUrl}
                        onChange={onChange}
                        className="border rounded p-2"
                        required
                    />
                    {errors.productImageUrl && (
                        <span className="text-red-500 text-xs mt-1">{errors.productImageUrl}</span>
                    )}
                </div>

                {/* Price */}
                <div className="flex flex-col">
                    <label htmlFor="price" className="mb-1">
                        Original Price ($) *
                    </label>

                    <input
                        type="number"
                        id="price"
                        name="price"
                        step="0.01"
                        min="0"
                        value={formData.price || ''}
                        onChange={onChange}
                        className="border rounded p-2"
                        required
                    />
                    {errors.price && (
                        <span className="text-red-500 text-xs mt-1">{errors.price}</span>
                    )}
                </div>

                {/* Stock */}
                <div className="flex flex-col">
                    <label htmlFor="stock" className="mb-1">
                        Stock *
                    </label>

                    <input
                        type="number"
                        id="stock"
                        name="stock"
                        min="0"
                        value={formData.stock || ''}
                        onChange={onChange}
                        className="border rounded p-2"
                        required
                    />
                    {errors.stock && (
                        <span className="text-red-500 text-xs mt-1">{errors.stock}</span>
                    )}
                </div>

                {/* Size */}
                <div className="flex flex-col">
                    <label htmlFor="size" className="mb-1">
                        Size *
                    </label>

                    <select
                        id="size"
                        name="size"
                        value={formData.size}
                        onChange={onChange}
                        className="border rounded p-2 bg-white"
                    >
                        <option value="">Select Size</option>
                        <option value="Small">Small</option>
                        <option value="Medium">Medium</option>
                        <option value="Large">Large</option>
                    </select>
                    {errors.size && (
                        <span className="text-red-500 text-xs mt-1">{errors.size}</span>
                    )}
                </div>

                {/* Popular */}
                <div className="flex items-center pt-8 gap-2">
                    <input
                        type="checkbox"
                        id="popular"
                        name="popular"
                        checked={formData.popular}
                        onChange={onChange}
                    />

                    <label htmlFor="popular">
                        Popular
                    </label>
                </div>

                

                {/* Offer Section Header */}
                <div className="md:col-span-2 border-t pt-4 mt-2">
                    <h3 className="font-semibold text-lg mb-1">Offer & Pricing Setup</h3>
                    <p className="text-xs text-gray-500 mb-3">
                        Set a discount percentage to automatically activate an offer for this product. Set to 0 if no offer.
                    </p>
                </div>
                
                {/* Discount Percentage Input */}
                <div className="flex flex-col">
                    <label htmlFor="discountPercentage" className="mb-1 text-sm font-medium">
                        Discount Percentage (%)
                    </label>
                    <input
                        type="number"
                        id="discountPercentage"
                        name="discountPercentage"
                        min="0"
                        max="100"
                        value={formData.discountPercentage || ''}
                        onChange={onChange}
                        className="border rounded p-2"
                        placeholder="0"
                    />
                    {errors.discountPercentage && (
                        <span className="text-red-500 text-xs mt-1">{errors.discountPercentage}</span>
                    )}
                </div>

                {/* Calculated Live Discount Preview (Disabled/Read-only) */}
                <div className="flex flex-col">
                    <label className="mb-1 text-sm font-medium text-gray-500">
                        Discount Amount (Auto-calculated)
                    </label>
                    <input
                        type="text"
                        disabled
                        value={`$${discountAmount}`}
                        className="border rounded p-2 bg-gray-100 text-gray-600 font-medium cursor-not-allowed"
                    />
                </div>

                {/* Calculated Final Price Preview (Disabled/Read-only) */}
                <div className="flex flex-col md:col-span-2">
                    <label className="mb-1 text-sm font-medium text-gray-500">
                        Final Price Customer Pays (Auto-calculated)
                    </label>
                    <input
                        type="text"
                        disabled
                        value={`$${finalPrice}`}
                        className="border rounded p-2 bg-gray-100 text-gray-900 font-bold text-lg cursor-not-allowed"
                    />
                </div>
            </div>

            <button
                disabled={isLoading}
                type="submit"
                className="w-full bg-black text-white py-3 mt-6 rounded hover:bg-gray-800 disabled:opacity-50"
            >
                {isLoading ? "Saving..." : "Submit"}
            </button>
        </>
    );
}