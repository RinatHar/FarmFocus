import { CalendarDays, CircleCheckBig, NotebookPen, ShoppingBag } from "lucide-react";
import type { SwiperPage } from "../../App";
import { ButtonAddTask } from "./ui/ButtonAddTask";
import clsx from "clsx";
import { useFarmStore } from "../../stores/useFarmStore";
import { BottomMenuButton } from "./ui/BottomMenuButton";

type Props = {
  onChangePage: (page: SwiperPage) => void;
  currentPage?: SwiperPage;
};

export const BottomMenu = ({ onChangePage, currentPage = "todo-list" }: Props) => {
  const { isDrought } = useFarmStore();

  return (
    <menu
      className={clsx(
        "fixed flex gap-2 xs:gap-6 justify-between bottom-0 bg-base-200 border-t",
        "left-0 right-0 z-50 bg-emerald-60 p-1 px-3 rounded-t-3xl shadow-lg",
        {
          "border-lime-400": !isDrought,
          "border-amber-400": isDrought,
        }
      )}
    >
      <div className="flex gap-1 xs:gap-8 justify-around w-full">
        <BottomMenuButton
          icon={CalendarDays}
          label="Календарь"
          page="calendar"
          currentPage={currentPage}
          isDrought={isDrought}
          onClick={() => onChangePage("calendar")}
        />

        <BottomMenuButton
          icon={NotebookPen}
          label="Дела"
          page="habit-list"
          currentPage={currentPage}
          isDrought={isDrought}
          onClick={() => onChangePage("habit-list")}
        />
      </div>

      <ButtonAddTask currentPage={currentPage} />

      <div className="flex gap-1 xs:gap-8 justify-around w-full">
        <BottomMenuButton
          icon={CircleCheckBig}
          label="Задачи"
          page="todo-list"
          currentPage={currentPage}
          isDrought={isDrought}
          onClick={() => onChangePage("todo-list")}
        />

        <BottomMenuButton
          icon={ShoppingBag}
          label="Магазин"
          page="shop"
          currentPage={currentPage}
          isDrought={isDrought}
          onClick={() => onChangePage("shop")}
        />
      </div>
    </menu>
  );
};