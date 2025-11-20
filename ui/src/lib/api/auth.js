import { sendRequest } from "./api";

export class AuthApi {
  static basePath = "/auth";

  static async me() {
    return await sendRequest(`${this.basePath}/me`, {
      method: "GET",
    });
  }

  static async login(email, password) {
    return await sendRequest(`${this.basePath}/login`, {
      method: "POST",
      body: JSON.stringify({ email, password }),
      headers: {
        "Content-Type": "application/json",
      },
    });
  }

  static async register(email, password) {
    return await sendRequest(`${this.basePath}/register`, {
      method: "POST",
      body: JSON.stringify({ email, password }),
      headers: {
        "Content-Type": "application/json",
      },
    });
  }
}
