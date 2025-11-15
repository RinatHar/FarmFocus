import { Sparkles } from "lucide-react";
import { toast } from "sonner";

export const showRewardToast = (xp = 0, coins = 0) => {
  if (xp <= 0 && coins <= 0) return;

  toast.custom(() => (
    <div className="flex items-center gap-3 px-4 py-3 rounded-xl shadow-lg bg-base-100">
      <div className="flex flex-col">
        <div className="text-sm font-medium opacity-90">Получено:</div>
        <div className="flex items-center gap-3 mt-0.5">
          {xp > 0 && (
            <div className="flex items-center gap-1">
              <span className="font-semibold">{xp} опыта</span>
              <Sparkles className="w-5 h-5 text-emerald-400" />
            </div>
          )}
          {coins > 0 && (
            <div className="flex items-center gap-1">
              <span className="font-semibold">{coins} монет</span>
              🪙
            </div>
          )}
        </div>
      </div>
    </div>
  ));
};