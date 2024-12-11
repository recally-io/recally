import { Github, Twitter } from "lucide-react";
import { Link } from "react-router-dom";

export default function Footer() {
  return (
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
  );
}
