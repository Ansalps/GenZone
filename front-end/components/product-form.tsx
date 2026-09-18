import { ProductFormData } from "@/types/productFormData"

interface Category {
    id: number
    category_name: string
}

export interface FormErrors {
    categoryName?: string
    productName?: string
    productDescription?: string
    productImage?: string
    price?: string
    stock?: string
    size?: string
    inventory?: string
    discountPercentage?: string
}

interface ProductFormProps {
    formData: ProductFormData

    onChange: (
        e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>
    ) => void

    onInventoryChange?: (size: string, value: number) => void

    onImageChange: (
        e: React.ChangeEvent<HTMLInputElement>
    ) => void

    isLoading: boolean
    categories: Category[]
    errors?: FormErrors

    // Used only when editing an existing product
    existingImageUrl?: string
}

export default function ProductForm({
    formData,
    onChange,
    onInventoryChange,
    onImageChange,
    isLoading,
    categories,
    errors = {},
    existingImageUrl,
}: ProductFormProps) {

    const price = Number(formData.price) || 0
    const inventoryConfig = [
        { size: 'Small', label: 'S' },
        { size: 'Medium', label: 'M' },
        { size: 'Large', label: 'L' },
    ]
    const discountPercent =
        Number(formData.discountPercentage) || 0

    const discountAmount =
        price > 0 && discountPercent > 0
            ? ((price * discountPercent) / 100).toFixed(2)
            : "0.00"

    const finalPrice =
        price > 0 && discountPercent > 0
            ? (
                  price - Number(discountAmount)
              ).toFixed(2)
            : price.toFixed(2)

    return (
        <>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-5">

                {/* Category */}
                <div className="flex flex-col">
                    <label htmlFor="categoryName" className="mb-2 text-sm font-medium text-slate-200">Category *</label>

                    <select id="categoryName" name="categoryName" value={formData.categoryName} onChange={onChange} className="rounded-xl border border-white/10 bg-slate-900/50 px-4 py-3 text-white outline-none">
                        <option value="">Select Category</option>

                        {categories.map((category) => (
                            <option key={category.id} value={category.category_name}>
                                {category.category_name}
                            </option>
                        ))}
                    </select>

                    {errors.categoryName && <span className="text-rose-400 text-xs mt-1">{errors.categoryName}</span>}
                </div>

                {/* Product Name */}
                <div className="flex flex-col">
                    <label htmlFor="productName" className="mb-2 text-sm font-medium text-slate-200">Product Name *</label>

                    <input type="text" id="productName" name="productName" value={formData.productName} onChange={onChange} className="rounded-xl border border-white/10 bg-slate-900/50 px-4 py-3 text-white outline-none" required />

                    {errors.productName && <span className="text-rose-400 text-xs mt-1">{errors.productName}</span>}
                </div>

                {/* Description */}
                <div className="flex flex-col">
                    <label htmlFor="productDescription" className="mb-2 text-sm font-medium text-slate-200">Description *</label>

                    <input type="text" id="productDescription" name="productDescription" value={formData.productDescription} onChange={onChange} className="rounded-xl border border-white/10 bg-slate-900/50 px-4 py-3 text-white outline-none" required />

                    {errors.productDescription && <span className="text-rose-400 text-xs mt-1">{errors.productDescription}</span>}
                </div>

                {/* Product Image */}
                <div className="flex flex-col">
                    <label htmlFor="productImage" className="mb-2 text-sm font-medium text-slate-200">Product Image *</label>

                    <input type="file" id="productImage" name="productImage" accept="image/*" onChange={onImageChange} className="rounded-lg border border-white/10 bg-slate-900/50 px-3 py-2 text-sm text-slate-200 cursor-pointer" />

                    {errors.productImage && <span className="text-rose-400 text-xs mt-1">{errors.productImage}</span>}

                    {/* Existing image - only appears on Edit */}
                    {existingImageUrl && (
                        <div className="mt-3">
                            <p className="text-sm text-slate-300 mb-2">Current image:</p>

                            <img src={existingImageUrl} alt="Current product" className="w-32 h-32 object-cover rounded border" />
                        </div>
                    )}
                </div>

                {/* Price */}
                <div className="flex flex-col">
                    <label htmlFor="price" className="mb-2 text-sm font-medium text-slate-200">Original Price ($) *</label>

                    <input type="number" id="price" name="price" step="0.01" min="0" value={formData.price || ""} onChange={onChange} className="rounded-xl border border-white/10 bg-slate-900/50 px-4 py-3 text-white outline-none" required />

                    {errors.price && <span className="text-rose-400 text-xs mt-1">{errors.price}</span>}
                </div>

                {!onInventoryChange ? (
                    <>
                        {/* Stock */}
                        <div className="flex flex-col">
                            <label htmlFor="stock" className="mb-2 text-sm font-medium text-slate-200">Stock *</label>

                            <input type="number" id="stock" name="stock" min="0" value={formData.stock || ""} onChange={onChange} className="rounded-xl border border-white/10 bg-slate-900/50 px-4 py-3 text-white outline-none" required />

                            {errors.stock && <span className="text-rose-400 text-xs mt-1">{errors.stock}</span>}
                        </div>

                        {/* Size */}
                        <div className="flex flex-col">
                            <label htmlFor="size" className="mb-2 text-sm font-medium text-slate-200">Size *</label>

                            <select id="size" name="size" value={formData.size} onChange={onChange} className="rounded-xl border border-white/10 bg-slate-900/50 px-4 py-3 text-white outline-none">
                                <option value="">Select Size</option>
                                <option value="Small">Small</option>
                                <option value="Medium">Medium</option>
                                <option value="Large">Large</option>
                            </select>

                            {errors.size && <span className="text-rose-400 text-xs mt-1">{errors.size}</span>}
                        </div>
                    </>
                ) : (
                    <div className="md:col-span-2">
                        <div className="mb-2 flex items-center justify-between">
                            <label className="text-sm font-medium text-slate-200">Inventory</label>
                        </div>
                        <div className="grid grid-cols-3 gap-3">
                            {inventoryConfig.map((item) => (
                                <div key={item.size} className="rounded-2xl border border-white/10 bg-slate-950/40 p-3">
                                    <div className="mb-2 text-center text-xs font-semibold uppercase tracking-[0.2em] text-slate-400">
                                        {item.label}
                                    </div>
                                    <input
                                        type="number"
                                        min="0"
                                        value={formData.inventory?.[item.size] ?? 0}
                                        onChange={(event) => onInventoryChange?.(item.size, Number(event.target.value) || 0)}
                                        className="w-full rounded-xl border border-white/10 bg-slate-900/50 px-3 py-2 text-center text-white outline-none"
                                    />
                                </div>
                            ))}
                        </div>
                        {errors.inventory && <span className="mt-2 block text-xs text-rose-400">{errors.inventory}</span>}
                    </div>
                )}

                {/* Popular */}
                <div className="flex items-center pt-8 gap-3">
                    <input type="checkbox" id="popular" name="popular" checked={formData.popular} onChange={onChange} className="h-4 w-4 rounded border-white/10 bg-slate-900/50" />

                    <label htmlFor="popular" className="text-sm text-slate-200">Popular</label>
                </div>

                {/* Offer Section */}
                <div className="md:col-span-2 border-t pt-4 mt-2">
                    <h3 className="font-semibold text-lg mb-1 text-white">Offer & Pricing Setup</h3>

                    <p className="text-xs text-slate-300 mb-3">Set a discount percentage to automatically activate an offer for this product. Set to 0 if no offer.</p>
                </div>

                {/* Discount Percentage */}
                <div className="flex flex-col">
                    <label
                        htmlFor="discountPercentage"
                        className="mb-1 text-sm font-medium"
                    >
                        Discount Percentage (%)
                    </label>

                    <input type="number" id="discountPercentage" name="discountPercentage" min="0" max="100" value={formData.discountPercentage || ""} onChange={onChange} className="rounded-xl border border-white/10 bg-slate-900/50 px-4 py-3 text-white outline-none" placeholder="0" />

                    {errors.discountPercentage && <span className="text-rose-400 text-xs mt-1">{errors.discountPercentage}</span>}
                </div>

                {/* Discount Amount */}
                <div className="flex flex-col">
                    <label className="mb-1 text-sm font-medium text-gray-500">
                        Discount Amount (Auto-calculated)
                    </label>

                    <input type="text" disabled value={`$${discountAmount}`} className="rounded-xl border border-white/10 px-4 py-2 bg-slate-900/30 text-slate-300 font-medium cursor-not-allowed" />
                </div>

                {/* Final Price */}
                <div className="flex flex-col md:col-span-2">
                    <label className="mb-1 text-sm font-medium text-gray-500">
                        Final Price Customer Pays (Auto-calculated)
                    </label>

                    <input type="text" disabled value={`$${finalPrice}`} className="rounded-xl border border-white/10 px-4 py-3 bg-slate-900/30 text-white font-bold text-lg cursor-not-allowed" />
                </div>
            </div>

            <button disabled={isLoading} type="submit" className="w-full rounded-xl bg-gradient-to-r from-emerald-500 to-teal-500 px-4 py-3 text-sm font-semibold text-white shadow-lg shadow-emerald-500/20 mt-6 disabled:opacity-60">
                {isLoading ? 'Saving...' : 'Submit'}
            </button>
        </>
    )
}