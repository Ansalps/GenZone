interface Category {
    id: number;
    category_name: string;
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
        hasOffer: boolean;
        discountPercentage: number;
        discountAmount: number;
        totalDiscountedAmount: number;
    };
    onChange: (
        e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>
    ) => void;
    isLoading: boolean;
    categories: Category[];
}

export default function ProductForm({
    formData,
    onChange,
    isLoading,
    categories,
}: ProductFormProps) {
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
                </div>

                {/* Price */}
                <div className="flex flex-col">
                    <label htmlFor="price" className="mb-1">
                        Price *
                    </label>

                    <input
                        type="number"
                        id="price"
                        name="price"
                        value={formData.price}
                        onChange={onChange}
                        className="border rounded p-2"
                        required
                    />
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
                        value={formData.stock}
                        onChange={onChange}
                        className="border rounded p-2"
                        required
                    />
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

                {/* Has Offer */}
                <div className="flex items-center pt-2 gap-2 md:col-span-2">
                    <input
                        type="checkbox"
                        id="hasOffer"
                        name="hasOffer"
                        checked={formData.hasOffer}
                        onChange={onChange}
                    />

                    <label htmlFor="hasOffer">
                        Has Offer
                    </label>
                </div>

                {/* Offer Fields */}
                {formData.hasOffer && (
                    <>
                        <div className="flex flex-col">
                            <label
                                htmlFor="discountPercentage"
                                className="mb-1"
                            >
                                Discount Percentage *
                            </label>

                            <input
                                type="number"
                                id="discountPercentage"
                                name="discountPercentage"
                                value={formData.discountPercentage}
                                onChange={onChange}
                                className="border rounded p-2"
                                required
                            />
                        </div>

                        <div className="flex flex-col">
                            <label
                                htmlFor="discountAmount"
                                className="mb-1"
                            >
                                Discount Amount *
                            </label>

                            <input
                                type="number"
                                id="discountAmount"
                                name="discountAmount"
                                value={formData.discountAmount}
                                onChange={onChange}
                                className="border rounded p-2"
                                required
                            />
                        </div>

                        <div className="flex flex-col md:col-span-2">
                            <label
                                htmlFor="totalDiscountedAmount"
                                className="mb-1"
                            >
                                Total Discounted Amount *
                            </label>

                            <input
                                type="number"
                                id="totalDiscountedAmount"
                                name="totalDiscountedAmount"
                                value={formData.totalDiscountedAmount}
                                onChange={onChange}
                                className="border rounded p-2"
                                required
                            />
                        </div>
                    </>
                )}
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