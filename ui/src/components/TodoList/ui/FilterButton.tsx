// components/FilterButton.tsx
import { Filter } from "lucide-react";
import clsx from "clsx";
import { useFilterModalStore } from "../../../stores/useFilterModalStore";
import { useTaskFiltersStore } from "../../../stores/useTaskFiltersStore";
import { useFarmStore } from "../../../stores/useFarmStore";


export const FilterButton = () => {
  const { isDrought } = useFarmStore();
  const { open } = useFilterModalStore();
  const { selectedDifficulties } = useTaskFiltersStore();

  const hasActiveFilters = selectedDifficulties.length > 0;

  return (
    <button
      onClick={open}
      className={clsx(
        "btn btn-sm btn-outline flex items-center gap-1",
        { 
          "text-emerald-600": !isDrought,
          "text-amber-600": isDrought,
          "btn-primary": hasActiveFilters,
        },
      )}
    >
      <Filter className="w-4 h-4" />
      Фильтры
      {hasActiveFilters && (
        <span className="badge badge-xs badge-primary ml-1">
          {selectedDifficulties.length}
        </span>
      )}
    </button>
  );
};