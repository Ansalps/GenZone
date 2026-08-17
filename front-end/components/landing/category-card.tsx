interface Category {
    id: number;
    category_name: string;
    category_description: string;
    category_image_url: string;
}

interface CategoryCardProps {
    category: Category;
}

export default function CategoryCard({
    category,
}: CategoryCardProps) {
    return (
        <div className="bg-white rounded-xl shadow-md overflow-hidden hover:shadow-xl transition duration-300 cursor-pointer">

            {/* Image */}
            <div className="h-48 overflow-hidden">
                <img
                    src={category.category_image_url}
                    alt={category.category_name}
                    className="w-full h-full object-cover hover:scale-105 transition duration-300"
                />
            </div>

            {/* Content */}
            <div className="p-4">
                <h3 className="text-lg font-semibold text-gray-800">
                    {category.category_name}
                </h3>

                <p className="text-sm text-gray-500 mt-2 line-clamp-2">
                    {category.category_description}
                </p>
            </div>
        </div>
    );
}