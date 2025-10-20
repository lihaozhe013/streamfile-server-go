import React from 'react';
import ReactDOM from 'react-dom/client';
import MarkdownViewer from '@/MarkdownViewer';
import '@/markdown-styles.css';

const rootElement = document.getElementById('root');

if (rootElement) {
  ReactDOM.createRoot(rootElement).render(
    <React.StrictMode>
      <MarkdownViewer />
    </React.StrictMode>
  );
}
