export default function Login(){
    return (
        /* 1. Added flex-col to stack items vertically */
        /* 2. Added h-screen so the box stretches from top to bottom of the screen */
        <div className="flex flex-col justify-center items-center h-screen gap-4">
            
            <div className="flex flex-col">
                <label htmlFor="email">Email</label>
                <input type="text" id="email" className="border-black border p-1"/>
            </div>

            <div className="flex flex-col">
                <label htmlFor="password">Password</label>
                <input type="text" id="password" className="border-black border p-1"/>
            </div>

        </div>
    )
}
