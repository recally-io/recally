import { create } from "zustand";

const useStore = create((set) => ({
  isLogin: false,
  setIsLogin: (isLogin) => set(() => ({ isLogin: isLogin })),

  isDarkMode: false,
  setIsDarkMode: (isDarkMode) => set(() => ({ isDarkMode: isDarkMode })),

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
  threadNewText: "",
  setThreadNewText: (text) => set(() => ({ threadNewText: text })),
  threadFileContent: "",
  setThreadFileContent: (content) =>
    set(() => ({ threadFileContent: content })),

  threadModels: [],
  setThreadModels: (models) => set(() => ({ threadModels: models })),
}));

export default useStore;
