import { m, AnimatePresence } from "framer-motion";
import { ArrowLeft } from "lucide-react";
import { Controller, useForm } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { useFarmStore } from "../../../stores/useFarmStore";
import { DifficultySelector } from "../../TaskForm/DifficultySelector";
import { TagSelect } from "../../TaskForm/TagSelect";
import { lazy, memo, Suspense } from "react";
import { useAddHabitModalStore } from "../../../stores/addHabitModal";
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

  every: z.number(),

  tag: z
    .object({
      id: z.number(),
      name: z.string(),
      color: z.string(),
    })
    .optional().nullable(),
});

export type NewHabitFormValues = z.infer<typeof newTaskSchema>;

export const AddHabitModal = memo(() => {
  const { isOpen, close } = useAddHabitModalStore();
  const { addHabit } = useFarmStore();

  const {
    handleSubmit,
    reset,
    control,
    register,
    formState: { errors },
  } = useForm<NewHabitFormValues>({
    resolver: zodResolver(newTaskSchema),
    defaultValues: {
      title: "",
      description: "",
      difficulty: "trifle",
      tag: undefined,
      every: 1,
      period: "day",
      startDate: new Date(),
    },
  });

  const onSubmit = (data: NewHabitFormValues) => {
    addHabit(
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

  return (
    <AnimatePresence>
      {isOpen && (
        <m.form
          onSubmit={handleSubmit(onSubmit)}
          className="fixed inset-x-0 bottom-0 top-0 z-50 flex flex-col bg-base-100 rounded-t-2xl shadow-xl"
          initial={{ y: "100%" }}
          animate={{ y: 0 }}
          exit={{ y: "100%" }}
          transition={{ duration: 0.8, ease: [0.16, 1, 0.3, 1] }}
        >
          {/* Header */}
          <div className="sticky top-0 bg-emerald-600 flex items-center justify-between gap-2 p-4 rounded-t-2xl">
            <button
              type="button"
              onClick={handleClose}
              className="text-emerald-50 hover:text-emerald-200 hover:bg-emerald-700 rounded-full p-1 transition-colors"
            >
              <ArrowLeft />
            </button>
            <span className="font-semibold text-emerald-50">Добавить дело</span>
            <button
              type="submit"
              className="btn btn-sm font-mono bg-emerald-50 text-emerald-700 outline-none border-0 hover:bg-emerald-100"
            >
              СОЗДАТЬ
            </button>
          </div>

          <div className="flex-1 overflow-y-auto">

            <div className="bg-emerald-600 flex flex-col gap-2 items-center p-4 w-full">
              <InputTitle autoFocus register={register} />
              <TextAreaDescription register={register} />
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
                  <p className="text-error text-sm mt-1">{errors.difficulty.message}</p>
                )}
              </div>
            </div>

            <div className="p-4 flex gap-2 w-full">
              <div className="flex flex-col gap-2 w-full">
                <p>Дата выполнения</p>
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
                  <p className="text-error text-sm mt-1">{errors.startDate.message}</p>
                )}
              </div>
            </div>

            <div className="flex gap-8 items-center p-4 pb-0 w-full">
              <SelectPeriod register={register} />
              <InputEvery register={register} />
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
                        selectedTag={value}
                        onSelect={onChange}
                      />
                    </Suspense>
                  )}
                />
                {errors.tag && (
                  <p className="text-error text-sm mt-1">Выберите тег</p>
                )}
              </div>
            </div>

          </div>

        </m.form>
      )}
    </AnimatePresence>
  );
});
