import { useFarmStore } from "../../../stores/useFarmStore";
import type { IGood, RaritySeed } from "../../../types/farm";
import { rarityConfig } from "../../../utils/rarity";

type Props = {
  item: IGood;
  handleBuy: (item: IGood) => void;
};

export const GoodCard = ({ item, handleBuy }: Props) => {
  const { coins } = useFarmStore();

  let name = "";
  let img = "";
  let rarity: RaritySeed = "common";
  let growTime = 0;

  if (item.type === "seed") {
    name = 'Cемена "' + item.item?.name + '"' || "Неизвестно";
    img = item.item?.icon || "";
    rarity = item.item?.rarity || "common";
    growTime = item.item?.targetGrowth || 0;
  }

  if (item.type === "bed") {
    name = "Новая грядка";
    img = "/assets/farm/bed.png";
    rarity = "Unique";
  }

  const { label, light, dark, bgLight, bgDark } = rarityConfig[rarity];
  const canBuy = coins >= item.cost;

  return (
    <div
      className={`relative min-h-84 rounded-2xl shadow-lg p-5 flex flex-col items-center text-center transition-all duration-300 
      bg-liner-to-b from-base-100 to-base-200 border border-base-300 hover:shadow-2xl hover:-translate-y-1`}
    >
      <div
        className={`absolute inset-0 rounded-2xl opacity-20 pointer-events-none ${bgLight} dark:${bgDark}`}
      ></div>

      <div className="relative w-24 h-24 mb-3">
        <img
          src={img}
          alt={name}
          className="w-full h-full object-contain drop-shadow-md"
          style={{ imageRendering: "pixelated" }}
        />
        {item.quantity > 1 && (
          <div className="absolute bottom-1 right-1 bg-base-300 bg-opacity-70 text-xs rounded-md px-1.5 py-0.5 font-mono font-semibold">
            ×{item.quantity}
          </div>
        )}
      </div>

      <h3 className="font-semibold text-base-content text-lg">{name}</h3>

      {item.type === "seed" && (
        <p className="text-sm text-gray-500">{growTime} задач до сбора</p>
      )}

      <div
        className={`mt-2 text-xs px-3 py-0.5 rounded-full font-semibold tracking-wide uppercase ${light} dark:${dark} ${bgLight} dark:${bgDark} shadow-sm`}
      >
        {label}
      </div>

      <div className="my-3 flex items-center gap-1 text-sm text-base-content font-medium">
        <span>{item.cost}</span>
        <span className="text-yellow-500">🪙</span>
      </div>

      <button
        onClick={() => handleBuy(item)}
        disabled={!canBuy}
        className={`btn btn-sm mt-auto w-full rounded-xl font-semibold transition-all duration-200 
          ${
            canBuy
              ? "bg-emerald-500 hover:bg-emerald-600 text-white shadow-md"
              : "bg-gray-400 text-gray-200 cursor-not-allowed"
          }`}
      >
        {canBuy ? "Купить" : "Недостаточно монет"}
      </button>
    </div>
  );
};
