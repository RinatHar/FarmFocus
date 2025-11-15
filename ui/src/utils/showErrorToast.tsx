
import { toast } from "sonner";

export const showErrorToast = (error: string) => {
  toast.custom(() => (
    <div className="flex items-center gap-3 px-4 py-3 rounded-xl shadow-lg bg-base-100">
      <div className="flex flex-col">
        <div className="text-sm font-medium opacity-90">Ошибка:</div>
        <div className="flex items-center gap-3 mt-0.5 text-sm">
          {error}
        </div>
      </div>
    </div>
  ));
};