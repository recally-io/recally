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
  assistant: {},
  setAssistant: (assistant) => set(() => ({ assistant: assistant })),
  threadId: "",
  setThreadId: (id) => set(() => ({ threadId: id })),
  thread: {},
  setThread: (thread) => set(() => ({ thread: thread })),

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

  threadTools: [],
  setThreadTools: (tools) => set(() => ({ threadTools: tools })),

  threadSettings: {
    name: "New Thread",
    description: "",
    system_prompt: "",
    temperature: 0.7,
    max_token: 4096,
    model: (state) => state.threadChatModel,
    metadata: {
      tools: [],
    },
  },
  setThreadSettings: (settings) => set(() => ({ threadSettings: settings })),

  threadChatImages: [],
  setThreadChatImages: (images) => set(() => ({ threadChatImages: images })),
  addThreadChatImage: (image) =>
    set((state) => ({
      threadChatImages: [...state.threadChatImages, image],
    })),

  files: [],
  setFiles: (files) => set(() => ({ files: files })),
  addFile: (file) => set((state) => ({ files: [...state.files, file] })),
}));

export default useStore;
