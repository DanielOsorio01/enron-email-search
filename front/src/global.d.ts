// global.d.ts
declare global {
  interface Window {
    APP_CONFIG: {
      backend_URL: string;
    };
  }
}

export {};
