import { Toaster } from "sonner";
import"./globals.css"
export default function fsljfs({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html>
      <body>
        {children}
         <Toaster
                    position="top-right"
                    richColors
                    closeButton
          />
      </body>
     
    </html>
  );
}
