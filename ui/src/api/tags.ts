import type { ITag } from "../types/farm";
import apiClient from "./client";

type TagNoId = Omit<ITag, "id">;

export const fetchCreatedTag = async (tag: TagNoId) => {
  const response = await apiClient.post("/tags",
    tag
  );

  return response.data;
}

export const fetchUpdateTag = async (tag: ITag) => {
  const response = await apiClient.put(`/tags/${tag.id}`, 
    {name: tag.name, color: tag.color}
  );

  return response.data;
}

export const fetchDeleteTag = async (id: number) => {
  const response = await apiClient.delete(`/tags/${id}`);

  return response.data;
}