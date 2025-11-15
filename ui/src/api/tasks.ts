import type { DoneTaskResponse } from "../types/api";
import type { ITask } from "../types/farm";
import { toLocalIsoWithOffset } from "../utils/toLocalIsoWithOffset";
import apiClient from "./client";

type TaskNoId = Omit<ITask, "id" | "done">;

export const fetchGetTasks = async (): Promise<ITask[]> => {
  const response = await apiClient.get("/tasks");

  return response.data;
}

export const fetchAddTasks = async (task: TaskNoId): Promise<ITask> => {
  const date = task.date ? toLocalIsoWithOffset(task.date) : null;

  const response = await apiClient.post("/tasks",
    {
      date: date,
      description: task.description || '',
      difficulty: task.difficulty,
      tagId: task.tag?.id,
      title: task.title,
    }
  );

  return response.data;
}

export const fetchEditTasks = async (task: ITask) => {
  const date = task.date ? toLocalIsoWithOffset(task.date) : null;

  const response = await apiClient.put(`/tasks/${task.id}`,
    {
      date: date,
      description: task.description || '',
      difficulty: task.difficulty,
      done: task.done,
      tagId: task.tag?.id,
      title: task.title,
    }
  );

  return response.data;
}

export const fetchDeleteTasks = async (id: number) => {
  const response = await apiClient.delete(`/tasks/${id}`);

  return response.data;
}

export const fetchDoneTask = async (id: number): Promise<DoneTaskResponse> => {
  const response = await apiClient.patch(`/tasks/${id}/done`);

  return response.data;
}

export const fetchUndoneTask = async (id: number): Promise<DoneTaskResponse> => {
  const response = await apiClient.patch(`/tasks/${id}/undone`);

  return response.data;
}
