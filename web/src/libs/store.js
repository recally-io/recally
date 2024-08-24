import { create } from "zustand";

const useStore = create((set) => ({
  isLogin: false,
  setIsLogin: (isLogin) => set(() => ({ isLogin: isLogin })),

  isDarkMode: false,
  setIsDarkMode: (isDarkMode) => set(() => ({ isDarkMode: isDarkMode })),

  isSidebarOpen: false,
  toggleSidebar: () =>
    set((state) => ({ isSidebarOpen: !state.isSidebarOpen })),

  assistantId: "",
  setAssistantId: (id) => set(() => ({ assistantId: id })),
  threadId: "",
  setThreadId: (id) => set(() => ({ threadId: id })),

  threadIsOpenSettings: false,
  setThreadIsOpenSettings: (isOpen) =>
    set(() => ({ threadIsOpenSettings: isOpen })),
  threadIsOpenModelSelecter: false,
  setThreadIsOpenModelSelecter: (isOpen) =>
    set(() => ({ threadIsOpenModelSelecter: isOpen })),
  threadIsTitleGenerated: false,
  setThreadIsTitleGenerated: (isGenerated) =>
    set(() => ({ threadIsTitleGenerated: isGenerated })),
  threadChatModel: "gpt-4o",
  setThreadChatModel: (model) => set(() => ({ threadChatModel: model })),
  threadMessageList: [],
  setThreadMessageList: (list) => set(() => ({ threadMessageList: list })),
  addThreadMessage: (message) =>
    set((state) => ({
      threadMessageList: [...state.threadMessageList, message],
    })),
  threadNewText: "",
  setThreadNewText: (text) => set(() => ({ threadNewText: text })),
  threadFileContent: "",
  setThreadFileContent: (content) =>
    set(() => ({ threadFileContent: content })),

  threadModels: [],
  setThreadModels: (models) => set(() => ({ threadModels: models })),

  threadSettings: {
    name: "New Thread",
    description: "",
    systemPrompt: "",
    temperature: 0.7,
    maxToken: 4096,
    model: (state) => state.threadChatModel,
  },
  setThreadSettings: (settings) => set(() => ({ threadSettings: settings })),
}));

export default useStore;
