import { getToken } from "$lib/stores/token";

export class ApiClient {
  private apiUrl: string;

  constructor(apiUrl: string) {
    this.apiUrl = apiUrl;
  }

  private getHeaders(): HeadersInit {
    const token = getToken();
    const headers: Record<string, string> = {
      "Content-Type": "application/json",
    };
    if (token) {
      headers.Authorization = `Bearer ${token}`;
    }
    return headers;
  }

  private async handleResponse(response: Response) {
    if (!response.ok) {
      throw new Error(
        `request failed: ${response.status} ${response.statusText}`,
      );
    }
    return response;
  }

  async get<T>(path: string): Promise<T> {
    const response = await fetch(`${this.apiUrl}${path}`, {
      headers: this.getHeaders(),
    });

    await this.handleResponse(response);
    return response.json();
  }

  async post<T>(path: string, body: unknown): Promise<T> {
    const response = await fetch(`${this.apiUrl}${path}`, {
      method: "POST",
      headers: this.getHeaders(),
      body: JSON.stringify(body),
    });

    await this.handleResponse(response);
    return response.json();
  }

  async put<T>(path: string, body: unknown): Promise<T> {
    const response = await fetch(`${this.apiUrl}${path}`, {
      method: "PUT",
      headers: this.getHeaders(),
      body: JSON.stringify(body),
    });

    await this.handleResponse(response);
    return response.json();
  }

  async delete<T>(path: string): Promise<T> {
    const response = await fetch(`${this.apiUrl}${path}`, {
      method: "DELETE",
      headers: this.getHeaders(),
    });

    await this.handleResponse(response);
    return response.json();
  }
}
