export const API_BASE_URL = import.meta.env.VITE_API_URL;


export const client = async <T>(endpoint: string): Promise<T> => {
  const response = await fetch(`${API_BASE_URL}${endpoint}`);
  if (!response.ok) {
    throw new Error(`Erreur HTTP api: ${response.status}`);
  }

  return response.json() as Promise<T>;
};