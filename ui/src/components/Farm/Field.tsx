import { useMemo } from "react";
import { useFarmStore } from "../../stores/useFarmStore";
import { PlantSelectModal } from "../Modal/Farm/PlantSelectModal";
import { Bed } from "./Bed";

export const Field = () => {
  const { field, rows, cols, isDrought } = useFarmStore();

  const beds = useMemo(
    () =>
      field.map((bed) => (
        <Bed
          key={bed.id}
          id={bed.id}
          plant={bed.plant}
          isLock={bed.isLock}
        />
      )),
    [field]
  );

  return (
    <>
      <div className="relative w-full max-w-80 md:max-w-lg mx-auto aspect-square overflow-hidden select-none">
        {/* Фон поля */}
        <img
          src="/assets/farm/pole.png"
          alt="farm"
          className={`w-full h-full object-cover transition-all duration-1000 ${
            isDrought ? "brightness-75 saturate-50" : ""
          }`}
          style={{ imageRendering: "pixelated" }}
        />

        {/* Сетка грядок */}
        <div
          className="absolute inset-0 grid py-8 px-4"
          style={{
            gridTemplateColumns: `repeat(${cols}, 1fr)`,
            gridTemplateRows: `repeat(${rows}, 1fr)`,
          }}
        >
          {beds}
        </div>

        {isDrought && (
          <div className="absolute inset-0 pointer-events-none">
            <div className="absolute inset-0 bg-yellow-600 rounded-3xl opacity-20" />
            <div className="absolute inset-0 bg-liner-to-t rounded-3xl from-orange-500/10 to-transparent animate-pulse" />
          </div>
        )}
      </div>

      <PlantSelectModal />
    </>
  );
};