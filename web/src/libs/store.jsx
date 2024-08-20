import { create } from "zustand";

const useStore = create((set) => ({
  isLogin: false,
  setIsLogin: (isLogin) => set(() => ({ isLogin: isLogin })),
}));

export default useStore;
