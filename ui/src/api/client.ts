// src/api/client.ts
import axios from 'axios';
import { useFarmStore } from '../stores/useFarmStore';

const apiClient = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL,
  timeout: 10000,
});

apiClient.interceptors.request.use((config) => {
  const userId = useFarmStore.getState().userId;

  if (userId && userId !== 0) {
    config.headers['X-User-ID'] = String(userId);
  } else {
    delete config.headers['X-User-ID'];
  }

  return config;
});

export default apiClient;