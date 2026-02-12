import { writable } from "svelte/store";
import { browser } from "$app/environment";

const STORAGE_KEY = "api_token";

let initialValue: string | null = null;

if (browser) {
  initialValue = localStorage.getItem(STORAGE_KEY);
}

const token = writable<string | null>(initialValue);

token.subscribe((value) => {
  if (browser) {
    if (value === null) {
      localStorage.removeItem(STORAGE_KEY);
    } else {
      localStorage.setItem(STORAGE_KEY, value);
    }
  }
});

export const getToken = (): string | null => {
  let tokenValue: string | null = null;
  token.subscribe((value) => (tokenValue = value))();
  return tokenValue;
};

export const setToken = (value: string | null) => {
  token.set(value);
};

export const clearToken = () => setToken(null);

export default token;
