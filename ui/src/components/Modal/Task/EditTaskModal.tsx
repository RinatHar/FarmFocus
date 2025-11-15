import { m, AnimatePresence } from "framer-motion";
import { ArrowLeft } from "lucide-react";
import { Controller, useForm } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { lazy, memo, Suspense, useEffect } from "react";
import clsx from "clsx";
import { useEditTaskModalStore } from "../../../stores/editTaskModal";
import { useFarmStore } from "../../../stores/useFarmStore";
import { difficultyColor } from "../../../utils/difficultyColor";
import { InputTitle } from "../../TaskForm/InputTitle";
import { TextAreaDescription } from "../../TaskForm/TextAreaDescription";
import { DifficultySelector } from "../../TaskForm/DifficultySelector";
import { TagSelect } from "../../TaskForm/TagSelect";

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

  date: z.date().optional().nullable(),

  tag: z
    .object({
      id: z.number(),
      name: z.string(),
      color: z.string(),
    })
    .optional().nullable(),
});

export type EditTaskFormValues = z.infer<typeof newTaskSchema>;

export const EditTaskModal = memo(() => {
  const { isOpen, close, idTask } = useEditTaskModalStore();
  const { tasks, editTask, removeTask } = useFarmStore();

  const currentTask = tasks.find(t => t.id === idTask)
  
  const {
    handleSubmit,
    reset,
    control,
    register,
    formState: { errors },
  } = useForm<EditTaskFormValues>({
    resolver: zodResolver(newTaskSchema),
    defaultValues: {
      title: currentTask?.title,
      description: currentTask?.description,
      difficulty: currentTask?.difficulty,
      date: currentTask?.date || null,
      tag: currentTask?.tag || null,
    },
  });

  useEffect(() => {
    reset(currentTask)
  }, [currentTask, reset])
  
  const onSubmit = (data: EditTaskFormValues) => {
    editTask(
      idTask,
      data.title,
      data.description,
      data.difficulty,
      data.date || null,
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
    removeTask(idTask);
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
            difficultyColor(currentTask?.difficulty)
          )}>
            <button
              type="button"
              onClick={handleClose}
              className={clsx(
                "rounded-full p-1 transition-colors",
                {
                  "text-slate-50 hover:text-slate-200 hover:bg-slate-600": currentTask?.difficulty === "trifle",
                  "text-emerald-50 hover:text-emerald-200 hover:bg-emerald-600": currentTask?.difficulty === "easy",
                  "text-sky-50 hover:text-sky-200 hover:bg-sky-600": currentTask?.difficulty === "normal",
                  "text-orange-50 hover:text-orange-200 hover:bg-orange-600": currentTask?.difficulty === "hard",
                } 
              )}
            >
              <ArrowLeft />
            </button>
            <span className={clsx(
              "font-semibold",
              {
                "text-slate-50": currentTask?.difficulty === "trifle",
                "text-emerald-50": currentTask?.difficulty === "easy",
                "text-sky-50": currentTask?.difficulty === "normal",
                "text-orange-50": currentTask?.difficulty === "hard",
              }
            )}>
              Редактировать задачу
            </span>

            <div className="flex gap-2">
              <button
                type="button"
                onClick={ deleteTask }
                className={clsx(
                  "btn btn-sm font-mono outline-none border-0",
                  {
                    "bg-slate-50 text-slate-700 hover:bg-slate-100": currentTask?.difficulty === "trifle",
                    "bg-emerald-50 text-emerald-700 hover:bg-emerald-100": currentTask?.difficulty === "easy",
                    "bg-sky-50 text-sky-700 hover:bg-sky-100": currentTask?.difficulty === "normal",
                    "bg-orange-50 text-orange-700 hover:bg-orange-100": currentTask?.difficulty === "hard",
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
                    "bg-slate-50 text-slate-700 hover:bg-slate-100": currentTask?.difficulty === "trifle",
                    "bg-emerald-50 text-emerald-700 hover:bg-emerald-100": currentTask?.difficulty === "easy",
                    "bg-sky-50 text-sky-700 hover:bg-sky-100": currentTask?.difficulty === "normal",
                    "bg-orange-50 text-orange-700 hover:bg-orange-100": currentTask?.difficulty === "hard",
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
              difficultyColor(currentTask?.difficulty)
            )}>
              <InputTitle
                register={register}
                difficulty={currentTask?.difficulty}
              />
              
              <TextAreaDescription
                register={register}
                difficulty={currentTask?.difficulty}
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
                <p>Дата выполнения</p>
                <Controller
                  name="date"
                  control={control}
                  render={({ field: { value, onChange } }) => (
                    <Suspense>
                      <DatePicker value={value} onChange={onChange} />
                    </Suspense>
                  )}
                />
                {errors.date && (
                  <p className="text-error text-sm mt-1">
                    {errors.date.message}
                  </p>
                )}
              </div>
            </div>

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
                        difficulty={currentTask?.difficulty}
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
