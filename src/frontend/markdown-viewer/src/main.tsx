import React from "react";
import ReactDOM from "react-dom/client";
import MarkdownViewer from "@/MarkdownViewer";
import "@/index.css"; // Tailwind CSS
import "@/markdown-styles.css"; // Markdown-specific styles
import { testMarkdownContent, shortMarkdownContent } from "@/test/testContent";

// Development mode: inject test content
if (import.meta.env.DEV) {
  console.log("🚀 Running in development mode");

  // Inject test markdown content
  window.markdownContent = testMarkdownContent;

  // Add dev tools to window for testing
  (window as any).devTools = {
    setContent: (content: string) => {
      window.markdownContent = content;
      // Force re-render
      const rootElement = document.getElementById("root");
      if (rootElement) {
        rootElement.innerHTML = "";
        ReactDOM.createRoot(rootElement).render(
          <React.StrictMode>
            <MarkdownViewer />
          </React.StrictMode>
        );
      }
    },
    loadTestContent: () => {
      (window as any).devTools.setContent(testMarkdownContent);
    },
    loadShortContent: () => {
      (window as any).devTools.setContent(shortMarkdownContent);
    },
    loadCustomContent: (content: string) => {
      (window as any).devTools.setContent(content);
    },
  };

  console.log("📝 Test content loaded!");
  console.log("🛠️  Available dev tools:");
  console.log("  - devTools.loadTestContent() - Load full test document");
  console.log("  - devTools.loadShortContent() - Load short test");
  console.log("  - devTools.loadCustomContent(content) - Load custom markdown");
}

const rootElement = document.getElementById("root");

if (rootElement) {
  ReactDOM.createRoot(rootElement).render(
    <React.StrictMode>
      <MarkdownViewer />
    </React.StrictMode>
  );
}
