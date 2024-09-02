import {
  AppShell,
  createTheme,
  Input,
  MantineProvider,
  virtualColor,
} from "@mantine/core";
import "@mantine/core/styles.css";
import { ModalsProvider } from "@mantine/modals";
import { Notifications } from "@mantine/notifications";
import "@mantine/notifications/styles.css";
import { QueryClientProvider } from "@tanstack/react-query";
import { QueryContextProvider } from "../libs/query-context";

import React from "react";
import { queryClient } from "../libs/api";
import useStore from "../libs/store";
import Header from "./header";

const theme = createTheme({
  components: {
    Input: Input.extend({
      styles: {
        input: {
          fontSize: "16px",
        },
      },
    }),
  },
  // autoContrast: true,
  primaryColor: "primary",
  colors: {
    // Default: For neutral, non-emphasized UI elements
    // Use for:
    // - Page backgrounds (light mode: light gray, dark mode: dark gray or near-black)
    // - Text color for body content (light mode: dark gray on light background, dark mode: light gray on dark background)
    // - Borders for input fields, cards, or dividers
    // - Inactive tabs or menu items
    default: virtualColor({
      name: "default",
      dark: "dark",
      light: "gray",
    }),

    // Primary: Main brand color, used for key interactive elements
    // Use for:
    // - Primary action buttons (e.g., "Submit", "Save", "Continue")
    // - Links within text
    // - Currently selected tab or menu item
    // - Progress bars
    // - Key data visualizations or charts
    primary: virtualColor({
      name: "primary",
      dark: "blue",
      light: "blue",
    }),

    // Secondary: Complementary to primary, for less emphasized interactions
    // Use for:
    // - Secondary action buttons (e.g., "Cancel", "Back")
    // - Alternate row backgrounds in tables
    // - Borders of selected items
    // - Secondary data in charts or graphs
    secondary: virtualColor({
      name: "secondary",
      dark: "indigo",
      light: "indigo",
    }),

    // Success: Indicates positive outcomes or completed actions
    // Use for:
    // - Success messages or notifications
    // - Checkmarks for completed tasks
    // - "Publish" or "Approve" buttons
    // - Positive trends in data visualizations
    // - Verified or active status indicators
    success: virtualColor({
      name: "success",
      dark: "green",
      light: "green",
    }),

    // Warning: Draws attention to potential issues or important notices
    // Use for:
    // - Warning messages that don't prevent further action
    // - "Caution" or "Are you sure?" prompts
    // - Highlighting changed or unsaved content
    // - Indicating a mid-range value in a scale or rating system
    warning: virtualColor({
      name: "warning",
      dark: "yellow",
      light: "orange",
    }),

    // Danger: Highlights critical issues or destructive actions
    // Use for:
    // - Error messages
    // - Delete or remove buttons
    // - Critical alerts that require immediate attention
    // - Highlighting validation errors in forms
    // - Indicating dangerously high values in monitoring systems
    danger: virtualColor({
      name: "danger",
      dark: "red",
      light: "red",
    }),

    // Info: For neutral informational content
    // Use for:
    // - Informational messages or notifications
    // - Help or hint text
    // - Tooltips
    // - Icons for "More Info" buttons
    // - Highlighting new features or content
    info: virtualColor({
      name: "info",
      dark: "cyan",
      light: "cyan",
    }),

    // Accent: For drawing attention without implying status
    // Use for:
    // - Highlighting selected items in a list
    // - Decorative elements that should stand out
    // - Badges or tags for categorization
    // - Call-to-action buttons that aren't the primary action
    // - Accent elements in data visualizations
    accent: virtualColor({
      name: "accent",
      dark: "violet",
      light: "violet",
    }),

    // Muted: For de-emphasized or disabled elements
    // Use for:
    // - Disabled buttons or form fields
    // - Secondary or helper text
    // - Backgrounds of inactive elements
    // - Less important grid lines in charts
    // - Placeholder text in input fields
    muted: virtualColor({
      name: "muted",
      dark: "gray",
      light: "gray",
    }),
  },
});

export default function Layout({ main, nav = null, header = null }) {
  let hasNavBar = nav !== null;
  const mobileSidebarOpen = useStore((state) => state.mobileSidebarOpen);
  const desktopSidebarOpen = useStore((state) => state.desktopSidebarOpen);

  return (
    <QueryClientProvider client={queryClient}>
      <QueryContextProvider>
        <MantineProvider theme={theme} defaultColorScheme="auto">
          <ModalsProvider>
            <Notifications />
            <AppShell
              header={{ height: 40 }}
              navbar={{
                width: "300",
                breakpoint: "sm",
                collapsed: {
                  mobile: !hasNavBar || !mobileSidebarOpen,
                  desktop: !hasNavBar || !desktopSidebarOpen,
                },
              }}
              padding="0"
              withBorder={true}
              layout="alt"
            >
              <AppShell.Header>
                {header ? header : <Header hasNavBar={hasNavBar} />}
              </AppShell.Header>
              <AppShell.Navbar>{nav}</AppShell.Navbar>
              <AppShell.Main>{main}</AppShell.Main>
            </AppShell>
          </ModalsProvider>
        </MantineProvider>
      </QueryContextProvider>
    </QueryClientProvider>
  );
}
