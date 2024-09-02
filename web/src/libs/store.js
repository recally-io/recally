import { create } from "zustand";

const useStore = create((set) => ({
  isLogin: false,
  setIsLogin: (isLogin) => set(() => ({ isLogin: isLogin })),

  isDarkMode: false,
  setIsDarkMode: (isDarkMode) => set(() => ({ isDarkMode: isDarkMode })),

  mobileSidebarOpen: false,
  toggleMobileSidebar: () =>
    set((state) => ({ mobileSidebarOpen: !state.mobileSidebarOpen })),

  desktopSidebarOpen: true,
  toggleDesktopSidebar: () =>
    set((state) => ({ desktopSidebarOpen: !state.desktopSidebarOpen })),

  assistantId: "",
  setAssistantId: (id) => set(() => ({ assistantId: id })),
  threadId: "",
  setThreadId: (id) => set(() => ({ threadId: id })),

  threadIsOpenSettings: false,
  setThreadIsOpenSettings: (isOpen) =>
    set(() => ({ threadIsOpenSettings: isOpen })),

  threadMessageList: [],
  setThreadMessageList: (list) => set(() => ({ threadMessageList: list })),
  addThreadMessage: (message) =>
    set((state) => ({
      threadMessageList: [...state.threadMessageList, message],
    })),

  threadNewText: "",
  setThreadNewText: (text) => set(() => ({ threadNewText: text })),

  threadChatImages: [],
  setThreadChatImages: (images) => set(() => ({ threadChatImages: images })),
  addThreadChatImage: (image) =>
    set((state) => ({
      threadChatImages: [...state.threadChatImages, image],
    })),
}));

export default useStore;
