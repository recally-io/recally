import { create } from "zustand";

const useStore = create((set) => ({
  user: null,
  isLogin: false,
  threads: [],

  setIsLogin: (isLogin) => set(() => ({ isLogin })),

  // 更新 user 的操作
  setUser: (user) => set(() => ({ user })),

  // 切换 theme 的操作
  toggleTheme: () =>
    set((state) => ({
      theme: state.theme === "light" ? "dark" : "light",
    })),
}));

export default useStore;
