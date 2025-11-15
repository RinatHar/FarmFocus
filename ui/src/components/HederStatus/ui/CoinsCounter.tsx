import { useFarmStore } from "../../../stores/useFarmStore";


export const CoinsCounter = () => {
  const { coins } = useFarmStore();


  const formatCoins = (num: number) => {
    if (num >= 1_000_000) {
      return `${(num / 1_000_000).toFixed(1).replace(/\.0$/, "")}M`;
    }
    if (num >= 1_000) {
      return `${(num / 1_000).toFixed(1).replace(/\.0$/, "")}K`;
    }
    return num.toString();
  };

  return (
    <div className="flex items-center gap-1.5">
        <p className="font-mono text-amber-500 dark:text-amber-400">
          {formatCoins(coins)}
        </p>
        🪙
      </div>
  );
};