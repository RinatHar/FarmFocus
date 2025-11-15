import type { DoneHabitsResponse } from "../types/api";
import type { IHabit } from "../types/farm";
import { toLocalIsoWithOffset } from "../utils/toLocalIsoWithOffset";
import apiClient from "./client";

type HabitNoId = Omit<IHabit, "id" | "done" | "count">;

export const fetchGetHabits = async (): Promise<IHabit[]> => {
  const response = await apiClient.get("/habits");
  return response.data;
};

export const fetchAddHabit = async (habit: HabitNoId): Promise<IHabit> => {
  const startDate = habit.startDate
    ? toLocalIsoWithOffset(habit.startDate)
    : null;

  const response = await apiClient.post("/habits", {
    title: habit.title,
    description: habit.description || "",
    difficulty: habit.difficulty,
    period: habit.period,
    every: habit.every,
    startDate,
    tagId: habit.tag?.id,
    count: 0,
  });

  return response.data;
};

export const fetchEditHabit = async (habit: IHabit): Promise<IHabit> => {
  const startDate = habit.startDate
    ? toLocalIsoWithOffset(habit.startDate)
    : null;

  const response = await apiClient.put(`/habits/${habit.id}`, {
    title: habit.title,
    description: habit.description || "",
    difficulty: habit.difficulty,
    done: habit.done,
    count: habit.count,
    period: habit.period,
    every: habit.every,
    startDate,
    tagId: habit.tag?.id,
  });

  return response.data;
};

export const fetchDeleteHabit = async (id: number) => {
  const response = await apiClient.delete(`/habits/${id}`);
  return response.data;
};

export const fetchDoneHabit = async (id: number): Promise<DoneHabitsResponse> => {
  const response = await apiClient.patch(`/habits/${id}/done`);
  return response.data;
};

export const fetchUndoneHabit = async (id: number): Promise<DoneHabitsResponse> => {
  const response = await apiClient.patch(`/habits/${id}/undone`);
  return response.data;
};
