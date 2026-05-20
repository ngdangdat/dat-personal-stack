// Simple IndexedDB wrapper for configuration persistence
const DB_NAME = 'EngineeringAssistant';
const DB_VERSION = 1;
const STORE_NAME = 'config';

class DBService {
  constructor() {
    this.db = null;
  }

  init() {
    if (this.db) return Promise.resolve(this.db);

    return new Promise((resolve, reject) => {
      const request = indexedDB.open(DB_NAME, DB_VERSION);

      request.onerror = (event) => {
        console.error('IndexedDB open error:', event.target.error);
        reject(event.target.error);
      };

      request.onsuccess = (event) => {
        this.db = event.target.result;
        resolve(this.db);
      };

      request.onupgradeneeded = (event) => {
        const db = event.target.result;
        if (!db.objectStoreNames.contains(STORE_NAME)) {
          db.createObjectStore(STORE_NAME, { keyPath: 'key' });
        }
      };
    });
  }

  async getConfig() {
    const db = await this.init();
    return new Promise((resolve, reject) => {
      const transaction = db.transaction(STORE_NAME, 'readonly');
      const store = transaction.objectStore(STORE_NAME);
      const request = store.get('github_config');

      request.onerror = (event) => reject(event.target.error);
      request.onsuccess = (event) => {
        const result = event.target.result;
        resolve(result ? result.value : null);
      };
    });
  }

  async saveConfig(value) {
    const db = await this.init();
    return new Promise((resolve, reject) => {
      const transaction = db.transaction(STORE_NAME, 'readwrite');
      const store = transaction.objectStore(STORE_NAME);
      const request = store.put({ key: 'github_config', value });

      request.onerror = (event) => reject(event.target.error);
      request.onsuccess = () => resolve();
    });
  }
}

export const db = new DBService();
