import { m, AnimatePresence } from "framer-motion";
import { ArrowLeft } from "lucide-react";
import { Controller, useForm } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { lazy, memo, Suspense, useEffect } from "react";
import clsx from "clsx";
import { useFarmStore } from "../../../stores/useFarmStore";
import { difficultyColor } from "../../../utils/difficultyColor";
import { DifficultySelector } from "../../TaskForm/DifficultySelector";
import { TagSelect } from "../../TaskForm/TagSelect";
import { useEditHabitModalStore } from "../../../stores/editHabitModal";
import { InputTitle } from "../../HabitForm/InputTitle";
import { TextAreaDescription } from "../../HabitForm/TextAreaDescription";
import { SelectPeriod } from "../../HabitForm/SelectPeriod";
import { InputEvery } from "../../HabitForm/InputEvery";
import { PeriodSummary } from "../../HabitForm/PeriodSummary";

const DatePicker = lazy(() => import("../../TaskForm/DayPicker"));

const newTaskSchema = z.object({
  title: z
    .string()
    .min(1, "Название должно содержать минимум 1 символ")
    .max(100, "Название слишком длинное"),

  description: z
    .string()
    .max(500, "Описание не должно превышать 500 символов"),

  difficulty: z.enum(["trifle", "easy", "normal", "hard"]),

  startDate: z.date(),

  period: z.enum(["day", "week", "month", "year"]),

  every: z.number().min(1).max(100),

  tag: z
    .object({
      id: z.number(),
      name: z.string(),
      color: z.string(),
    })
    .optional().nullable(),
});


export type EditHabitFormValues = z.infer<typeof newTaskSchema>;

export const EditHabitModal = memo(() => {
  const { isOpen, close, idHabit } = useEditHabitModalStore();
  const { habits, editHabit, removeHabit } = useFarmStore();

  const currentHabit = habits.find(h => h.id === idHabit)
  
  const {
    handleSubmit,
    reset,
    control,
    register,
    formState: { errors },
  } = useForm<EditHabitFormValues>({
    resolver: zodResolver(newTaskSchema),
    defaultValues: {
      title: currentHabit?.title,
      description: currentHabit?.description,
      difficulty: currentHabit?.difficulty,
      startDate: currentHabit?.startDate || new Date(),
      tag: currentHabit?.tag || null,
      every: currentHabit?.every || 0,
    },
  });

  useEffect(() => {
    reset(currentHabit)
  }, [currentHabit, reset])
  
  const onSubmit = (data: EditHabitFormValues) => {
    editHabit(
      idHabit,
      data.title,
      data.description,
      data.difficulty,
      data.period,
      data.every,
      data.startDate || new Date(),
      data.tag || null
    );
    reset();
    close();
  };

  const handleClose = () => {
    reset();
    close();
  };

  const deleteTask = () => {
    removeHabit(idHabit);
    close();
  };

  return (
    <AnimatePresence>
      {isOpen && (
        <m.form
          onSubmit={handleSubmit(onSubmit)}
          className="fixed inset-x-0 bottom-0 top-0 z-60 flex flex-col bg-base-100 rounded-t-2xl shadow-xl"
          initial={{ y: "100%" }}
          animate={{ y: 0 }}
          exit={{ y: "100%" }}
          transition={{ duration: 0.8, ease: [0.16, 1, 0.3, 1] }}
        >
          <div className={clsx(
            "sticky top-0 flex items-center justify-between gap-2 p-4 rounded-t-2xl",
            difficultyColor(currentHabit?.difficulty)
          )}>
            <button
              type="button"
              onClick={handleClose}
              className={clsx(
                "rounded-full p-1 transition-colors",
                {
                  "text-slate-50 hover:text-slate-200 hover:bg-slate-600": currentHabit?.difficulty === "trifle",
                  "text-emerald-50 hover:text-emerald-200 hover:bg-emerald-600": currentHabit?.difficulty === "easy",
                  "text-sky-50 hover:text-sky-200 hover:bg-sky-600": currentHabit?.difficulty === "normal",
                  "text-orange-50 hover:text-orange-200 hover:bg-orange-600": currentHabit?.difficulty === "hard",
                } 
              )}
            >
              <ArrowLeft />
            </button>
            <span className={clsx(
              "font-semibold",
              {
                "text-slate-50": currentHabit?.difficulty === "trifle",
                "text-emerald-50": currentHabit?.difficulty === "easy",
                "text-sky-50": currentHabit?.difficulty === "normal",
                "text-orange-50": currentHabit?.difficulty === "hard",
              }
            )}>
              Редактировать дело
            </span>

            <div className="flex gap-2">
              <button
                type="button"
                onClick={ deleteTask }
                className={clsx(
                  "btn btn-sm font-mono outline-none border-0",
                  {
                    "bg-slate-50 text-slate-700 hover:bg-slate-100": currentHabit?.difficulty === "trifle",
                    "bg-emerald-50 text-emerald-700 hover:bg-emerald-100": currentHabit?.difficulty === "easy",
                    "bg-sky-50 text-sky-700 hover:bg-sky-100": currentHabit?.difficulty === "normal",
                    "bg-orange-50 text-orange-700 hover:bg-orange-100": currentHabit?.difficulty === "hard",
                  } 
                )}
              >
                Удалить
              </button>
              <button
                type="submit"
                className={clsx(
                  "btn btn-sm font-mono outline-none border-0",
                  {
                    "bg-slate-50 text-slate-700 hover:bg-slate-100": currentHabit?.difficulty === "trifle",
                    "bg-emerald-50 text-emerald-700 hover:bg-emerald-100": currentHabit?.difficulty === "easy",
                    "bg-sky-50 text-sky-700 hover:bg-sky-100": currentHabit?.difficulty === "normal",
                    "bg-orange-50 text-orange-700 hover:bg-orange-100": currentHabit?.difficulty === "hard",
                  } 
                )}
              >
                Сохранить
              </button>
            </div>

          </div>

          <div className="flex-1 overflow-y-auto">
            <div className={clsx(
              "flex flex-col gap-2 items-center p-4 w-full",
              difficultyColor(currentHabit?.difficulty)
            )}>
              <InputTitle
                register={register}
                difficulty={currentHabit?.difficulty}
              />
              
              <TextAreaDescription
                register={register}
                difficulty={currentHabit?.difficulty}
              />
            </div>

            <div className="p-4 flex gap-2 w-full">
              <div className="flex flex-col gap-1 w-full">
                <p>Сложность</p>
                <Controller
                  name="difficulty"
                  control={control}
                  render={({ field: { value, onChange } }) => (
                    <DifficultySelector value={value} onChange={onChange} />
                  )}
                />
                {errors.difficulty && (
                  <p className="text-error text-sm mt-1">
                    {errors.difficulty.message}
                  </p>
                )}
              </div>
            </div>

            <div className="p-4 flex gap-2 w-full">
              <div className="flex flex-col gap-2 w-full">
                <p>Дата начало отсчёта</p>
                <Controller
                  name="startDate"
                  control={control}
                  render={({ field: { value, onChange } }) => (
                    <Suspense>
                      <DatePicker value={value} onChange={onChange} />
                    </Suspense>
                  )}
                />
                {errors.startDate && (
                  <p className="text-error text-sm mt-1">
                    {errors.startDate.message}
                  </p>
                )}
              </div>
            </div>

            <div className="flex gap-8 items-center p-4 pb-0 w-full">
              <SelectPeriod
                register={register}
                difficulty={currentHabit?.difficulty}
              />
              
              <InputEvery
                register={register}
                difficulty={currentHabit?.difficulty}
              />
                              {errors.every && (
                  <p className="text-error text-sm mt-1">
                    {errors.every.message}
                  </p>
                )}
            </div>
            <Controller
              name="period"
              control={control}
              render={({ field: { value: period } }) => (
                <Controller
                name="every"
                control={control}
                render={({ field: { value: every } }) => (
                  <PeriodSummary period={period} every={every} />
                )}
                />
              )}
            />

            <div className="p-4 flex gap-2 w-full">
              <div className="flex flex-col gap-2 w-full">
                <p>Тег</p>
                <Controller
                  name="tag"
                  control={control}
                  render={({ field: { value, onChange } }) => (
                    <Suspense>
                      <TagSelect
                        selectedTag={value || null}
                        difficulty={currentHabit?.difficulty}
                        onSelect={onChange}
                      />
                    </Suspense>
                  )}
                />
              </div>
            </div>

          </div>


        </m.form>
      )}
    </AnimatePresence>
  );
});
