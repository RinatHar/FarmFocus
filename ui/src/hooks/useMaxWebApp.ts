// src/hooks/useMaxWebApp.ts
import { useEffect, useState, useCallback } from 'react';

interface MaxUser {
  id: number;
  first_name: string;
  last_name?: string;
  username?: string;
  language_code?: string;
  photo_url?: string;
}

interface MaxInitData {
  user?: MaxUser;
  query_id?: string;
  chat?: {
    id: number;
    type: string;
  };
  start_param?: string;
}

declare global {
  interface Window {
    WebApp?: {
      ready: () => void;
      initDataUnsafe: MaxInitData;
      initData: string;
      platform: string;
      version: string;
    };
  }
}

export const useMaxWebApp = () => {
  const [isReady, setIsReady] = useState(false);
  const [user, setUser] = useState<MaxUser | null>(null);
  const [initData, setInitData] = useState<MaxInitData | null>(null);
  const [platform, setPlatform] = useState<string | null>(null);

  const ready = useCallback(() => {
    if (window.WebApp) {
      window.WebApp.ready();
      setIsReady(true);
    }
  }, []);

  useEffect(() => {
    if (!window.WebApp) {
      console.warn('MAX WebApp SDK не загружен.');
      return;
    }

    const data = window.WebApp.initDataUnsafe;
    const plat = window.WebApp.platform;

    setInitData(data);
    setUser(data.user || null);
    setPlatform(plat);
    ready();
  }, [ready]);

  return {
    /** SDK инициализирован и вызван ready() */
    isReady,

    /** Данные пользователя */
    user,

    /** Полные данные инициализации (включая query_id, chat и т.д.) */
    initData,

    /** Платформа: ios | android | desktop | web */
    platform,

    /** Принудительно вызвать ready() (если нужно) */
    ready,
  };
};