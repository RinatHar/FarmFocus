import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import App from './App.tsx'
import { MaxUI } from '@maxhub/max-ui'
import { Toaster } from 'sonner'
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { domAnimation, LazyMotion } from 'framer-motion'

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: 1,
      refetchOnWindowFocus: false,
      staleTime: 1000 * 60,
    },
  },
});

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <LazyMotion features={domAnimation}>
      <QueryClientProvider client={queryClient}>
        <MaxUI>
          <App />
          <Toaster position='bottom-center' duration={ 1000 } />
        </MaxUI>
      </QueryClientProvider>
    </LazyMotion>
    <script src="https://st.max.ru/js/max-web-app.js"></script>
  </StrictMode>,
)
