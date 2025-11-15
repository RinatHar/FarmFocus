import type { HarvestPlantResponse } from "../types/api";
import type { IPlant } from "../types/farm";
import apiClient from "./client";





export const fetchSetPlant = async (cellNumber: number, seedId: number): Promise<IPlant> => {
  const response = await apiClient.post("/user-plants",
    {
      cellNumber,
      seedId,
    });

  return response.data;
}

export const fetchHarvestPlant = async (id: number): Promise<HarvestPlantResponse> => {
  const response = await apiClient.post(`/user-plants/${id}/harvest`);

  return response.data;
}