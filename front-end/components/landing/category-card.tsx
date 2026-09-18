interface Category {
    id: number;
    category_name: string;
    category_description: string;
    category_image_url: string;
}

interface CategoryCardProps {
    category: Category;
}

export default function CategoryCard({ category }: CategoryCardProps) {
    return (
        <div className="group overflow-hidden rounded-[1.5rem] border border-white/10 bg-slate-900/80 shadow-lg shadow-slate-950/20 transition duration-300 hover:-translate-y-1 hover:border-cyan-400/40 hover:shadow-cyan-950/20">
            <div className="relative h-52 overflow-hidden">
                <img
                    src={category.category_image_url}
                    alt={category.category_name}
                    className="h-full w-full object-cover transition duration-500 group-hover:scale-105"
                />
                <div className="absolute inset-0 bg-gradient-to-t from-slate-950 via-slate-950/20 to-transparent" />
            </div>

            <div className="space-y-3 p-5">
                <div className="flex items-center justify-between gap-3">
                    <h3 className="text-lg font-semibold text-white">{category.category_name}</h3>
                    <span className="rounded-full border border-cyan-400/30 bg-cyan-500/10 px-2 py-1 text-[10px] font-semibold uppercase tracking-[0.2em] text-cyan-200">
                        New
                    </span>
                </div>

                <p className="line-clamp-2 text-sm text-slate-300">
                    {category.category_description}
                </p>
            </div>
        </div>
    );
}