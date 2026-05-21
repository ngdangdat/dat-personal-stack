// Simple IndexedDB wrapper for configuration persistence
const DB_NAME = 'EngineeringAssistant';
const DB_VERSION = 2;
const CONFIG_STORE = 'config';
const NODES_STORE = 'nodes';

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
        if (!db.objectStoreNames.contains(CONFIG_STORE)) {
          db.createObjectStore(CONFIG_STORE, { keyPath: 'key' });
        }
        if (!db.objectStoreNames.contains(NODES_STORE)) {
          const store = db.createObjectStore(NODES_STORE, { keyPath: 'id' });
          store.createIndex('name', 'name', { unique: false });
          store.createIndex('address', 'address', { unique: false });
        }
      };
    });
  }

  async getConfig() {
    const db = await this.init();
    return new Promise((resolve, reject) => {
      const transaction = db.transaction(CONFIG_STORE, 'readonly');
      const store = transaction.objectStore(CONFIG_STORE);
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
      const transaction = db.transaction(CONFIG_STORE, 'readwrite');
      const store = transaction.objectStore(CONFIG_STORE);
      const request = store.put({ key: 'github_config', value });

      request.onerror = (event) => reject(event.target.error);
      request.onsuccess = () => resolve();
    });
  }

  // Node operations
  async getNodes() {
    const db = await this.init();
    return new Promise((resolve, reject) => {
      const transaction = db.transaction(NODES_STORE, 'readonly');
      const store = transaction.objectStore(NODES_STORE);
      const request = store.getAll();

      request.onerror = (event) => reject(event.target.error);
      request.onsuccess = (event) => resolve(event.target.result || []);
    });
  }

  async saveNode(node) {
    const db = await this.init();
    return new Promise((resolve, reject) => {
      const transaction = db.transaction(NODES_STORE, 'readwrite');
      const store = transaction.objectStore(NODES_STORE);
      const request = store.put(node);

      request.onerror = (event) => reject(event.target.error);
      request.onsuccess = () => resolve();
    });
  }

  async deleteNode(id) {
    const db = await this.init();
    return new Promise((resolve, reject) => {
      const transaction = db.transaction(NODES_STORE, 'readwrite');
      const store = transaction.objectStore(NODES_STORE);
      const request = store.delete(id);

      request.onerror = (event) => reject(event.target.error);
      request.onsuccess = () => resolve();
    });
  }
}

export const db = new DBService();
