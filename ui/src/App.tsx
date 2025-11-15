import './App.css';
import 'swiper/css';
import '@maxhub/max-ui/dist/styles.css';
import { Suspense, useEffect, useState } from 'react';
import { TodoList } from './components/TodoList/TodoList';
import { BottomMenu } from './components/BottomMenu/BottomMenu';
import { Field } from './components/Farm/Field';
import { Shop } from './components/Shop/Shop';
import { Swiper, SwiperSlide } from 'swiper/react';
import type { Swiper as SwiperClass } from 'swiper';
import { HeaderStatus } from './components/HederStatus/HeaderStatus';
import { EditTaskModal } from './components/Modal/Task/EditTaskModal';
import { AddTaskModal } from './components/Modal/Task/AddTaskModal';
import { useFarmData } from './hooks/useFarmData';
import { useFarmStore } from './stores/useFarmStore';
import { FilterModal } from './components/TodoList/ui/FilterModal';
import { HabitList } from './components/HabitList/HabitList';
import { EditHabitModal } from './components/Modal/Habit/EditHabitModal';
import { AddHabitModal } from './components/Modal/Habit/AddHabitModal';
import { TaskCalendar } from './components/Calendar/TaskCalendar';
import { useMaxWebApp } from './hooks/useMaxWebApp';

export type SwiperPage = "calendar" | "habit-list" | "todo-list" | "shop"

function App() {
  const { data, isLoading, error, } = useFarmData();
  const { isReady, user } = useMaxWebApp();
  const setFromServer = useFarmStore((state) => state.setFromServer);

useEffect(() => {
  if (isReady) {
    useFarmStore.getState().setUserId(user?.id ?? 1);
  }
}, [isReady, user?.id]);

  useEffect(() => {
    if (data) {
      setFromServer(data);
    }
  }, [data, setFromServer]);

  const [swiperInstance, setSwiperInstance] = useState<SwiperClass | null>(null);
  const [currentPage, setCurrentPage] = useState<SwiperPage>("todo-list");

  const handlePageChange = (page: SwiperPage) => {
    if (!swiperInstance) return;
    if (page === 'calendar') swiperInstance.slideTo(0);
    if (page === 'habit-list') swiperInstance.slideTo(1);
    if (page === 'todo-list') swiperInstance.slideTo(2);
    if (page === 'shop') swiperInstance.slideTo(3);
    setCurrentPage(page);
  };
  
  if (isLoading) {
    return (
      <div className="flex flex-col items-center justify-center h-screen bg-base-200 text-base-content">
        <span className="loading loading-spinner loading-lg mb-4"></span>
        <p className="font-medium text-lg animate-pulse">Загрузка фермы...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex flex-col items-center justify-center h-screen bg-base-200 text-error text-center">
        <p className="text-xl font-semibold">Ошибка при загрузке данных 😢</p>
        <p className="text-sm opacity-80 mt-2">{String(error)}</p>
      </div>
    );
  }


  return (
    <div className="bg-base-100 py-2 overflow-hidden flex flex-col">

      <HeaderStatus />

      <div className="flex-1 overflow-hidden">
        <Field />
        <Swiper
          threshold={10}
          spaceBetween={50}
          slidesPerView={1}
          initialSlide={2}
          onSwiper={(swiper) => setSwiperInstance(swiper)}
          onSlideChange={(swiper) => {
            switch (swiper.activeIndex){
              case 0:
                setCurrentPage("calendar");
                break;
              case 1:
                setCurrentPage("habit-list");
                break;
              case 2:
                setCurrentPage("todo-list");
                break;
              case 3:
                setCurrentPage("shop");
                break;
            }
          }}
        >
          <SwiperSlide>
            <TaskCalendar />
          </SwiperSlide>

          <SwiperSlide>
            <HabitList />
          </SwiperSlide>

          <SwiperSlide>
            <TodoList />
          </SwiperSlide>

          <SwiperSlide>
            <Shop />
          </SwiperSlide>
        </Swiper>
      </div>

      <BottomMenu
        onChangePage={handlePageChange}
        currentPage={currentPage}
      />

    {(currentPage === "todo-list" || currentPage === "calendar") && (
      <Suspense>
        <AddTaskModal />
        <EditTaskModal />
        <FilterModal />
      </Suspense>
    )}
    {currentPage === "habit-list" && (
      <Suspense>
        <AddHabitModal />
        <EditHabitModal />
      </Suspense>
    )}
    </div>
  );
}

export default App;
