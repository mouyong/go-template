import { useState, useEffect } from 'react';

const STORAGE_KEY = 'sidebar_collapsed';

export default function useSidebarCollapsed() {
  const [collapsed, setCollapsed] = useState(() => {
    const stored = localStorage.getItem(STORAGE_KEY);
    return stored === 'true';
  });

  const toggle = () => {
    setCollapsed((prev) => !prev);
  };

  useEffect(() => {
    localStorage.setItem(STORAGE_KEY, String(collapsed));
  }, [collapsed]);

  return [collapsed, toggle, setCollapsed];
}
