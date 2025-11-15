import { Flower2 } from "lucide-react";
import { useFarmStore } from "../../stores/useFarmStore";
import { GoodCard } from "./ui/GoodCard";
import type { IGood } from "../../types/farm";
import { fetchBuyItem } from "../../api/shop";
import { showErrorToast } from "../../utils/showErrorToast";



export const Shop = () => {
  const { shopItem, decreaseShopItem, addSeedToInventory, unlockNextBed, coins, setCoins } = useFarmStore();

  const handleBuy = async (item: IGood) => {
    if (coins < item.cost) return;
    
    const status = await fetchBuyItem(item.id)

    if (status !== 200) {
      showErrorToast("Не удалось купить товар");
      return;
    }

    const newCoins = coins - item.cost;

    switch (item.type) {
      case "seed":
        setCoins(newCoins);
        if (item.item) {
          addSeedToInventory(item.item);
        }
        break;
      case "bed":
        setCoins(newCoins);
        unlockNextBed();
    }

    decreaseShopItem(item.type, item.id);
  }

  return (
    <div className="">
      <div className="px-4 py-6 scrollbar-hide">
        <h2 className="text-2xl font-bold text-center mb-6 text-emerald-700">
          Магазин
        </h2>

        {shopItem.length === 0 ? (
          <div className="flex flex-col items-center justify-center h-full text-center space-y-4 py-16">
            <Flower2 strokeWidth={1} size={80} className="text-base-content/80" />
            <p className="text-xl font-semibold text-base-content/60">
              Магазин пуст
            </p>
            <p className="text-sm text-base-content/50 max-w-xs">
              Заходите завтра!
            </p>
          </div>
        ) : (
          <ul className="grid grid-cols-2 gap-4 pb-4 h-[calc(100dvh-200px)] overflow-y-auto scrollbar-hide">
            {shopItem.map((item) => (
              <li key={item.id} >
                <GoodCard item={item} handleBuy={handleBuy} />
              </li>
            ))}
          </ul>
        )}
      </div>
    </div>
  );
};
