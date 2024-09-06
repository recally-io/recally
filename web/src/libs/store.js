import { create } from "zustand";

const useStore = create((set) => ({
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
}));

export default useStore;
