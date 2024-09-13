import { create } from "zustand";

export const defaultThreadSettings = {
  name: "Name",
  description: "Assistant description",
  system_prompt: "You are a helpful assistant.",
  model: "gpt-4o-mini",
  metadata: {
    tools: [],
  },
};

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
  toggleThreadIsOpenSettings: () =>
    set((state) => ({ threadIsOpenSettings: !state.threadIsOpenSettings })),

  threadMessageList: [],
  setThreadMessageList: (list) => set(() => ({ threadMessageList: list })),
  addThreadMessage: (message) =>
    set((state) => ({
      threadMessageList: [...state.threadMessageList, message],
    })),
  updateLastThreadMessage: (message) =>
    set((state) => ({
      threadMessageList: [...state.threadMessageList.slice(0, -1), message],
    })),

  threadSettings: defaultThreadSettings,
  setThreadSettings: (settings) => set(() => ({ threadSettings: settings })),
  resetThreadSettings: () =>
    set(() => ({ threadSettings: defaultThreadSettings })),
}));

export default useStore;
