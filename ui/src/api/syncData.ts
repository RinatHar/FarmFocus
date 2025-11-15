import type { ServerFarmData } from "../types/api";
import apiClient from "./client";





export const fetchGetServerData = async (): Promise<ServerFarmData> => {
  const response = await apiClient.get("/users/sync");

  return response.data;
}
