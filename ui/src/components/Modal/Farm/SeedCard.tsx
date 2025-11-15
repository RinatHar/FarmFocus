import type { ISeed } from "../../../types/farm";
import { rarityConfig } from "../../../utils/rarity";


interface SeedCardProps {
  seed: ISeed;
  onSelect: () => void;
}


export const SeedCard = ({ seed, onSelect }: SeedCardProps) => {
  const { label, light, dark, bgLight, bgDark } = rarityConfig[seed.rarity];

  return (
    <button
      onClick={onSelect}
      className="group relative flex flex-col items-center justify-center p-4 rounded-xl bg-base-200 dark:bg-base-300 shadow-md hover:bg-base-300 dark:hover:bg-base-200 transition-all duration-200 overflow-hidden"
    >
      {/* Полоса редкости сверху */}
      <div
        className={`absolute top-0 left-0 right-0 h-1 ${bgLight} dark:${bgDark} opacity-70 group-hover:opacity-100 transition-opacity`}
      />

      {/* Иконка */}
      <div className="relative z-10">
        <img
          src={seed.icon}
          alt={seed.name}
          className="w-14 h-14 object-contain"
          style={{ imageRendering: "pixelated" }}
        />
      </div>

      {/* Название */}
      <div className="flex items-center gap-1">
        <span className="mt-2 text-sm font-medium text-base-content">{seed.name}</span>
        {
          seed.quantity > 1 
          ? <div className="flex items-center mt-2.5 gap-1 text-sm">
            <span className="opacity-60">×</span>
            <span className="font-mono font-semibold text-base-content">{seed.quantity}</span>
          </div> : null
        }
      </div>

      <div className={`mt-1 px-2 py-0.5 rounded-full text-xs font-semibold ${light} dark:${dark} ${bgLight} dark:${bgDark}`}>
        {label}
      </div>

      <span className="mt-0.5 text-xs opacity-60">{seed.targetGrowth} задач</span>

    </button>
  );
};