import { type ReactNode, StrictMode } from "react";
import "./index.css";

import { ThemeProvider } from "@/components/theme-provider";
import { Toaster } from "@/components/ui/toaster";
import { NuqsAdapter } from "nuqs/adapters/react";
import { SWRConfig } from "swr";
import fetcher from "./lib/apis/fetcher";

export default function App({ children }: { children: ReactNode }) {
	return (
		<StrictMode>
			<ThemeProvider defaultTheme="system" storageKey="vite-ui-theme">
				<SWRConfig
					value={{
						// Define your global configuration options here
						fetcher: fetcher,
						// ...other global configurations...
					}}
				>
					<NuqsAdapter>{children}</NuqsAdapter>
					<Toaster />
				</SWRConfig>
			</ThemeProvider>
		</StrictMode>
	);
}
