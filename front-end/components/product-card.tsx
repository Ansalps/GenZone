export default function ProductCard(){
    return (
        <div className="w-70 h-100 border-2 border-b-blue-700">
           
            <label htmlFor="">name</label>
            <label htmlFor="">In stock</label>
            <label htmlFor="">price</label>
            <button className="cursor-pointer">Add to Cart</button>

        </div>
    )
}