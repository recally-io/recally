import { ThemeProvider } from "@/components/theme-provider";
import { Route, BrowserRouter as Router, Routes } from 'react-router-dom';
import BookmarkPage from './pages/BookmarkPage';
import HomePage from './pages/HomePage';

export default function App() {
  return (
    <ThemeProvider attribute="class" defaultTheme="system" enableSystem>
      <Router>
        <Routes>
          <Route path="/" element={<HomePage />} />
          <Route path="/bookmarks/:id" element={<BookmarkPage />} />
        </Routes>
      </Router>
    </ThemeProvider>
  )
}

