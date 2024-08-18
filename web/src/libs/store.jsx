import { act } from "react";
import { create } from "zustand";

const useStore = create((set) => ({
  isLogin: false,
  setIsLogin: (isLogin) => set(() => ({ isLogin })),

  activateThreadId: "",
  setActivateThreadId: (activateThreadId) => set(() => ({ activateThreadId })),

  activateAssistantId: "",
  setActivateAssistantId: (activateAssistantId) =>
    set(() => ({ activateAssistantId })),
}));

export default useStore;
