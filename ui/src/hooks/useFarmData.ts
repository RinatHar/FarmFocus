import { useQuery } from "@tanstack/react-query";
import { fetchGetServerData } from "../api/syncData";

export const useFarmData = () => {

  return useQuery({
    queryKey: ["farm-state"],
    queryFn: fetchGetServerData,
  });

};