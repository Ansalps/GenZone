interface CategoryFormProps {
    formData: {
        categoryName: string;
        categoryDescription: string;
        categoryImage: File | null;
        categoryImageUrl?: string;
    };

    onChange: (e: React.ChangeEvent<HTMLInputElement>) => void;

    onImageChange: (e: React.ChangeEvent<HTMLInputElement>) => void;

    isLoading: boolean;

    isEdit?: boolean;
}

export default function CategoryForm({
    formData,
    onChange,
    onImageChange,
    isLoading,
    isEdit = false,
}: CategoryFormProps) {
    return (
        <>
            {/* Category Name */}
            <div className="flex flex-col mb-4">
                <label htmlFor="categoryName" className="mb-2 text-sm font-medium text-slate-200">
                    Category name *
                </label>

                <input
                    type="text"
                    id="categoryName"
                    name="categoryName"
                    value={formData.categoryName}
                    onChange={onChange}
                    className="rounded-xl border border-white/10 bg-slate-900/50 px-4 py-3 text-white outline-none focus:border-cyan-400 focus:ring-2 focus:ring-cyan-400/20"
                    required
                />
            </div>

            {/* Category Description */}
            <div className="flex flex-col mb-4">
                <label htmlFor="categoryDescription" className="mb-2 text-sm font-medium text-slate-200">
                    Category description *
                </label>

                <input
                    type="text"
                    id="categoryDescription"
                    name="categoryDescription"
                    value={formData.categoryDescription}
                    onChange={onChange}
                    className="rounded-xl border border-white/10 bg-slate-900/50 px-4 py-3 text-white outline-none focus:border-cyan-400 focus:ring-2 focus:ring-cyan-400/20"
                    required
                />
            </div>

            {/* Existing Image - Edit */}
            {isEdit && formData.categoryImageUrl && (
                <div className="mb-4">
                    <label className="block mb-2 font-medium">
                        Current image
                    </label>

                    <img
                        src={formData.categoryImageUrl}
                        alt={formData.categoryName}
                        className="w-40 h-40 object-cover rounded border"
                    />
                </div>
            )}

            {/* Image */}
            <div className="flex flex-col mb-4">
                <label htmlFor="categoryImage" className="mb-2 text-sm font-medium text-slate-200">
                    {isEdit ? 'Change category image' : 'Category image *'}
                </label>

                <input
                    type="file"
                    id="categoryImage"
                    name="categoryImage"
                    accept="image/*"
                    onChange={onImageChange}
                    className="rounded-lg border border-white/10 bg-slate-900/50 px-3 py-2 text-sm text-slate-200 cursor-pointer"
                    required={!isEdit}
                />

                {formData.categoryImage && (
                    <p className="text-sm text-slate-300 mt-2">Selected: {formData.categoryImage.name}</p>
                )}
            </div>

            {/* Submit */}
            <button
                disabled={isLoading}
                type="submit"
                className="w-full rounded-xl bg-gradient-to-r from-cyan-500 to-violet-500 px-4 py-3 text-sm font-semibold text-white shadow-lg shadow-cyan-500/20 mt-2 disabled:opacity-60"
            >
                {isLoading ? 'Saving...' : 'Submit'}
            </button>
        </>
    );
}