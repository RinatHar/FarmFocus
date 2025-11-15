import { useFarmStore } from "../../../stores/useFarmStore";

const TOTAL_SEGMENTS = 8;

export const XPBar = () => {
  const { currentLevelXp, countXpLvLUp, currentLevel, isDrought } = useFarmStore();

  const progress = Math.min(currentLevelXp / countXpLvLUp, 1);
  const filledSegments = progress * TOTAL_SEGMENTS;

  const textColor = isDrought
    ? "text-amber-700 dark:text-amber-400"
    : "text-emerald-800 dark:text-emerald-500";

  const bgColor = isDrought ? "bg-amber-200" : "bg-emerald-200";
  const fillColor = isDrought ? "bg-amber-500" : "bg-emerald-400";

  return (
    <div className="flex flex-col gap-1 pb-2 w-fit">
      <div className="flex justify-between">
        <p className={`font-mono text-xs ${textColor}`}>УРОВЕНЬ {currentLevel}</p>
        <p className={`font-mono text-xs ${textColor}`}>
          {currentLevelXp} / {countXpLvLUp}
        </p>
      </div>

      <div className="w-full flex gap-1">
        {Array.from({ length: TOTAL_SEGMENTS }, (_, index) => {
          const segmentProgress = Math.min(Math.max(filledSegments - index, 0), 1);
          return (
            <div
              key={index}
              className={`relative w-6 h-3 rounded-sm ${bgColor} overflow-hidden`}
            >
              <div
                className={`absolute top-0 left-0 h-full ${fillColor} transition-all duration-300`}
                style={{ width: `${segmentProgress * 100}%` }}
              />
            </div>
          );
        })}
      </div>
    </div>
  );
};