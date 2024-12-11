import { Button } from "@/components/ui/button";

import { Github, Twitter } from "lucide-react";
import { Link } from "react-router-dom";
import { ThemeToggle } from "../ThemeToggle";

interface BaseLayoutProps {
  children: React.ReactNode;
}

export function BaseLayout({ children }: BaseLayoutProps) {
  return (
    <div className="flex flex-col min-h-screen bg-background">
      <header className="sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
        <div className="container mx-auto flex px-2 h-16 items-center justify-between">
          <div className="flex items-center space-x-4">
            <Link to="/" className="flex items-center space-x-2">
              <span className="font-bold text-xl">Vibrain</span>
            </Link>
            <nav className="hidden md:flex space-x-4">
              <Link
                to="/features"
                className="text-sm font-medium hover:text-primary"
              >
                Features
              </Link>
              <Link
                to="/pricing"
                className="text-sm font-medium hover:text-primary"
              >
                Pricing
              </Link>
              <Link
                to="/blog"
                className="text-sm font-medium hover:text-primary"
              >
                Blog
              </Link>
            </nav>
          </div>
          <div className="flex items-center space-x-2">
            <ThemeToggle />
            <Button variant="outline" size="sm">
              Sign In
            </Button>
          </div>
        </div>
      </header>

      <main className="flex-grow container py-8">{children}</main>

      <footer className="sticky bottom-0 z-50 w-full border-t w-full bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
        <div className="container mx-auto items-center justify-between">
          <div className="md:m-4 my-2 flex flex-col md:flex-row justify-between items-center">
            <p className="text-sm text-muted-foreground">
              © 2024 Vibrain. All rights reserved.
            </p>
            <div className="flex space-x-4 mt-4 md:mt-0">
              <Link
                to="https://twitter.com/vibrain"
                className="text-muted-foreground hover:text-primary"
              >
                <Twitter className="h-5 w-5" />
              </Link>
              <Link
                to="https://github.com/vibrain"
                className="text-muted-foreground hover:text-primary"
              >
                <Github className="h-5 w-5" />
              </Link>
            </div>
          </div>
        </div>
      </footer>
    </div>
  );
}
