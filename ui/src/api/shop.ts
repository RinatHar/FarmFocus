import apiClient from "./client";






export const fetchBuyItem = async (id: number): Promise<number> => {
  const response = await apiClient.post(`/goods/${id}/buy`);

  return response.status;
}