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
                <label
                    htmlFor="categoryName"
                    className="mb-1 font-medium"
                >
                    Category name *
                </label>

                <input
                    type="text"
                    id="categoryName"
                    name="categoryName"
                    value={formData.categoryName}
                    onChange={onChange}
                    className="border border-black p-2 rounded"
                    required
                />
            </div>

            {/* Category Description */}
            <div className="flex flex-col mb-4">
                <label
                    htmlFor="categoryDescription"
                    className="mb-1 font-medium"
                >
                    Category description *
                </label>

                <input
                    type="text"
                    id="categoryDescription"
                    name="categoryDescription"
                    value={formData.categoryDescription}
                    onChange={onChange}
                    className="border border-black p-2 rounded"
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
                <label
                    htmlFor="categoryImage"
                    className="mb-1 font-medium"
                >
                    {isEdit
                        ? 'Change category image'
                        : 'Category image *'}
                </label>

                <input
                    type="file"
                    id="categoryImage"
                    name="categoryImage"
                    accept="image/*"
                    onChange={onImageChange}
                    className="border border-black p-2 rounded cursor-pointer"
                    required={!isEdit}
                />

                {formData.categoryImage && (
                    <p className="text-sm text-gray-500 mt-1">
                        Selected: {formData.categoryImage.name}
                    </p>
                )}
            </div>

            {/* Submit */}
            <button
                disabled={isLoading}
                type="submit"
                className="bg-black text-white p-2 mt-2 rounded disabled:bg-gray-400 cursor-pointer"
            >
                {isLoading ? 'Saving...' : 'Submit'}
            </button>
        </>
    );
}